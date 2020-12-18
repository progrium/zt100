package image

import (
	"path"

	"github.com/progrium/zt100/pkg/manifold"
	"github.com/progrium/zt100/pkg/manifold/image/gen"
	"github.com/spf13/afero"
)

func (i *Image) DestroyObjectPackage(obj manifold.Object) error {
	i.pkgFs = afero.NewBasePathFs(i.fs, PackageDir)
	if err := i.pkgFs.RemoveAll(path.Join(ObjectDir, obj.ID())); err != nil {
		return err
	}
	return i.IndexObjectPackages()
}

func (i *Image) CreateObjectPackage(obj manifold.Object) error {
	i.pkgFs = afero.NewBasePathFs(i.fs, PackageDir)

	dir := path.Join(ObjectDir, obj.ID())
	if err := i.pkgFs.MkdirAll(dir, 0755); err != nil {
		return err
	}

	filepath := path.Join(dir, "component.go")
	if ok, _ := afero.Exists(i.pkgFs, filepath); ok {
		return i.IndexObjectPackages()
	}

	src := "package object\n\ntype Main struct {\n\n}\n"
	if err := afero.WriteFile(i.pkgFs, filepath, []byte(src), 0644); err != nil {
		return err
	}

	return i.IndexObjectPackages()
}

func (i *Image) IndexObjectPackages() error {
	i.pkgFs = afero.NewBasePathFs(i.fs, PackageDir)

	if err := i.pkgFs.MkdirAll(ObjectDir, 0755); err != nil {
		return err
	}

	objs := []string{}
	fi, err := afero.ReadDir(i.pkgFs, ObjectDir)
	if err != nil {
		return err
	}
	for _, info := range fi {
		if info.IsDir() {
			objs = append(objs, info.Name())
		}
	}
	src, err := gen.ObjectPackageIndex(objs)
	if err != nil {
		return err
	}

	return afero.WriteFile(i.pkgFs, path.Join(ObjectDir, "import.go"), src, 0644)
}
