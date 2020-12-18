package image

import (
	"fmt"
	"path"

	"github.com/progrium/zt100/pkg/manifold/image/gen"
	"github.com/spf13/afero"
)

func (i *Image) UserPackagePath(name string) string {
	return path.Join(i.filepath, PackageDir, UserPkgDir, name)
}

func (i *Image) DestroyUserPackage(name string) error {
	i.pkgFs = afero.NewBasePathFs(i.fs, PackageDir)
	if err := i.pkgFs.RemoveAll(path.Join(UserPkgDir, name)); err != nil {
		return err
	}
	return i.IndexUserPackages()
}

func (i *Image) CreateUserPackage(name string) error {
	i.pkgFs = afero.NewBasePathFs(i.fs, PackageDir)

	dir := path.Join(UserPkgDir, name)
	if err := i.pkgFs.MkdirAll(dir, 0755); err != nil {
		return err
	}

	filepath := path.Join(dir, "component.go")
	if ok, _ := afero.Exists(i.pkgFs, filepath); ok {
		return nil
	}

	src := fmt.Sprintf("package %s\n\ntype Main struct {\n\n}\n", name)
	if err := afero.WriteFile(i.pkgFs, filepath, []byte(src), 0644); err != nil {
		return err
	}

	return i.IndexUserPackages()
}

func (i *Image) IndexUserPackages() error {
	i.pkgFs = afero.NewBasePathFs(i.fs, PackageDir)

	if err := i.pkgFs.MkdirAll(UserPkgDir, 0755); err != nil {
		return err
	}

	pkgs := []string{}
	fi, err := afero.ReadDir(i.pkgFs, UserPkgDir)
	if err != nil {
		return err
	}
	for _, info := range fi {
		if info.IsDir() {
			pkgs = append(pkgs, info.Name())
		}
	}
	src, err := gen.UserPackageIndex(pkgs)
	if err != nil {
		return err
	}

	return afero.WriteFile(i.pkgFs, path.Join(UserPkgDir, "import.go"), src, 0644)
}
