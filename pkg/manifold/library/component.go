package library

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path"
	"reflect"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/mitchellh/mapstructure"
	"github.com/progrium/zt100/pkg/manifold"
	"github.com/progrium/zt100/pkg/misc/jsonpointer"
	"github.com/progrium/zt100/pkg/misc/notify"
	"github.com/progrium/zt100/pkg/ui"
	"github.com/rs/xid"
)

type Initializer interface {
	Initialize()
}

type Iconer interface {
	ComponentIcon() string
}

type component struct {
	object  manifold.Object
	name    string
	id      string
	enabled bool
	icon    string
	value   interface{}
	typed   bool
	loaded  bool // if Reload has been called once
}

type ComponentEnabler interface {
	ComponentEnable()
}
type ComponentDisabler interface {
	ComponentDisable()
}

type SetFieldObserver interface {
	ComponentFieldSet(path string, value interface{})
}

type ChildProvider interface {
	ChildNodes() []manifold.Object
}

func newComponent(name string, value interface{}, id string) *component {
	if id == "" {
		id = xid.New().String()
	}
	var typed bool
	if _, ok := value.(map[string]interface{}); !ok {
		typed = true
		if i, ok := value.(Initializer); ok {
			i.Initialize()
		}
	}
	rc := Lookup(name)
	if rc == nil {
		log.Panicf("could not find registered component: %s", name)
	}
	return &component{
		name:    name,
		enabled: false,
		value:   value,
		id:      id,
		typed:   typed,
		icon:    rc.Icon,
	}
}

func NewComponent(name string, value interface{}, id string) manifold.Component {
	return newComponent(name, value, id)
}

func (c *component) GetField(path string) (interface{}, reflect.Type, error) {
	// TODO: check if field exists
	v := jsonpointer.Reflect(c.Pointer(), path)
	// if v != nil {
	// 	return v, reflect.TypeOf(v), nil
	// }
	// return nil, c.FieldType(path), nil
	//log.Println("GetField:", reflect.TypeOf(v), c.FieldType(path))
	return v, c.FieldType(path), nil
}

func (c *component) SetField(path string, value interface{}) error {
	old, t, _ := c.GetField(path)
	if t.Kind() != reflect.Slice && t.Kind() != reflect.Map && old == value {
		return nil
	}
	// when you have a custom string type for enums
	if t.Kind() == reflect.String && t.Name() != "string" {
		rv := reflect.ValueOf(value)
		value = rv.Convert(t).Interface()
	}
	if err := jsonpointer.SetReflect(c.value, path, value); err != nil {
		return err
	}
	if sfo, ok := c.value.(SetFieldObserver); ok {
		sfo.ComponentFieldSet(path, value)
	}
	notify.Send(c.object, manifold.ObjectChange{
		Object: c.object,
		Path:   fmt.Sprintf("%s/%s", c.name, path),
		Old:    old,
		New:    value,
	})
	return nil
}

func (c *component) FieldType(path string) reflect.Type {
	parts := strings.Split(path, "/")
	rt := reflect.TypeOf(c.Pointer())
	for _, part := range parts {
		switch rt.Kind() {
		case reflect.Slice:
			rt = rt.Elem()
		case reflect.Struct:
			field, _ := rt.FieldByName(part)
			rt = field.Type
		case reflect.Ptr:
			// TODO: might not be a struct!
			field, _ := rt.Elem().FieldByName(part)
			rt = field.Type
		case reflect.Map:
			// TODO: eventually we can't assume string keys!
			rt = reflect.TypeOf("")
		default:
			panic("unhandled type: " + rt.String())
		}
	}
	return rt
}

func (c *component) CallMethod(path string, args []interface{}, reply interface{}) error {
	// TODO: support methods on sub paths / data structures
	rval := reflect.ValueOf(c.Pointer())
	method := rval.MethodByName(path)
	var params []reflect.Value
	for _, arg := range args {
		params = append(params, reflect.ValueOf(arg))
	}
	retVals := method.Call(params)
	errorInterface := reflect.TypeOf((*error)(nil)).Elem()
	// assuming up to 2 return values, one being an error
	rreply := reflect.ValueOf(reply)
	var errVal error
	for _, v := range retVals {
		if v.Type().Implements(errorInterface) {
			if !v.IsNil() {
				errVal = v.Interface().(error)
			}
		} else {
			if reply != nil {
				rreply.Elem().Set(v)
			}
		}
	}
	return errVal
}

func (c *component) Index() int {
	for idx, com := range c.object.Components() {
		if com == c {
			return idx
		}
	}
	return 0
}

func (c *component) Key() string {
	return fmt.Sprintf("%d:%s", c.Index(), c.Name())
}

func (c *component) SetIndex(idx int) {
	if idx == -1 {
		idx = len(c.object.Components()) - 1
	}
	old := c.Index()
	if old == idx {
		return
	}
	c.object.RemoveComponent(c)
	c.object.InsertComponentAt(idx, c)
	notify.Send(c.object, manifold.ObjectChange{
		Object: c.object,
		Path:   fmt.Sprintf("%s/::Index", c.name),
		Old:    old,
		New:    idx,
	})
}

func (c *component) Name() string {
	return c.name
}

func (c *component) ID() string {
	return c.id
}

func (c *component) Enabled() bool {
	return c.enabled
}

func (c *component) SetEnabled(enable bool) {
	old := c.enabled
	if old == enable {
		return
	}
	c.enabled = enable
	if enable {
		if e, ok := c.Pointer().(ComponentEnabler); ok {
			e.ComponentEnable()
		}
	} else {
		if e, ok := c.Pointer().(ComponentDisabler); ok {
			e.ComponentDisable()
		}
	}
	notify.Send(c.object, manifold.ObjectChange{
		Object: c.object,
		Path:   fmt.Sprintf("%s/::Enabled", c.name),
		Old:    old,
		New:    enable,
	})

}

func (c *component) Container() manifold.Object {
	return c.object
}

func (c *component) SetContainer(obj manifold.Object) {
	c.object = obj
}

// TODO: rename to Value()?
func (c *component) Pointer() interface{} {
	if !c.typed {
		c.value = typedComponentValue(c.value, c.name, c.id)
		c.typed = true
		if i, ok := c.value.(Initializer); ok {
			i.Initialize()
		}
	}
	return c.value
}

func (c *component) Type() reflect.Type {
	return reflect.TypeOf(c.Pointer())
}

func (c *component) Reload() error {
	if c.loaded && c.enabled {
		if e, ok := c.Pointer().(ComponentDisabler); ok {
			e.ComponentDisable()
		}
	}
	if e, ok := c.Pointer().(ComponentEnabler); ok {
		e.ComponentEnable()
	}
	// TODO: maybe combine LazyChilder+ChildProvider interfaces
	lazyChildren := false
	if lc, ok := c.Pointer().(ui.LazyChilder); ok {
		lazyChildren = lc.LazyChildren()
	}
	if len(c.object.Children()) == 0 && !lazyChildren {
		if cp, ok := c.Pointer().(ChildProvider); ok {
			for _, obj := range cp.ChildNodes() {
				c.object.AppendChild(obj)
			}
		}
	}
	c.SetEnabled(true)
	c.loaded = true
	return nil
}

func (c *component) Icon() string {
	if i, ok := c.Pointer().(Iconer); ok {
		return i.ComponentIcon()
	}
	return c.icon
}

// TODO
func (c *component) Fields() {}

// TODO
func (c *component) Methods() {}

// TODO
func (c *component) RelatedPrefabs() {}

func (c *component) Snapshot() manifold.ComponentSnapshot {
	if !c.typed {
		panic("snapshot before component value is typed")
	}
	com := manifold.ComponentSnapshot{
		Name:    c.name,
		ID:      c.id,
		Enabled: c.enabled,
		Value:   c.value,
	}
	if c.object != nil {
		com.ObjectID = c.object.ID()
		com.Value, com.Refs = extractRefs(c.object, c.Key(), com.Value)
	}
	return com
}

func extractRefs(obj manifold.Object, basePath string, v interface{}) (out interface{}, refs []manifold.SnapshotRef) {
	if obj.Root() == nil {
		return
	}
	if _, ok := v.(json.Marshaler); ok {
		out = v
		return
	}
	sub := make(map[string]interface{})
	rv := reflect.ValueOf(v)
	for _, field := range keys(rv) {
		fv := prop(rv, field)
		ft := fv.Type()
		fieldPath := path.Join(basePath, field)
		var subrefs []manifold.SnapshotRef
		switch ft.Kind() {
		case reflect.Func:
			return
		case reflect.Slice:
			sub[field], subrefs = extractRefsSlice(obj, fieldPath, fv.Interface())
			refs = append(refs, subrefs...)
		case reflect.Struct, reflect.Map:
			sub[field], subrefs = extractRefs(obj, fieldPath, fv.Interface())
			refs = append(refs, subrefs...)
		case reflect.Ptr, reflect.Interface:
			if fv.IsNil() {
				continue
			}
			target, com := obj.Root().FindComponent(fv.Interface())
			if target != nil {
				refs = append(refs, manifold.SnapshotRef{
					ObjectID: obj.ID(),
					Path:     fieldPath,
					TargetID: fmt.Sprintf("%s/%s", target.ID(), com.Key()),
				})
				sub[field] = nil
			}
		default:
			sub[field] = fv.Interface()
		}
	}
	out = sub
	return
}

func extractRefsSlice(obj manifold.Object, basePath string, v interface{}) (out []interface{}, refs []manifold.SnapshotRef) {
	if obj.Root() == nil {
		return
	}
	rv := reflect.ValueOf(v)
	if rv.Kind() != reflect.Slice {
		return
	}
	for i := 0; i < rv.Len(); i++ {
		field := rv.Index(i)
		ft := field.Type()
		fieldPath := path.Join(basePath, strconv.Itoa(i))
		var subrefs []manifold.SnapshotRef
		var vv interface{}
		switch ft.Kind() {
		case reflect.Slice:
			vv, subrefs = extractRefsSlice(obj, fieldPath, field.Interface())
			out = append(out, vv)
			refs = append(refs, subrefs...)
		case reflect.Struct, reflect.Map:
			vv, subrefs = extractRefs(obj, fieldPath, field.Interface())
			out = append(out, vv)
			refs = append(refs, subrefs...)
		case reflect.Ptr, reflect.Interface:
			if field.IsNil() {
				continue
			}
			target, com := obj.Root().FindComponent(field.Interface())
			if target != nil {
				refs = append(refs, manifold.SnapshotRef{
					ObjectID: obj.ID(),
					Path:     fieldPath,
					TargetID: fmt.Sprintf("%s/%s", target.ID(), com.Key()),
				})
				out = append(out, nil)
			}
		default:
			out = append(out, field.Interface())
		}
	}
	return
}

func typedComponentValue(value interface{}, name, id string) interface{} {
	var typedValue interface{}
	if rc := Lookup(name); rc != nil {
		typedValue = rc.NewValue()
	}
	// if id != "" {
	// 	if rc := LookupID(id); rc != nil {
	// 		typedValue = rc.NewValue()
	// 	}
	// }
	if typedValue == nil {
		panic("unable to find registered component: " + name)
	}
	decoder, err := mapstructure.NewDecoder(&mapstructure.DecoderConfig{
		Metadata:   nil,
		DecodeHook: mapstructure.ComposeDecodeHookFunc(mapstructure.StringToTimeHookFunc(time.RFC3339)),
		Result:     typedValue,
	})
	if err != nil {
		panic(err)
	}
	if err := decoder.Decode(value); err == nil {
		return typedValue
	} else {
		panic(err)
	}
}

func keys(v reflect.Value) []string {
	switch v.Type().Kind() {
	case reflect.Map:
		var keys []string
		for _, key := range v.MapKeys() {
			k, ok := key.Interface().(string)
			if !ok {
				continue
			}
			keys = append(keys, k)
		}
		sort.Sort(sort.StringSlice(keys))
		return keys
	case reflect.Struct:
		t := v.Type()
		var f []string
		for i := 0; i < t.NumField(); i++ {
			name := t.Field(i).Name
			// first letter capitalized means exported
			if name[0] == strings.ToUpper(name)[0] {
				f = append(f, name)
			}
		}
		return f
	case reflect.Slice, reflect.Array:
		var k []string
		for n := 0; n < v.Len(); n++ {
			k = append(k, strconv.Itoa(n))
		}
		return k
	case reflect.Ptr:
		if !v.IsNil() {
			return keys(v.Elem())
		}
		return []string{}
	case reflect.String, reflect.Bool, reflect.Float64, reflect.Float32, reflect.Interface:
		return []string{}
	default:
		fmt.Fprintf(os.Stderr, "unexpected type: %s\n", v.Type().Kind())
		return []string{}
	}
}

func prop(robj reflect.Value, key string) reflect.Value {
	rtyp := robj.Type()
	switch rtyp.Kind() {
	case reflect.Slice, reflect.Array:
		idx, err := strconv.Atoi(key)
		if err != nil {
			panic("non-numeric index given for slice")
		}
		rval := robj.Index(idx)
		if rval.IsValid() {
			return reflect.ValueOf(rval.Interface())
		}
	case reflect.Ptr:
		return prop(robj.Elem(), key)
	case reflect.Map:
		rval := robj.MapIndex(reflect.ValueOf(key))
		if rval.IsValid() {
			return reflect.ValueOf(rval.Interface())
		}
	case reflect.Struct:
		rval := robj.FieldByName(key)
		if rval.IsValid() {
			return rval
		}
		for i := 0; i < rtyp.NumField(); i++ {
			field := rtyp.Field(i)
			tag := strings.Split(field.Tag.Get("json"), ",")
			if tag[0] == key || field.Name == key {
				return robj.FieldByName(field.Name)
			}
		}
		panic("struct field not found: " + key)
	}
	//spew.Dump(robj, key)
	panic("unexpected kind: " + rtyp.Kind().String())
}
