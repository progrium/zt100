package zt100

import (
	_ "image/png"

	"github.com/manifold/tractor/pkg/manifold"
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

func (s *Section) Mounted(obj manifold.Object) error {
	s.object = obj
	_, com := obj.FindComponent(s)
	s.Key = com.Key()
	s.OID = obj.ID()
	return nil
}

var Template = ``
