package feature

import (
	"github.com/progrium/zt100/pkg/misc/registry"
)

var (
	Registry *registry.Registry
)

func init() {
	Registry = &registry.Registry{}
}

func Register(v Feature) {
	if err := Registry.Register(v); err != nil {
		panic(err)
	}
}

type Initializer interface {
	Initialize()
}

func Load(srv interface{}) {
	r, err := registry.New()
	if err != nil {
		panic(err)
	}
	for _, e := range Registry.Entries() {
		i, ok := e.Ref.(Initializer)
		if ok {
			i.Initialize()
		}
		r.Register(e.Ref)
	}

	r.Register(srv)
	r.SelfPopulate()
}

type Feature interface {
	Flag() Flag
}

type Flag struct {
	Name     string
	Desc     string
	Subflags []Flag
}

func FlattenFlags(feats []Feature) (flags []Flag) {
	for _, feat := range feats {
		flag := feat.Flag()
		flags = append(flags, Flag{
			Name: flag.Name,
			Desc: flag.Desc,
		})
		for _, sub := range flag.Subflags {
			flags = append(flags, Flag{
				Name: sub.Name,
				Desc: sub.Desc,
			})
		}
	}
	return flags
}
