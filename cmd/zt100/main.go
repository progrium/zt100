package main

import (
	"log"

	_ "github.com/progrium/zt100"

	"github.com/progrium/zt100/pkg/core/cmd"
	"github.com/progrium/zt100/pkg/core/daemon"
	"github.com/progrium/zt100/pkg/core/fswatch"
	"github.com/progrium/zt100/pkg/core/obj"
	"github.com/progrium/zt100/pkg/core/rpc"
	"github.com/progrium/zt100/pkg/core/service"
	"github.com/progrium/zt100/pkg/manifold"
	"github.com/progrium/zt100/pkg/manifold/object"
	"github.com/progrium/zt100/pkg/misc/registry"
	"github.com/progrium/zt100/pkg/stdlib"
	"github.com/progrium/zt100/pkg/ui"
	"github.com/progrium/zt100/pkg/ui/state"
)

type Initializer interface {
	Initialize()
}

func main() {
	stdlib.Load()

	ui.State = &state.Framework{}
	svcs := []interface{}{
		&fswatch.Service{},
		&daemon.Framework{},
		&rpc.Framework{},
		&cmd.Framework{},
		&service.Framework{},
		&obj.Service{},
		ui.State,
	}

	for _, svc := range svcs {
		if i, ok := svc.(Initializer); ok {
			i.Initialize()
		}
	}

	reg, err := registry.New(svcs...)
	if err != nil {
		log.Fatal(err)
	}

	object.RegistryPreloader = func(o manifold.Object) []interface{} {
		return svcs
	}

	if err := daemon.Run(reg); err != nil {
		log.Fatal(err)
	}
}
