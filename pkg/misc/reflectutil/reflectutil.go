package reflectutil

import (
	"fmt"
	"path"
	"reflect"
	"sort"
	"strings"
)

func MemberKind(o reflect.Value, field string) reflect.Kind {
	if o.Type().Kind() == reflect.Struct {
		return FieldType(o.Type(), field).Kind()
	} else {
		v := Get(o, field)
		if !v.IsValid() {
			return reflect.Invalid
		} else {
			t, _ := UnwrapType(v)
			return t.Kind()
		}
	}
}

func UnwrapType(v reflect.Value) (t reflect.Type, ok bool) {
	defer func() {
		// recover from "call of reflect.Value.Type on zero Value"
		// since I think it only happens when you have an empty
		// value in a map[]interface{}
		if r := recover(); r != nil {
			t = v.Type()
			ok = false
		}
	}()
	t = UnwrapValue(v).Type()
	ok = true
	return
}

func UnwrapValue(v reflect.Value) reflect.Value {
	// this will "unwrap" the interface value
	// so we don't get interface kind for a string
	// value in a map[]interface{}
	return reflect.ValueOf(v.Interface())
}

// FieldsTagged returns field names that have a struct tag
// including a particular key, or if value is provided it returns
// fields that include that key and value.
func FieldsTagged(t reflect.Type, key, value string) []string {
	var f []string
	for i := 0; i < t.NumField(); i++ {
		v, ok := t.Field(i).Tag.Lookup(key)
		if !ok {
			continue
		}
		if value != "" && v != value {
			continue
		}
		f = append(f, t.Field(i).Name)
	}
	return f
}

// Fields returns exported field names for Type t.
// If it is not a struct, it returns an empty slice.
func Fields(t reflect.Type) []string {
	var f []string
	for i := 0; i < t.NumField(); i++ {
		name := t.Field(i).Name
		if name[0] == strings.ToUpper(name)[0] {
			f = append(f, name)
		}
	}
	return f
}

func FieldType(t reflect.Type, f string) reflect.Type {
	field, _ := t.FieldByName(f)
	return field.Type
}

// Get returns the member value by name m of value v. Members include map keys, struct fields, and methods.
func Get(v reflect.Value, m string) reflect.Value {
	if v.Type().Kind() == reflect.Map {
		return v.MapIndex(reflect.ValueOf(m))
	} else {
		if HasMethod(v, m) {
			return v.MethodByName(m)
		} else {
			return reflect.Indirect(v).FieldByName(m)
		}
	}
}

// Methods returns the names of methods on Value v.
func Methods(v reflect.Value) []string {
	var methods []string
	for idx := 0; idx < v.NumMethod(); idx++ {
		methods = append(methods, v.Type().Method(idx).Name)
	}
	return methods
}

// HasMethod returns whether Value v has the method m.
func HasMethod(v reflect.Value, m string) bool {
	for _, meth := range Methods(v) {
		if m == meth {
			return true
		}
	}
	return false
}

// HasKey returns whether Value v has the field or map key k.
func HasKey(v reflect.Value, k string) bool {
	for _, key := range Keys(v) {
		if k == key {
			return true
		}
	}
	return false
}

// Keys returns the names of settable fields or map keys of Value v.
func Keys(v reflect.Value) []string {
	if v.Type().Kind() == reflect.Map {
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
	}
	return Fields(v.Type())
}

// Values returns a slice of Values for the values contained in Value v if it is an Array or Slice.
// It panics if v's Kind is not Array or Slice.
func Values(v reflect.Value) []reflect.Value {
	kind := v.Type().Kind()
	// TODO: maps?
	if kind != reflect.Array && kind != reflect.Slice {
		panic("Values called on value that is not an Array or Slice")
	}
	var vals []reflect.Value
	for i := 0; i < v.Len(); i++ {
		vals = append(vals, v.Index(i))
	}
	return vals
}

func TypeName(v interface{}) string {
	t := reflect.Indirect(reflect.ValueOf(v)).Type()
	return fmt.Sprintf("%s.%s", path.Base(t.PkgPath()), t.Name())
}
