// Use `tractor run` to start
package main

import (
	_ "worksite/pkg/obj"
	_ "worksite/pkg/usr"

	_ "github.com/progrium/zt100"

	"github.com/manifold/tractor/pkg/worksite"
)

func main() {
	worksite.Run()
}
