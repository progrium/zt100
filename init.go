package zt100

import (
	"path"
	"runtime"

	"github.com/manifold/tractor/pkg/manifold/library"
	"github.com/progrium/zt100/features"
	"github.com/progrium/zt100/pkg/feature"
	"github.com/progrium/zt100/pkg/zt100"
)

func filepath(subpath string) string {
	_, filename, _, _ := runtime.Caller(1)
	return path.Join(path.Dir(filename), subpath)
}

func init() {
	library.Register(&zt100.Demo{}, "", filepath("pkg/zt100/demo.go"), "fas fa-user")
	library.Register(&zt100.App{}, "", filepath("pkg/zt100/app.go"), "fas fa-browser")
	library.Register(&zt100.Page{}, "", filepath("pkg/zt100/page.go"), "fas fa-file-alt")
	library.Register(&zt100.Block{}, "", filepath("pkg/zt100/block.go"), "far fa-cube")
	library.Register(&zt100.Server{}, "", filepath("pkg/zt100/server.go"), "fas fa-users")

	feature.Register(&features.LoginFeature{})
	feature.Register(&features.ConsentFeature{})
}
