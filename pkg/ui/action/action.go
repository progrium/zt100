package action

import (
	"path/filepath"
	"reflect"

	"github.com/progrium/zt100/pkg/manifold"
	"github.com/progrium/zt100/pkg/misc/reflectutil"
	"github.com/progrium/zt100/pkg/ui"
)

func FromMethods(com manifold.Component) (actions []ui.Action) {
	obj := com.Container()
	rv := reflect.ValueOf(com.Pointer())
	path := filepath.Join(obj.Path(), com.Key())
	for _, methodname := range reflectutil.Methods(rv) {
		m, _ := rv.Type().MethodByName(methodname)
		var out reflect.Type = nil
		if m.Type.NumOut() > 0 {
			out = m.Type.Out(0)
		}
		actions = append(actions, ui.Action{
			Name: methodname,
			Path: filepath.Join(path, methodname),
			Out:  out,
		})
	}
	return
}
