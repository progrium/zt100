package zt100

import (
	_ "image/png"

	"github.com/manifold/tractor/pkg/manifold"
	"github.com/manifold/tractor/pkg/misc/notify"
)

type AppLibrary struct {
}

type BlockLibrary struct {
}

type Theme struct {
	Name string
}

type MenuItem struct {
	Title string
	Page  string
}

type Section struct {
	Block     *Block
	Overrides map[string]string
	Key       string `tractor:"hidden"`
	OID       string `tractor:"hidden"`
	object    manifold.Object
}

func (s *Section) Initialize() {
	if s.Overrides == nil {
		s.Overrides = make(map[string]string)
	}
}

func (s *Section) Mounted(obj manifold.Object) error {
	s.object = obj
	_, com := obj.FindComponent(s)
	s.Key = com.ID()
	s.OID = obj.ID()
	notify.Observe(obj, notify.Func(func(event interface{}) error {
		_, com := obj.FindComponent(s)
		if com != nil {
			s.Key = com.ID()
			s.OID = obj.ID()
		}
		return nil
	}))
	return nil
}
