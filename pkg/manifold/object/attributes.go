package object

import (
	"github.com/progrium/zt100/pkg/manifold"
	"github.com/progrium/zt100/pkg/misc/notify"
)

type attributes struct {
	m map[string]interface{}
	o *object
}

func (a *attributes) Has(attr string) bool {
	_, ok := a.m[attr]
	return ok
}

func (a *attributes) Get(attr string) interface{} {
	return a.m[attr]
}

func (a *attributes) Set(attr string, value interface{}) {
	prev := a.m[attr]
	if prev != value {
		a.m[attr] = value
		notify.Send(a.o, manifold.ObjectChange{
			Object: a.o,
			Path:   "::attrs:" + attr,
			Old:    prev,
			New:    value,
		})
	}
}

func (a *attributes) Del(attr string) {
	prev := a.m[attr]
	if prev != nil {
		delete(a.m, attr)
		notify.Send(a.o, manifold.ObjectChange{
			Object: a.o,
			Path:   "::attrs:" + attr,
			Old:    prev,
		})
	}
}

func (a *attributes) Snapshot() map[string]interface{} {
	// TODO: copy into new map w/ lock
	return a.m
}
