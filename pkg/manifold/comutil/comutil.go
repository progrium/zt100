package comutil

import (
	"log"
	"reflect"

	"github.com/progrium/zt100/pkg/manifold"
)

var Root manifold.Object

func Object(v interface{}) manifold.Object {
	obj, _ := Root.FindComponent(v)
	return obj
}

func MustLookup(path string) manifold.Object {
	obj := Root.FindChild(path)
	if obj == nil {
		log.Panicf("no object at path '%s'", path)
	}
	return obj
}

func Ancestors(obj manifold.Object) []manifold.Object {
	var ancestors []manifold.Object
	for obj.Parent() != nil {
		obj = obj.Parent()
		ancestors = append(ancestors, obj)
	}
	return ancestors
}

func AncestorValue(obj manifold.Object, v interface{}) manifold.Object {
	rv := reflect.ValueOf(v)
	for _, a := range Ancestors(obj) {
		if a.ValueTo(rv) {
			return a
		}
	}
	return nil
}

func Enabled(obj manifold.Object) []manifold.Component {
	var out []manifold.Component
	for _, c := range obj.Components() {
		if c.Enabled() {
			out = append(out, c)
		}
	}
	return out
}
