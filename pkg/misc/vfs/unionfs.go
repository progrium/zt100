package vfs

import (
	"os"
	"syscall"
	"time"

	"github.com/spf13/afero"
)

type UnionFS struct {
	Layers []afero.Fs
}

func NewUnionFS(layers ...afero.Fs) *UnionFS {
	return &UnionFS{
		Layers: layers,
	}
}

func isNotExist(err error) bool {
	if e, ok := err.(*os.PathError); ok {
		err = e.Err
	}
	if err == os.ErrNotExist || err == syscall.ENOENT || err == syscall.ENOTDIR {
		return true
	}
	return false
}

func (fs *UnionFS) Name() string {
	return "UnionFS"
}

func (u *UnionFS) isBaseFile(name string) (bool, error) {
	for i := len(u.Layers) - 1; i >= 0; i-- {
		if i > 0 {
			if _, err := u.Layers[i].Stat(name); err == nil {
				return false, nil
			}
			continue
		}
		_, err := u.Layers[i].Stat(name)
		if err != nil {
			if oerr, ok := err.(*os.PathError); ok {
				if oerr.Err == os.ErrNotExist || oerr.Err == syscall.ENOENT || oerr.Err == syscall.ENOTDIR {
					return false, nil
				}
			}
			if err == syscall.ENOENT {
				return false, nil
			}
		}
		return true, err
	}
	return false, nil
}

func (u *UnionFS) Open(name string) (afero.File, error) {
	// Since the overlay overrides the base we check that first
	b, err := u.isBaseFile(name)
	if err != nil {
		return nil, err
	}

	// If overlay doesn't exist, return the base (base state irrelevant)
	if b {
		return u.Layers[0].Open(name)
	}

	// If overlay is a file, return it (base state irrelevant)
	for i := len(u.Layers) - 1; i >= 1; i-- {
		dir, err := afero.IsDir(u.Layers[i], name)
		if err != nil {
			if isNotExist(err) {
				continue
			}
			return nil, err
		}
		if !dir {
			return u.Layers[i].Open(name)
		}
	}

	// Overlay is a directory, base state now matters.
	// Base state has 3 states to check but 2 outcomes:
	// A. It's a file or non-readable in the base (return just the overlay)
	// B. It's an accessible directory in the base (return a UnionFile)

	// If base is file or nonreadable, return overlay
	dir, err := afero.IsDir(u.Layers[0], name)
	if !dir || err != nil {
		for i := len(u.Layers) - 1; i >= 1; i-- {
			f, err := u.Layers[i].Open(name)
			if err != nil {
				if isNotExist(err) {
					continue
				}
				return nil, err
			}
			return f, nil
		}
		return u.Layers[0].Open(name)
	}

	// Should be directories for this name in all layers
	// Return union file (if opens are without error)
	base := u.Layers[0]
	file, err := base.Open(name)
	if err != nil {
		return nil, err
	}
	for i := 1; i < len(u.Layers); i++ {
		lfile, err := u.Layers[i].Open(name)
		if err != nil {
			return nil, err
		}
		file = &afero.UnionFile{Base: file, Layer: lfile}
	}
	return file, nil
}

func (u *UnionFS) OpenFile(name string, flag int, perm os.FileMode) (afero.File, error) {
	b, err := u.isBaseFile(name)
	if err != nil {
		return nil, err
	}

	if flag&(os.O_WRONLY|os.O_RDWR|os.O_APPEND|os.O_CREATE|os.O_TRUNC) != 0 {
		return nil, syscall.EPERM
	}

	if b {
		return u.Layers[0].OpenFile(name, flag, perm)
	}

	for i := len(u.Layers) - 1; i >= 1; i-- {
		f, err := u.Layers[i].OpenFile(name, flag, perm)
		if err != nil {
			if isNotExist(err) {
				continue
			}
			return nil, err
		}
		return f, nil
	}

	return u.Layers[0].OpenFile(name, flag, perm)
}

func (u *UnionFS) Stat(name string) (fi os.FileInfo, err error) {
	for i := len(u.Layers) - 1; i >= 0; i-- {
		fi, err = u.Layers[i].Stat(name)
		if err != nil {
			if isNotExist(err) && i > 0 {
				continue
			}
			return nil, err
		}
		return fi, nil
	}
	return
}

func (fs *UnionFS) Create(name string) (afero.File, error) {
	return nil, syscall.EPERM
}

func (fs *UnionFS) Mkdir(name string, perm os.FileMode) error {
	return syscall.EPERM
}

func (fs *UnionFS) MkdirAll(path string, perm os.FileMode) error {
	return syscall.EPERM
}

func (fs *UnionFS) Remove(name string) error {
	return syscall.EPERM
}

func (fs *UnionFS) RemoveAll(path string) error {
	return syscall.EPERM
}

func (fs *UnionFS) Rename(oldname, newname string) error {
	return syscall.EPERM
}

func (fs *UnionFS) Chmod(name string, mode os.FileMode) error {
	return syscall.EPERM
}

func (fs *UnionFS) Chtimes(name string, atime time.Time, mtime time.Time) error {
	return syscall.EPERM
}
