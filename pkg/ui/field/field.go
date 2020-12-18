package field

import (
	"fmt"
	"math"
	"path"
	"path/filepath"
	"reflect"
	"strconv"
	"time"

	"github.com/progrium/zt100/pkg/manifold"
	"github.com/progrium/zt100/pkg/misc/reflectutil"
	"github.com/progrium/zt100/pkg/ui"
)

type MinValuer interface {
	MinValue() int
}

type Enumer interface {
	Enum() []string
}

type field struct {
	ui.Field

	v        interface{}
	rv       reflect.Value
	parent   reflect.Value
	basepath string
	obj      manifold.Object
	tag      reflect.StructTag
}

var typeBuilders map[string]func(*field)

func init() {
	typeBuilders = map[string]func(*field){
		"invalid": func(f *field) {
			f.Type = "string"
			f.Value = "INVALID"
		},
		"unknown": func(f *field) {
			f.Type = "string"
			f.Value = "UNKNOWN"
		},
		"boolean": func(f *field) {
			f.Value = f.v
		},
		"number": func(f *field) {
			applyMinMax(&(f.Field), f.rv)
			e, ok := f.v.(Enumer)
			if ok {
				f.Enum = e.Enum()
				f.Max = uint(len(f.Enum) - 1)
				f.Min = 0
			}
			mv, ok := f.v.(MinValuer)
			if ok {
				f.Min = mv.MinValue()
			}
			f.Value = f.v

		},
		"string": func(f *field) {
			e, ok := f.v.(Enumer)
			if ok {
				f.Enum = e.Enum()
			}
			f.Value = f.v
		},
		"color": func(f *field) {
			e, ok := f.v.(Enumer)
			if ok {
				f.Enum = e.Enum()
			}
			f.Value = f.v
		},
		"password": func(f *field) {
			f.Value = f.v
		},
		"time": func(f *field) {
			f.Value = f.v.(time.Time).Format("15:04")
		},
		"date": func(f *field) {
			f.Value = f.v.(time.Time).Format("2006-01-02")
		},
		"struct": func(f *field) {
			for _, fieldname := range reflectutil.Fields(f.rv.Type()) {
				f.Fields = append(f.Fields, subField(f, fieldname))
			}
		},
		"map": func(f *field) {
			f.SubType = collectionSubtypeField(f.rv)
			if f.tag.Get("fields") != "" {
				f.SubType.Type = f.tag.Get("fields")
			}
			// unwrap in case of big untyped (interface{}) value maps
			f.rv = reflectutil.UnwrapValue(f.rv)
			for _, keyname := range reflectutil.Keys(f.rv) {
				f.Fields = append(f.Fields, subField(f, keyname))
			}
		},
		"array": func(f *field) {
			f.SubType = collectionSubtypeField(f.rv)
			if f.tag.Get("fields") != "" {
				f.SubType.Type = f.tag.Get("fields")
			}
			for idx, e := range reflectutil.Values(f.rv) {
				f.Fields = append(f.Fields, subFieldElem(f, idx, e))
			}
		},
		"checkboxes": func(f *field) {
			f.SubType = collectionSubtypeField(f.rv)
			if f.tag.Get("fields") != "" {
				f.SubType.Type = f.tag.Get("fields")
			}
			for idx, e := range reflectutil.Values(f.rv) {
				f.Fields = append(f.Fields, subFieldElem(f, idx, e))
			}
			rv := reflect.New(f.rv.Type().Elem())
			e, ok := rv.Interface().(Enumer)
			if ok {
				f.Enum = e.Enum()
			}
		},
		"reference": func(f *field) {
			if !f.rv.IsValid() {
				panic("referenceField: invalid value used")
			}
			t := f.rv.Type()
			if t.Kind() == reflect.Ptr {
				t = t.Elem()
			}
			var refPath string
			refNode, com := f.obj.Root().FindComponent(f.v)
			if refNode != nil {
				refPath = filepath.Join(refNode.Path(), com.Name())
			}
			f.Type = fmt.Sprintf("reference:%s", fmt.Sprintf("%s.%s", path.Base(t.PkgPath()), t.Name())) // TODO: stop doing this
			f.SubType = &ui.Field{
				Type: t.Name(),
			}
			f.Value = refPath
		},
	}

}

func FromComponent(com manifold.Component) (fields []ui.Field) {
	obj := com.Container()
	rv := reflect.Indirect(reflect.ValueOf(com.Pointer()))
	path := filepath.Join(obj.Path(), com.Key())
	hiddenFields := reflectutil.FieldsTagged(rv.Type(), "tractor", "hidden")
	for _, fieldname := range reflectutil.Fields(rv.Type()) {
		fields = append(fields, newField(obj, rv, strInSlice(hiddenFields, fieldname), path, fieldname))
	}
	return
}

func subtypeName(rv reflect.Value) (name string) {
	defer func() {
		// recover from "Elem of invalid type interface {}"
		// since I think it only happens when you have an empty
		// value in a map[]interface{}
		if r := recover(); r != nil {
			name = "any"
		}
	}()
	rt := rv.Type()
	kind := rt.Elem().Kind()
	name = typeFromKind(kind)
	if kind == reflect.Interface && rt.Kind() == reflect.Map {
		name = "any"
	}
	return
}

func collectionSubtypeField(rv reflect.Value) *ui.Field {
	f := &ui.Field{
		Type: subtypeName(rv),
	}
	applyMinMax(f, rv)
	e, ok := rv.Interface().(Enumer)
	if ok {
		f.Enum = e.Enum()
	}
	if f.Type == "reference" {
		// same as from typeBuilders, but assumes no value
		t := rv.Type().Elem()
		if t.Kind() == reflect.Ptr {
			t = t.Elem()
		}
		f.Type = fmt.Sprintf("reference:%s", fmt.Sprintf("%s.%s", path.Base(t.PkgPath()), t.Name())) // TODO: stop doing this
		f.SubType = &ui.Field{
			Type: t.Name(),
		}
		f.Value = ""
	}
	return f
}

func applyMinMax(f *ui.Field, rv reflect.Value) {
	switch rv.Type().Kind() {
	case reflect.Int, reflect.Int64:
		f.Min = math.MinInt64
		f.Max = math.MaxInt64
	case reflect.Int8:
		f.Min = math.MinInt8
		f.Max = math.MaxInt8
	case reflect.Int16:
		f.Min = math.MinInt16
		f.Max = math.MaxInt16
	case reflect.Int32:
		f.Min = math.MinInt32
		f.Max = math.MaxInt32
	case reflect.Uint, reflect.Uint64:
		f.Max = math.MaxUint64
	case reflect.Uint8:
		f.Max = math.MaxUint8
	case reflect.Uint16:
		f.Max = math.MaxUint16
	case reflect.Uint32:
		f.Max = math.MaxUint32
	case reflect.Float32:
		f.Max = math.MaxUint32 // math.MaxFloat32 someday
	case reflect.Float64:
		f.Max = math.MaxUint64 // math.MaxFloat64 someday
	}
}

func newField(obj manifold.Object, parent reflect.Value, hidden bool, basepath, fieldname string) ui.Field {
	return subField(&field{
		Field: ui.Field{
			Path:   basepath,
			Hidden: hidden,
		},
		rv:  parent,
		obj: obj,
	}, fieldname)
}

func subField(f *field, fieldname string) ui.Field {
	sf := &field{
		Field: ui.Field{
			Type:   typeFromKind(reflectutil.MemberKind(f.rv, fieldname)),
			Path:   filepath.Join(f.Path, fieldname),
			Hidden: f.Hidden,
		},
		parent: f.rv,
		obj:    f.obj,
		rv:     reflectutil.Get(f.rv, fieldname),
	}
	sf.Name = filepath.Base(sf.Path)
	sf.v = sf.rv.Interface()

	if f.rv.Kind() == reflect.Struct {
		structfield, _ := f.rv.Type().FieldByName(fieldname)
		sf.tag = structfield.Tag
		if structfield.Tag.Get("field") != "" {
			sf.Type = structfield.Tag.Get("field")
		}
	}
	builder, ok := typeBuilders[sf.Type]
	if !ok {
		builder = typeBuilders["unknown"]
	}
	builder(sf)
	return sf.Field
}

func subFieldElem(f *field, idx int, value reflect.Value) ui.Field {
	sf := &field{
		Field: ui.Field{
			Type: typeFromKind(value.Type().Kind()),
			Path: filepath.Join(f.Path, strconv.Itoa(idx)),
		},
		parent: f.rv,
		obj:    f.obj,
		rv:     value,
		v:      value.Interface(),
	}
	sf.Name = filepath.Base(sf.Path)
	if f.tag.Get("fields") != "" {
		sf.Type = f.tag.Get("fields")
	}
	builder, ok := typeBuilders[sf.Type]
	if !ok {
		builder = typeBuilders["unknown"]
	}
	builder(sf)
	return sf.Field
}

func typeFromKind(kind reflect.Kind) string {
	switch kind {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64,
		reflect.Float32, reflect.Float64:
		return "number"
	case reflect.Bool:
		return "boolean"
	case reflect.String:
		return "string"
	case reflect.Struct:
		return "struct"
	case reflect.Map:
		return "map"
	case reflect.Slice:
		return "array"
	case reflect.Ptr, reflect.Interface:
		// reflect.Interface does not always mean "reference",
		// but that is an edge case dealing with maps and (possibly)
		// arrays that use interface{}, which we deal with elsewhere
		return "reference"
	case reflect.Invalid:
		return "invalid"
	default:
		return "unknown"
	}
}

func strInSlice(strs []string, str string) bool {
	for _, s := range strs {
		if s == str {
			return true
		}
	}
	return false
}
