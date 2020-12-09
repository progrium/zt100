package zt100

import (
	"path"
	"runtime"

	"github.com/manifold/tractor/pkg/manifold/library"
	"github.com/progrium/zt100/pkg/zt100"
)

func filepath(subpath string) string {
	_, filename, _, _ := runtime.Caller(1)
	return path.Join(path.Dir(filename), subpath)
}

func init() {
	library.Register(&zt100.Theme{}, "", filepath("pkg/zt100/misc.go"), "fas fa-palette")
	library.Register(&zt100.Prospect{}, "", filepath("pkg/zt100/prospect.go"), "fas fa-user")
	library.Register(&zt100.App{}, "", filepath("pkg/zt100/app.go"), "fas fa-browser")
	library.Register(&zt100.Page{}, "", filepath("pkg/zt100/page.go"), "fas fa-file-alt")
	library.Register(&zt100.Block{}, "", filepath("pkg/zt100/block.go"), "far fa-cube")
	library.Register(&zt100.Section{}, "", filepath("pkg/zt100/misc.go"), "far fa-cube-alt")
	library.Register(&zt100.AppLibrary{}, "", filepath("pkg/zt100/misc.go"), "fas fa-books")
	library.Register(&zt100.BlockLibrary{}, "", filepath("pkg/zt100/misc.go"), "fas fa-cubes")
	library.Register(&zt100.Server{}, "", filepath("pkg/zt100/zt100.go"), "fas fa-users")
}
