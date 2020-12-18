package make

import (
	"errors"
	"log"
	"os"
	"reflect"
	"sync"
	"time"

	"github.com/progrium/watcher"
	"github.com/spf13/afero"
)

type WatchFS interface {
	afero.Fs
	Watch(name string, watch func(watcher.Event)) (func(), error)
}

type Builder interface {
	Match(dst string) (string, bool)
	Build(fs afero.Fs, dst, src string) error
}

type Filesystem struct {
	ReadFS   WatchFS
	Builders []Builder

	built   map[string]string
	writeFs afero.Fs

	sync.Mutex
}

func (fs *Filesystem) Initialize() {
	fs.writeFs = afero.NewMemMapFs()
	fs.built = make(map[string]string)
}

func (fs *Filesystem) unioned() afero.Fs {
	return afero.NewCopyOnWriteFs(
		afero.NewReadOnlyFs(fs.ReadFS),
		fs.writeFs,
	)
}

func (fs *Filesystem) tryBuild(name string) (bool, error) {
	for _, builder := range fs.Builders {
		if src, ok := builder.Match(name); ok {
			if exists, err := afero.Exists(fs.unioned(), src); !exists || err != nil {
				return false, err
			}
			err := builder.Build(fs.unioned(), name, src)
			if err != nil {
				return false, err
			}
			log.Println(reflect.TypeOf(builder).Elem().Name(), name, "=>", src)
			fs.Lock()
			fs.built[src] = name
			fs.Unlock()
			return true, nil
		}
	}
	return false, nil
}

func (fs *Filesystem) ExportFs() afero.Fs {
	// called on referencing + object reload
	return nil
}

func (fs *Filesystem) Watch(name string, watch func(watcher.Event)) (func(), error) {
	fs.Lock()
	defer fs.Unlock()
	for src, dst := range fs.built {
		if name == dst {
			return fs.ReadFS.Watch(src, func(event watcher.Event) {
				if event.Op == watcher.Write || event.Op == watcher.Remove {
					fs.writeFs.Remove(dst)
				}
				event.Path = dst
				watch(event)
			})
		}
	}
	return fs.ReadFS.Watch(name, watch)
}

func (fs *Filesystem) Create(name string) (afero.File, error) {
	return fs.unioned().Create(name)
}

func (fs *Filesystem) Mkdir(name string, perm os.FileMode) error {
	return fs.unioned().Mkdir(name, perm)
}

func (fs *Filesystem) MkdirAll(path string, perm os.FileMode) error {
	return fs.unioned().MkdirAll(path, perm)
}

func (fs *Filesystem) Open(name string) (afero.File, error) {
	f, err := fs.unioned().Open(name)
	if errors.Is(err, os.ErrNotExist) {
		var ok bool
		ok, err = fs.tryBuild(name)
		if err != nil {
			return f, err
		}
		if ok {
			f, err = fs.unioned().Open(name)
		}
	}
	return f, err
}

func (fs *Filesystem) OpenFile(name string, flag int, perm os.FileMode) (afero.File, error) {
	f, err := fs.unioned().OpenFile(name, flag, perm)
	if errors.Is(err, os.ErrNotExist) {
		ok, err := fs.tryBuild(name)
		if err != nil {
			return f, err
		}
		if ok {
			f, err = fs.unioned().OpenFile(name, flag, perm)
		}
	}
	return f, err
}

func (fs *Filesystem) Remove(name string) error {
	return fs.unioned().Remove(name)
}

func (fs *Filesystem) RemoveAll(path string) error {
	return fs.unioned().RemoveAll(path)
}

func (fs *Filesystem) Rename(oldname, newname string) error {
	return fs.unioned().Rename(oldname, newname)
}

func (fs *Filesystem) Stat(name string) (os.FileInfo, error) {
	fi, err := fs.unioned().Stat(name)
	if errors.Is(err, os.ErrNotExist) {
		ok, err := fs.tryBuild(name)
		if err != nil {
			return fi, err
		}
		if ok {
			fi, err = fs.unioned().Stat(name)
		}
	}
	return fi, err
}

func (fs *Filesystem) Name() string {
	return fs.unioned().Name()
}

func (fs *Filesystem) Chmod(name string, mode os.FileMode) error {
	return fs.unioned().Chmod(name, mode)
}

func (fs *Filesystem) Chtimes(name string, atime time.Time, mtime time.Time) error {
	return fs.unioned().Chtimes(name, atime, mtime)
}
