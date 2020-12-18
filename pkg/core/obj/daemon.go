package obj

import (
	"context"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/progrium/zt100/pkg/manifold"
	"github.com/progrium/zt100/pkg/manifold/comutil"
	"github.com/progrium/zt100/pkg/manifold/image"
	"github.com/progrium/zt100/pkg/misc/debouncer"
	"github.com/progrium/zt100/pkg/misc/notify"
)

func (s *Service) InitializeDaemon() (err error) {
	wd, err := os.Getwd() // TODO: override with env var
	if err != nil {
		return err
	}
	if dirExists(filepath.Join(wd, ".tractor")) {
		wd = filepath.Join(wd, ".tractor")
	}

	s.Image = image.New(wd)

	s.Root, err = s.Image.Load()
	if err != nil {
		return err
	}

	comutil.Root = s.Root

	if err := manifold.Walk(s.Root, s.InitializeObject); err != nil {
		return err
	}

	debounce := debouncer.New(2 * time.Second)
	notify.Observe(s.Root, notify.Func(func(event interface{}) error {
		// initialize any new objects
		change, ok := event.(manifold.ObjectChange)
		if ok && change.Path == "::Children" && change.New != nil {
			s.InitializeObject(change.New.(manifold.Object))
		}

		// update any ui state
		s.UI.UpdateState()

		// snapshot behind 2 sec debounce
		debounce(func() {
			start := time.Now()
			if err := s.Snapshot(); err != nil {
				log.Println("snapshot error:", err)
				return
			}
			log.Print("snapshot taken in ", time.Since(start))
		})

		return nil
	}))

	return nil
}

func (s *Service) TerminateDaemon(ctx context.Context) error {
	return s.Snapshot()
}

func dirExists(path string) bool {
	info, err := os.Stat(path)
	if os.IsNotExist(err) {
		return false
	}
	return info.IsDir()
}
