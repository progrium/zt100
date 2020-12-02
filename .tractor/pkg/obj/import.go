package obj

import (
	"github.com/manifold/tractor/pkg/manifold/library"
	"path"
	"runtime"
	buujj4ug10l8pq4cqegg "worksite/pkg/obj/buujj4ug10l8pq4cqegg"
)

// GENERATED INDEX
func relPath(subpath string) string {
	_, filename, _, _ := runtime.Caller(1)
	return path.Join(path.Dir(filename), subpath)
}

func init() {
	library.Register(&buujj4ug10l8pq4cqegg.Main{}, "buujj4ug10l8pq4cqegg", relPath("buujj4ug10l8pq4cqegg/component.go"), "fas fa-user-robot")
}
