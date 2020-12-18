package fswatch

import (
	"context"
	"io/ioutil"
	"log"
	"os"
	"time"

	"github.com/progrium/watcher"
)

func (s *Service) InitializeDaemon() error {
	s.watches = make(map[*Watch]struct{})

	s.watcher = watcher.New()
	s.watcher.IgnoreHiddenFiles(true)
	go func() {
		if err := s.watcher.Start(500 * time.Millisecond); err != nil {
			log.Println(err)
		}
	}()

	var err error
	s.tmpDir, err = ioutil.TempDir(os.TempDir(), "fswatch")
	return err
}

func (s *Service) TerminateDaemon(ctx context.Context) error {
	if s.tmpDir != "" {
		os.RemoveAll(s.tmpDir)
	}
	s.watcher.Close()
	return nil
}
