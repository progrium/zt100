package ui

import (
	"reflect"

	"github.com/progrium/zt100/pkg/manifold"
)

type LazyChilder interface {
	LazyChildren() bool
}

type Field struct {
	Type    string
	SubType *Field
	Hidden  bool
	Name    string
	Path    string
	Value   interface{}
	Enum    []string
	Min     int
	Max     uint
	Fields  []Field

	rv  reflect.Value
	obj manifold.Object
}

type Action struct {
	Name string
	Path string
	Out  reflect.Type
}

type Component struct {
	Name     string
	Index    int
	Enabled  bool
	Key      string
	Filepath string
	Icon     string
	Fields   []Field
	Actions  []Action
	Related  []string
	CustomUI []Element
}

type Object struct {
	Name        string
	ParentID    string
	Path        string
	Dir         string
	ID          string
	Icon        string
	Index       int
	Active      bool
	Components  []Component
	HasChildren bool
	Attrs       map[string]interface{}
	Children    []string
}

type Prefab struct {
	Name string `msgpack:"name"`
	ID   string `msgpack:"id"`
}

type ComponentType struct {
	Filepath string `msgpack:"filepath"`
	Name     string `msgpack:"name"`
}
