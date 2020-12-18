package obj

import (
	"fmt"
	"log"
	"os"
	"path"
	"path/filepath"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/progrium/zt100/pkg/core/cmd"
	"github.com/progrium/zt100/pkg/manifold"
	"github.com/progrium/zt100/pkg/manifold/library"
	"github.com/progrium/zt100/pkg/manifold/object"
	"github.com/progrium/zt100/pkg/stdlib/file"
	"github.com/progrium/zt100/pkg/ui"

	L "github.com/progrium/zt100/pkg/misc/logging"
)

func (s *Service) ContributeCommands(cmds *cmd.Registry) {
	cmds.Register(cmd.Definition{
		ID:       "debug.exit",
		Label:    "Exit",
		Category: "Debug",
		Desc:     "Exit",
		Run: func(params interface{}) {
			if err := s.Snapshot(); err != nil {
				log.Fatal(err)
			}
			os.Exit(0)
		},
	})

	cmds.Register(cmd.Definition{
		ID:       "object.select",
		Label:    "Select",
		Category: "Objects",
		Desc:     "Selects an object",
		Run: func(params struct {
			ID   string
			Path string
		}) {
			if params.ID == "" && params.Path == "" {
				return
			}
			var obj manifold.Object
			if params.ID != "" {
				obj = s.Root.FindID(params.ID)
			} else {
				obj = s.Root.FindChild(params.Path)
			}
			if obj == nil {
				return
			}
			s.view.ActiveID = obj.ID()
			ref := &file.Reference{}
			rv := reflect.ValueOf(ref)
			obj.ValueTo(rv)
			if ref.Filepath != "" {
				cmds.Execute("editor.open", map[string]interface{}{
					"Filename": ref.Filepath,
				})
			}
		},
	})

	cmds.Register(cmd.Definition{
		ID:       "object.new.empty",
		Label:    "New Empty",
		Category: "Objects",
		Desc:     "Creates an empty object",
		Run:      s.objectNewEmpty,
	})

	cmds.Register(cmd.Definition{
		ID:       "object.new.from",
		Label:    "New Prefab",
		Category: "Objects",
		Desc:     "Creates an object from a prefab",
		Run:      s.objectNewFrom,
	})

	cmds.Register(cmd.Definition{
		ID:       "object.expand",
		Label:    "Expand",
		Category: "Objects",
		Desc:     "Expands an object loading lazy children",
		Run:      s.objectExpand,
	})

	cmds.Register(cmd.Definition{
		ID:       "object.refresh",
		Label:    "Refresh",
		Category: "Objects",
		Desc:     "Refreshes an object",
		Run:      s.objectRefresh,
	})

	cmds.Register(cmd.Definition{
		ID:       "object.move",
		Label:    "Move",
		Category: "Objects",
		Desc:     "Moves an object",
		Run:      s.objectMove,
	})

	cmds.Register(cmd.Definition{
		ID:       "object.duplicate",
		Label:    "Duplicate",
		Category: "Objects",
		Desc:     "Duplicates an object",
		Run:      s.objectDuplicate,
	})

	cmds.Register(cmd.Definition{
		ID:       "object.delete",
		Label:    "Delete",
		Category: "Objects",
		Desc:     "Deletes an object",
		Run:      s.objectDelete,
	})

	cmds.Register(cmd.Definition{
		ID:       "object.rename",
		Label:    "Rename",
		Category: "Objects",
		Desc:     "Renames an object",
		Run:      s.objectRename,
	})

	cmds.Register(cmd.Definition{
		ID:       "object.main.edit",
		Label:    "Edit Main",
		Category: "Objects",
		Desc:     "Edits an object main component",
		Run:      s.objectMainEdit,
	})

	cmds.Register(cmd.Definition{
		ID:       "object.component.move",
		Label:    "Move",
		Category: "Components",
		Desc:     "Moves a component",
		Run:      s.objectComponentMove,
	})

	cmds.Register(cmd.Definition{
		ID:       "object.component.enable",
		Label:    "Toggle",
		Category: "Components",
		Desc:     "Toggles a component",
		Run:      s.objectComponentEnable,
	})

	cmds.Register(cmd.Definition{
		ID:       "object.component.reload",
		Label:    "Reload",
		Category: "Components",
		Desc:     "Reloads a component",
		Run:      s.objectComponentReload,
	})

	cmds.Register(cmd.Definition{
		ID:       "object.component.add",
		Label:    "Add Component",
		Category: "Objects",
		Desc:     "Adds a component to an object",
		Run:      s.objectComponentAdd,
	})

	cmds.Register(cmd.Definition{
		ID:       "object.component.remove",
		Label:    "Remove Component",
		Category: "Objects",
		Desc:     "Removes a component from an object",
		Run:      s.objectComponentRemove,
	})

	cmds.Register(cmd.Definition{
		ID:       "object.method.call",
		Label:    "Call Method",
		Category: "Objects",
		Desc:     "Calls a method on an object",
		Run:      s.objectMethodCall,
	})

	// TODO: refactor this, too much going on here
	cmds.Register(cmd.Definition{
		ID:       "object.value.set",
		Label:    "Set Value",
		Category: "Objects",
		Desc:     "Sets a component value for an object",
		Run:      s.objectValueSet,
	})

	// TODO: refactor this, too much going on here
	cmds.Register(cmd.Definition{
		ID:       "object.value.remove",
		Label:    "Remove Value",
		Category: "Objects",
		Desc:     "Removes a component value on an object",
		Run:      s.objectValueRemove,
	})

	// TODO: refactor this, too much going on here
	cmds.Register(cmd.Definition{
		ID:       "object.value.add",
		Label:    "Add Value",
		Category: "Objects",
		Desc:     "Adds a component value on an object",
		Run:      s.objectValueAdd,
	})

	cmds.Register(cmd.Definition{
		ID:       "object.value.delete",
		Label:    "Delete Map Value",
		Category: "Objects",
		Desc:     "Deletes a component value key on an object",
		Run:      s.objectValueDelete,
	})

}

func (s *Service) objectNewEmpty(params struct {
	ID   string
	Name string
}) {
	if params.Name == "" {
		params.Name = "EmptyObject"
	}
	p := s.Root.FindID(params.ID)
	if p == nil {
		p = s.Root
	}
	n := object.New(params.Name)
	p.AppendChild(n)
}

func (s *Service) objectNewFrom(params struct {
	ID     string
	Prefab string
}) {
	if params.Prefab == "" {
		return
	}
	p := s.Root.FindID(params.ID)
	if p == nil {
		p = s.Root
	}
	switch params.Prefab {
	case "userpkg":
		// TODO: move package creation into component on object
		pkgName := "mypackage"
		if err := s.Image.CreateUserPackage(pkgName); err != nil {
			panic(err)
		}
		pkgPath := s.Image.UserPackagePath(pkgName)

		n := object.New(pkgName)
		c1 := library.Lookup("file.Path").New()
		c1.SetEnabled(true)
		if err := c1.SetField("Filepath", pkgPath); err != nil {
			panic(err)
		}
		if err := c1.SetField("SyncName", true); err != nil {
			panic(err)
		}
		n.AppendComponent(c1)
		if err := s.MountedComponent(c1, n); err != nil {
			panic(err)
		}
		c2 := library.Lookup("file.Explorer").New()
		c2.SetEnabled(true)
		n.AppendComponent(c2)
		if err := s.MountedComponent(c2, n); err != nil {
			panic(err)
		}
		p.AppendChild(n)
	case "worksitesrc":
		wd, err := os.Getwd()
		if err != nil {
			panic(err)
		}
		n := object.New("Source")
		c1 := library.Lookup("file.Path").New()
		c1.SetEnabled(true)
		if err := c1.SetField("Filepath", wd); err != nil {
			panic(err)
		}
		n.AppendComponent(c1)
		if err := s.MountedComponent(c1, n); err != nil {
			panic(err)
		}
		c2 := library.Lookup("file.Explorer").New()
		c2.SetEnabled(true)
		n.AppendComponent(c2)
		if err := s.MountedComponent(c2, n); err != nil {
			panic(err)
		}
		p.AppendChild(n)
	case "tractorsrc":
		src := os.Getenv("TRACTOR_SRC")
		if src == "" {
			log.Println("unable to find tractor source")
			return
		}
		n := object.New("Tractor")
		// TODO: eventually look at source bundled in binary
		c1 := library.Lookup("file.Path").New()
		c1.SetEnabled(true)
		if err := c1.SetField("Filepath", os.Getenv("TRACTOR_SRC")); err != nil {
			panic(err)
		}
		n.AppendComponent(c1)
		if err := s.MountedComponent(c1, n); err != nil {
			panic(err)
		}
		c2 := library.Lookup("file.Explorer").New()
		c2.SetEnabled(true)
		n.AppendComponent(c2)
		if err := s.MountedComponent(c2, n); err != nil {
			panic(err)
		}
		p.AppendChild(n)
	default:
	}
}

func (s *Service) objectExpand(params struct {
	ID string
}) {
	if params.ID == "" {
		return
	}
	obj := s.Root.FindID(params.ID)
	if obj == nil {
		return
	}

	for _, com := range obj.Components() {
		lc, ok := com.Pointer().(ui.LazyChilder)
		if ok && lc.LazyChildren() && len(obj.Children()) == 0 {
			if cp, ok := com.Pointer().(library.ChildProvider); ok {
				for _, o := range cp.ChildNodes() {
					obj.AppendChild(o)
				}
			}
		}
	}
}

func (s *Service) objectMove(params struct {
	ID     string
	Index  int
	Parent string
}) {
	if params.ID == "" {
		return
	}
	obj := s.Root.FindID(params.ID)
	if obj == nil {
		return
	}
	if params.Parent != "" {
		p := s.Root.FindID(params.Parent)
		if p != nil {
			obj.Parent().RemoveChild(obj)
			p.AppendChild(obj)
		}
	}
	if params.Index > -1 {
		if err := obj.SetSiblingIndex(params.Index); err != nil {
			log.Println(err)
		}
	}
}

func (s *Service) objectMethodCall(params struct {
	Path string
}) (interface{}, error) {
	if params.Path == "" {
		return nil, nil // TODO: error
	}
	obj := s.Root.FindChild(params.Path)
	localPath := params.Path[len(obj.Path())+1:]
	// TODO: support args
	var ret interface{}
	err := obj.CallMethod(localPath, nil, &ret)
	return ret, err
}

func (s *Service) objectValueAdd(params struct {
	Path        string
	Type        string
	Value       interface{}
	IntValue    *int
	RefValue    *string
	IntSlice    *[]int
	StringSlice *[]string
	Key         string // maps only
}) {
	obj := s.Root.FindChild(params.Path)
	if obj == nil {
		return
	}
	localPath := params.Path[len(obj.Path())+1:]
	f, t, err := obj.GetField(localPath)
	if err != nil {
		log.Println(fmt.Errorf("unable to get field: %s", localPath))
		return
	}

	switch t.Kind() {
	case reflect.Map:
		// assumes map[string]string
		rm := reflect.ValueOf(f)
		rk := reflect.ValueOf(params.Key)
		rv := reflect.ValueOf(params.Value)
		rm.SetMapIndex(rk, rv)
		obj.SetField(localPath, rm.Interface())
	case reflect.Slice, reflect.Array:
		rv := reflect.ValueOf(f)
		idx := rv.Len()
		nv := reflect.New(t.Elem())
		rv = reflect.Append(rv, reflect.Indirect(nv))
		obj.SetField(localPath, rv.Interface())
		// NOTE: this must also done in setValue
		v := params.Value
		switch params.Type {
		case "time":
			v, err = time.Parse("15:04", v.(string))
			if err != nil {
				log.Println(err)
				return
			}
		case "date":
			v, err = time.Parse("2006-01-02", v.(string))
			if err != nil {
				log.Println(err)
				return
			}
		}
		obj.SetField(filepath.Join(localPath, strconv.Itoa(idx)), v)
	}

}

func (s *Service) objectValueRemove(params struct {
	Path        string
	Type        string
	Value       interface{}
	IntValue    *int
	RefValue    *string
	IntSlice    *[]int
	StringSlice *[]string
}) {
	obj := s.Root.FindChild(params.Path)
	if obj == nil {
		return
	}
	localPath := params.Path[len(obj.Path())+1:]
	f, _, err := obj.GetField(localPath)
	if err != nil {
		log.Println(fmt.Errorf("unable to get field: %s", localPath))
		return
	}
	rv := reflect.ValueOf(f)
	idx := *params.IntValue
	a := rv.Slice(0, idx)
	b := rv.Slice(idx+1, rv.Len())
	rv = reflect.AppendSlice(a, b)
	obj.SetField(localPath, rv.Interface())
}

func (s *Service) objectValueSet(params struct {
	Path        string
	Type        string
	Value       interface{}
	IntValue    *int
	RefValue    *string
	RefType     *string
	IntSlice    *[]int
	StringSlice *[]string
}) {
	obj := s.Root.FindChild(params.Path)
	if obj == nil {
		log.Printf("obj not found: %+v", params)
		return
	}
	localPath := params.Path[len(obj.Path())+1:]
	switch {
	case params.IntSlice != nil:
		obj.SetField(localPath, *params.IntSlice)
	case params.StringSlice != nil:
		obj.SetField(localPath, *params.StringSlice)
	case params.IntValue != nil:
		obj.SetField(localPath, *params.IntValue)
	case params.RefValue != nil:
		refPath := filepath.Dir(*params.RefValue) // TODO: support subfields
		refNode := s.Root.FindChild(refPath)
		parts := strings.SplitN(localPath, "/", 2)
		refType := obj.Component(parts[0]).FieldType(parts[1])
		if refNode != nil {
			typeSelector := (*params.RefValue)[len(refNode.Path())+1:]
			c := refNode.Component(typeSelector)
			if c != nil {
				if err := obj.SetField(localPath, c.Pointer()); err != nil {
					log.Println(err)
				}
			} else {
				// interface reference
				ptr := reflect.New(refType)
				refNode.ValueTo(ptr)
				if ptr.IsValid() {
					if err := obj.SetField(localPath, reflect.Indirect(ptr).Interface()); err != nil {
						log.Println(err)
					}
				}
			}

		}
	default:
		// NOTE: this must also done in addValue
		v := params.Value
		var err error
		switch params.Type {
		case "time":
			v, err = time.Parse("15:04", v.(string))
			if err != nil {
				log.Println(err)
				return
			}
		case "date":
			v, err = time.Parse("2006-01-02", v.(string))
			if err != nil {
				log.Println(err)
				return
			}
		}
		obj.SetField(localPath, v)
	}
}

func (s *Service) objectValueDelete(params struct {
	Path string
}) {
	obj := s.Root.FindChild(params.Path)
	if obj == nil {
		log.Printf("obj not found: %+v", params)
		return
	}
	localPath := params.Path[len(obj.Path())+1:]
	key := path.Base(localPath)
	mapPath := path.Dir(localPath)
	m, t, err := obj.GetField(mapPath)
	if err != nil {
		log.Println(fmt.Errorf("unable to get map field: %s", mapPath))
		return
	}
	if t.Kind() != reflect.Map {
		log.Println(fmt.Errorf("not a map field: %s", mapPath))
		return
	}

	// TODO: uhhhh
	mm, ok := m.(map[string]string)
	if !ok {
		log.Println(fmt.Errorf("not a map[string]string field: %s", mapPath))
		return
	}

	delete(mm, key)
	obj.SetField(mapPath, mm)
}

func (s *Service) objectComponentRemove(params struct {
	ID        string
	Component string
}) {
	if params.ID == "" || params.Component == "" {
		return
	}
	obj := s.Root.FindID(params.ID)
	if obj == nil {
		return
	}
	com := obj.Component(params.Component)
	obj.RemoveComponent(com)
	if com.ID() == obj.ID() {
		if err := s.Image.DestroyObjectPackage(obj); err != nil {
			L.Error(s.Log, err)
		}
	}
}

func (s *Service) objectComponentAdd(params struct {
	ID   string
	Name string
}) {
	if params.ID == "" || params.Name == "" {
		return
	}
	obj := s.Root.FindID(params.ID)
	if obj == nil {
		return
	}
	v := library.Lookup(params.Name).New()
	obj.AppendComponent(v)
	if err := s.MountedComponent(v, obj); err != nil {
		log.Println(err)
	}
}

func (s *Service) objectComponentEnable(params struct {
	ID        string
	Component string
	Enable    bool
}) {
	if params.ID == "" {
		return
	}
	obj := s.Root.FindID(params.ID)
	if obj == nil {
		return
	}
	com := obj.Component(params.Component)
	if com != nil {
		com.SetEnabled(params.Enable)
	}
}

func (s *Service) objectComponentMove(params struct {
	ID   string
	From int
	To   int
}) {
	if params.ID == "" {
		return
	}
	obj := s.Root.FindID(params.ID)
	if obj == nil {
		return
	}
	obj.MoveComponentAt(params.From, params.To)
}

func (s *Service) objectComponentReload(params struct {
	ID        string
	Component string
}) {
	if params.ID == "" {
		return
	}
	obj := s.Root.FindID(params.ID)
	if obj == nil {
		return
	}
	com := obj.Component(params.Component)
	if com != nil {
		if err := com.Reload(); err != nil {
			// TODO: return error
			return
		}
	}
}

func (s *Service) objectMainEdit(params struct {
	ID string
}) {
	if params.ID == "" {
		return
	}
	obj := s.Root.FindID(params.ID)
	if obj == nil {
		return
	}
	s.Image.CreateObjectPackage(obj)
	// TODO: actually open it. prob need URIs since images might not be on filesystem.
}

func (s *Service) objectRename(params struct {
	ID   string
	Name string
}) {
	if params.ID == "" {
		return
	}
	obj := s.Root.FindID(params.ID)
	if obj == nil {
		return
	}
	if params.Name != "" {
		obj.SetName(params.Name)
	}
}

func (s *Service) objectDuplicate(params struct {
	ID string
}) {
	if params.ID == "" {
		return
	}
	obj := s.Root.FindID(params.ID)
	if obj == nil {
		return
	}
	parent := obj.Parent()
	dup := object.Duplicate(obj)
	parent.AppendChild(dup)
}

func (s *Service) objectDelete(params struct {
	ID string
}) {
	if params.ID == "" {
		return
	}
	s.Root.RemoveID(params.ID)
}

func (s *Service) objectRefresh(params struct {
	ID string
}) {
	if params.ID == "" {
		return
	}
	obj := s.Root.FindID(params.ID)
	if obj == nil {
		return
	}
	if err := obj.Refresh(); err != nil {
		// TODO: retrun error
		return
	}
}
