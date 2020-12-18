package obj

import (
	"fmt"
	"path"
	"strings"
	"sync"

	"github.com/progrium/zt100/pkg/manifold"
	"github.com/progrium/zt100/pkg/manifold/library"
	"github.com/progrium/zt100/pkg/manifold/prefab"
	"github.com/progrium/zt100/pkg/ui"
	"github.com/progrium/zt100/pkg/ui/action"
	"github.com/progrium/zt100/pkg/ui/field"
)

type UIProvider interface {
	InspectorUI() []ui.Element
}

type ActionProvider interface {
	InspectorActions() []string
}

type ViewState struct {
	Components []ui.ComponentType
	Prefabs    []ui.Prefab

	Objects map[string]ui.Object

	RootID   string
	ActiveID string

	mu sync.Mutex
}

func (s *Service) InitializeState() (name string, v interface{}) {
	s.view = &ViewState{
		Objects:  make(map[string]ui.Object),
		ActiveID: "",
	}

	name = "obj"
	v = s.view

	for _, com := range library.Registered() {
		if strings.HasSuffix(com.Name, ".Main") {
			continue
		}
		s.view.Components = append(s.view.Components, ui.ComponentType{
			Name:     com.Name,
			Filepath: com.Filepath,
		})
	}

	for _, pf := range prefab.Registered() {
		s.view.Prefabs = append(s.view.Prefabs, ui.Prefab{
			Name: pf.Name,
			ID:   pf.ID,
		})
	}

	return
}

func (s *Service) UpdateState() error {
	// reset/clear nodes
	s.view.Objects = make(map[string]ui.Object)

	// walk every object in the tree
	return manifold.Walk(s.Root, func(o manifold.Object) error {
		// start a node struct based on node passed in
		obj := ui.Object{
			Name:       o.Name(),
			Active:     true,
			Attrs:      o.Attrs().Snapshot(),
			Icon:       o.Icon(),
			Path:       o.Path(),
			Index:      o.SiblingIndex(),
			ID:         o.ID(),
			Components: []ui.Component{},
		}

		if o.Parent() != nil {
			if o.Parent().ID() != o.Root().ID() {
				obj.ParentID = o.Parent().ID()
			} else {
				s.view.RootID = o.ID()
			}
		}

		lazyChildren := false
		for _, com := range o.Components() {
			// get all the fields of the component
			fields := field.FromComponent(com)

			// see if component provides custom ui
			var customUI []ui.Element
			uip, ok := com.Pointer().(UIProvider)
			if ok {
				customUI = uip.InspectorUI()
			}

			lc, ok := com.Pointer().(ui.LazyChilder)
			if ok {
				lazyChildren = lazyChildren || lc.LazyChildren()
			}

			var actions []ui.Action
			ap, ok := com.Pointer().(ActionProvider)
			if ok {
				names := ap.InspectorActions()
				for _, a := range action.FromMethods(com) {
					if strInSlice(names, a.Name) {
						actions = append(actions, a)
					}
				}
			} else {
				for _, a := range action.FromMethods(com) {
					if a.Out != nil && fmt.Sprintf("%s.%s", path.Base(a.Out.PkgPath()), a.Out.Name()) == "ui.Script" {
						actions = append(actions, a)
					}
				}
			}

			// look up the filepath for this component
			var filepath string
			// if com.ID() != "" {
			// 	rc := library.LookupID(com.ID())
			// 	if rc == nil {
			// 		return fmt.Errorf("component ID not in library: %q", com.ID())
			// 	}
			// 	filepath = rc.Filepath
			// } else {
			rc := library.Lookup(com.Name())
			if rc == nil {
				return fmt.Errorf("component name not in library: %q", com.Name())
			}
			filepath = rc.Filepath
			// }

			// look up related components to this component
			var related []string
			for _, rc := range library.Related(library.Lookup(com.Name())) {
				related = append(related, rc.Name)
			}

			// add component to frontend node's components
			obj.Components = append(obj.Components, ui.Component{
				Name:     com.Name(),
				Icon:     com.Icon(),
				Index:    com.Index(),
				Key:      com.Key(),
				Enabled:  com.Enabled(),
				Filepath: filepath,
				Fields:   fields,
				Actions:  actions,
				Related:  related,
				CustomUI: customUI,
			})
		}

		children := o.Children()
		if len(children) > 0 || lazyChildren {
			obj.HasChildren = true
		}
		for _, child := range children {
			obj.Children = append(obj.Children, child.ID())
		}

		// add the node to state
		s.view.mu.Lock()
		s.view.Objects[o.ID()] = obj
		s.view.mu.Unlock()

		return nil
	})
}

func strInSlice(strs []string, str string) bool {
	for _, s := range strs {
		if s == str {
			return true
		}
	}
	return false
}
