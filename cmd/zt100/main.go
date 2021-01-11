package main

import (
	"fmt"
	"log"
	"os"

	_ "github.com/progrium/zt100"

	"github.com/progrium/zt100/pkg/stdlib/net/tcp"

	"github.com/progrium/zt100/pkg/core/cmd"
	"github.com/progrium/zt100/pkg/core/daemon"
	"github.com/progrium/zt100/pkg/core/fswatch"
	"github.com/progrium/zt100/pkg/core/obj"
	"github.com/progrium/zt100/pkg/core/rpc"
	"github.com/progrium/zt100/pkg/core/service"
	"github.com/progrium/zt100/pkg/manifold"
	"github.com/progrium/zt100/pkg/manifold/image"
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
		&Bootstrapper{},
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

	fmt.Println("zt100 starting...")
	if err := daemon.Run(reg); err != nil {
		log.Fatal(err)
	}
}

type Bootstrapper struct {
	Objects *obj.Service
}

func (b *Bootstrapper) InitializeDaemon() (err error) {
	os.MkdirAll("local/uploads", 0755)

	sys := b.Objects.Root.ChildAt(0)
	if len(sys.Children()) == 0 && len(sys.Components()) == 0 {
		obj, refs, err := image.FromSnapshot(manifold.ObjectSnapshot{
			ID:   "bvbr2fug10l125mlqnfg",
			Name: "Zartan",
			Components: []manifold.ComponentSnapshot{
				{
					ID:      "bvbr2nmg10l125mlqng0",
					Name:    "zt100.Server",
					Enabled: true,
				},
				{
					ID:      "bvbr2qeg10l125mlqngg",
					Name:    "tcp.Listener",
					Enabled: true,
					Value: &tcp.Listener{
						Address: ":8080",
					},
				},
				{
					ID:      "bvbr2tmg10l125mlqnh0",
					Name:    "http.Server",
					Enabled: true,
					Refs: []manifold.SnapshotRef{
						{
							ObjectID: "bvbr2fug10l125mlqnfg",
							Path:     "2:http.Server/Listener",
							TargetID: "bvbr2fug10l125mlqnfg/1:tcp.Listener",
						},
						{
							ObjectID: "bvbr2fug10l125mlqnfg",
							Path:     "2:http.Server/Handler",
							TargetID: "bvbr2fug10l125mlqnfg/0:zt100.Server",
						},
					},
				},
			},
		})
		if err != nil {
			return err
		}
		b.Objects.Root.RemoveChild(sys)
		b.Objects.Root.AppendChild(obj)
		image.ApplyRefs(b.Objects.Root, refs)
		if err := manifold.Walk(b.Objects.Root, b.Objects.InitializeObject); err != nil {
			return err
		}
		log.Println("new object tree")
	}
	return nil
}
