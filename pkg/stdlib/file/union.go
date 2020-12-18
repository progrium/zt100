package file

import (
	"errors"
	"fmt"

	"github.com/progrium/watcher"
	"github.com/progrium/zt100/pkg/manifold"
	"github.com/progrium/zt100/pkg/misc/vfs"
	"github.com/spf13/afero"
)

type WatchFS interface {
	afero.Fs
	Watch(name string, watch func(watcher.Event)) (func(), error)
}

type UnionFS struct {
	vfs.UnionFS

	Base  afero.Fs
	Layer afero.Fs
}

func (c *UnionFS) Open(name string) (afero.File, error) {
	f, err := c.UnionFS.Open(name)
	if f == nil {
		return nil, fmt.Errorf("no file")
	}
	return f, err
}

func (c *UnionFS) Mounted(obj manifold.Object) error {
	c.UnionFS.Layers = []afero.Fs{c.Base, c.Layer}
	return nil
}

func (c *UnionFS) ComponentEnable() {
	c.UnionFS.Layers = []afero.Fs{c.Base, c.Layer}
}

func (c *UnionFS) Watch(name string, watch func(watcher.Event)) (func(), error) {
	if ok, _ := afero.Exists(c.Layer, name); ok {
		if wfs, ok := c.Layer.(WatchFS); ok {
			return wfs.Watch(name, watch)
		}
	}
	if wfs, ok := c.Base.(WatchFS); ok {
		return wfs.Watch(name, watch)
	}
	return nil, errors.New("no watchable filesystems in union")
}
