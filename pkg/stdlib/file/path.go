package file

import (
	"log"
	"os"
	"path"
	"path/filepath"
	"strings"
	"time"

	"github.com/progrium/watcher"
	"github.com/progrium/zt100/pkg/core/fswatch"
	"github.com/progrium/zt100/pkg/manifold"
	"github.com/progrium/zt100/pkg/misc/notify"
	"github.com/spf13/afero"
)

type Path struct {
	Filepath   string
	SyncName   bool
	SyncDelete bool

	Watcher *fswatch.Service `tractor:"hidden"`

	fs afero.Fs
}

func (c *Path) Mounted(obj manifold.Object) error {
	notify.Observe(obj, notify.Func(func(event interface{}) error {
		change, ok := event.(manifold.ObjectChange)
		if !ok || change.Path != "::Name" || obj.ID() != change.Object.ID() {
			return nil
		}
		newName := change.New.(string)
		base := path.Base(c.Filepath)
		dir := path.Dir(c.Filepath)
		if c.SyncName && base != newName {
			newPath := path.Join(dir, newName)
			if err := os.Rename(c.Filepath, newPath); err != nil {
				log.Println(err)
			}
			c.Filepath = newPath
		}
		return nil
	}))
	return nil
}

func (c *Path) exists() bool {
	if c.Filepath == "" {
		return false
	}
	if _, err := os.Stat(c.Filepath); os.IsNotExist(err) {
		return false
	}
	return true
}

func (c *Path) ComponentEnable() {
	c.fs = nil
	if c.exists() {
		c.fs = afero.NewBasePathFs(afero.NewOsFs(), c.Filepath)
	}
}

func (c *Path) Watch(name string, watch func(watcher.Event)) (func(), error) {
	watcher_ := c.Watcher
	path := c.Filepath
	if name != "." {
		path = filepath.Join(path, name)
	}

	w := &fswatch.Watch{
		Path: path,
		Handler: func(event watcher.Event) {
			event.Path = strings.TrimPrefix(event.Path, c.Filepath)
			watch(event)
		},
	}
	if err := watcher_.Watch(w); err != nil {
		return nil, err
	}

	return func() {
		watcher_.Unwatch(w)
	}, nil
}

// func (c *Path) Open(name string) (http.File, error) {
// 	return http.Dir(c.Filepath).Open(name)
// }

func (c *Path) Create(name string) (afero.File, error) {
	return c.fs.Create(name)
}

func (c *Path) Mkdir(name string, perm os.FileMode) error {
	return c.fs.Mkdir(name, perm)
}

func (c *Path) MkdirAll(path string, perm os.FileMode) error {
	return c.fs.MkdirAll(path, perm)
}

func (c *Path) Open(name string) (afero.File, error) {
	return c.fs.Open(name)
}

func (c *Path) OpenFile(name string, flag int, perm os.FileMode) (afero.File, error) {
	return c.fs.OpenFile(name, flag, perm)
}

func (c *Path) Remove(name string) error {
	return c.fs.Remove(name)
}

func (c *Path) RemoveAll(path string) error {
	return c.fs.RemoveAll(path)
}

func (c *Path) Rename(oldname, newname string) error {
	return c.fs.Rename(oldname, newname)
}

func (c *Path) Stat(name string) (os.FileInfo, error) {
	return c.fs.Stat(name)
}

func (c *Path) Name() string {
	return c.fs.Name()
}

func (c *Path) Chmod(name string, mode os.FileMode) error {
	return c.fs.Chmod(name, mode)
}

func (c *Path) Chtimes(name string, atime time.Time, mtime time.Time) error {
	return c.fs.Chtimes(name, atime, mtime)
}
