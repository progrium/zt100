package fswatch

import (
	"context"
	"log"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/progrium/watcher"
)

type Watch struct {
	Path    string
	Handler func(watcher.Event)
}

type Service struct {
	watcher *watcher.Watcher
	tmpDir  string
	watches map[*Watch]struct{}
	mu      sync.Mutex
}

func (s *Service) Serve(ctx context.Context) {
	for {
		select {
		case event := <-s.watcher.Event:
			log.Println(event)
			for _, w := range s.matchWatches(event.Path) {
				w.Handler(event)
			}
		case err := <-s.watcher.Error:
			if err == watcher.ErrWatchedFileDeleted {
				log.Println(err)
			} else {
				panic(err)
			}
		case <-s.watcher.Closed:
			return
		case <-ctx.Done():
			return
		}
	}
}

func (s *Service) matchWatches(path string) (watches []*Watch) {
	s.mu.Lock()
	defer s.mu.Unlock()
	for w := range s.watches {
		if strings.HasPrefix(path, w.Path) {
			watches = append(watches, w)
		}
	}
	return
}

func (s *Service) Tempfile(name string, handler func(watcher.Event)) (*Watch, error) {
	w := &Watch{
		Path:    filepath.Join(s.tmpDir, name),
		Handler: handler,
	}
	f, err := os.Create(w.Path)
	if err != nil {
		return nil, err
	}
	f.Close()
	return w, s.Watch(w)
}

func (s *Service) Watch(w *Watch) error {
	w.Path, _ = filepath.Abs(w.Path)
	files := s.watcher.WatchedFiles()
	_, exists := files[w.Path]
	if !exists {
		if err := s.watcher.Add(w.Path); err != nil {
			return err
		}
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	s.watches[w] = struct{}{}
	return nil
}

func (s *Service) Unwatch(w *Watch) error {
	watches := s.matchWatches(w.Path)
	if len(watches) == 1 {
		if err := s.watcher.RemoveRecursive(w.Path); err != nil {
			return err
		}
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.watches, w)
	return nil
}
