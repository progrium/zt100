package cmd

import (
	"fmt"
	"reflect"
	"sync"
)

type Definition struct {
	ID       string
	Label    string `json:",omitempty"`
	Category string `json:",omitempty"`
	Desc     string `json:",omitempty"`
	Run      interface{}
}

type Registry struct {
	cmds map[string]Definition
	sync.Mutex
}

func (r *Registry) Get(cmdID string) Definition {
	r.Lock()
	defer r.Unlock()
	return r.cmds[cmdID]
}

func (r *Registry) Register(def Definition) {
	r.Lock()
	defer r.Unlock()
	// TODO: assert valid Run
	r.cmds[def.ID] = def
}

func (r *Registry) Execute(cmdID string, params map[string]interface{}) (result interface{}, err error) {
	cmd := r.Get(cmdID)
	if cmd.ID == "" {
		err = fmt.Errorf("cmd not found: %s", cmdID)
		return
	}

	fn := reflect.ValueOf(cmd.Run)
	args := []reflect.Value{buildParams(fn.Type().In(0), params)}

	return parseReturn(fn.Call(args))
}

func buildParams(t reflect.Type, v map[string]interface{}) reflect.Value {
	rp := reflect.New(t)
	for k, vv := range v {
		_, ok := rp.Elem().Type().FieldByName(k)
		if ok {
			f := rp.Elem().FieldByName(k)
			ft := f.Type()
			if vv == nil {
				f.Set(reflect.Zero(ft))
				continue
			}
			if f.Kind() == reflect.Ptr {
				rpv := reflect.New(ft.Elem())
				rpv.Elem().Set(reflect.ValueOf(vv).Convert(ft.Elem()))
				f.Set(rpv)
			} else {
				f.Set(reflect.ValueOf(vv).Convert(ft))
			}
		}
	}
	return rp.Elem()
}

func parseReturn(ret []reflect.Value) (interface{}, error) {
	switch len(ret) {
	case 0:
		return nil, nil
	case 2:
		if ret[1].IsNil() {
			return ret[0].Interface(), nil
		}
		return ret[0].Interface(), ret[1].Interface().(error)
	case 1:
		err, ok := ret[0].Interface().(error)
		if ok {
			return nil, err
		}
		return ret[0].Interface(), nil
	default:
		panic("unexpected return value count")
	}
}
