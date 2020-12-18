package stdlib

import (
	"path"
	"runtime"

	"github.com/progrium/zt100/pkg/stdlib/file"
	"github.com/progrium/zt100/pkg/stdlib/file/make"
	"github.com/progrium/zt100/pkg/stdlib/http"
	"github.com/progrium/zt100/pkg/stdlib/net/tcp"

	"github.com/progrium/zt100/pkg/manifold/library"
)

func filepath(subpath string) string {
	_, filename, _, _ := runtime.Caller(1)
	return path.Join(path.Dir(filename), subpath)
}

func Load() {

	// file
	library.Register(&file.Local{}, "", filepath("file/local.go"), "fas fa-file")
	library.Register(&file.Path{}, "", filepath("file/path.go"), "fas fa-folder")
	library.Register(&file.Explorer{}, "", filepath("file/explorer.go"), "fas fa-folder-tree")
	library.Register(&file.Reference{}, "", filepath("file/reference.go"), "fas fa-file")
	library.Register(&file.UnionFS{}, "", filepath("file/union.go"), "fas fa-folders")

	// file/make
	library.Register(&make.Filesystem{}, "", filepath("file/make/filesystem.go"), "fas fa-folder-tree")
	library.Register(&make.UpperHTML{}, "", filepath("file/make/upper_html.go"), "fas fa-file-export")
	library.Register(&make.JSX{}, "", filepath("file/make/jsx.go"), "fas fa-file-export")

	// http
	library.Register(&http.SingleUserBasicAuth{}, "", filepath("http/basicauth.go"), "fas fa-id-card")
	library.Register(&http.FileServer{}, "", filepath("http/fileserver.go"), "fas fa-file-download")
	library.Register(&http.Logger{}, "", filepath("http/logger.go"), "fas fa-monitor-heart-rate")
	library.Register(&http.Mux{}, "", filepath("http/mux.go"), "fas fa-layer-group")
	library.Register(&http.Server{}, "", filepath("http/server.go"), "fas fa-server")
	library.Register(&http.TemplateRenderer{}, "", filepath("http/templaterenderer.go"), "fas fa-file-code")

	// net
	library.Register(&tcp.Listener{}, "", filepath("net/tcp/listener.go"), "fas fa-ethernet")

}
