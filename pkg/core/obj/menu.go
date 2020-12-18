package obj

import (
	"context"
	"strings"

	"github.com/progrium/zt100/pkg/manifold"
	"github.com/progrium/zt100/pkg/manifold/library"
	"github.com/progrium/zt100/pkg/ui/menu"
)

type MenuProvider interface {
	ObjectMenu(menuID string) []menu.Item
}

func (s *Service) ContextObject(ctx context.Context) (manifold.Object, bool) {
	v := ctx.Value("obj")
	if v == nil {
		return nil, false
	}
	id, ok := v.(string)
	if !ok {
		return nil, false
	}
	obj := s.Root.FindID(id)
	if obj == nil {
		return nil, false
	}
	return obj, true
}

func ObjectMenu(obj manifold.Object, menuID string) ([]menu.Item, bool) {
	var items []menu.Item
	for _, com := range obj.Components() {
		if mp, ok := com.Pointer().(MenuProvider); ok {
			objItems := mp.ObjectMenu(menuID)
			if len(objItems) > 0 {
				items = append(items, menu.Item{Separator: true})
				items = append(items, objItems...)
			}
		}
	}
	return items, true
}

func (s *Service) ContributeMenus(menus *menu.Registry) {
	components := make(map[string][]menu.Item)
	for _, com := range library.Registered() {
		if strings.HasSuffix(com.Name, ".Main") {
			continue
		}
		parts := strings.Split(com.Name, ".")
		if _, ok := components[parts[0]]; !ok {
			components[parts[0]] = []menu.Item{}
		}
		components[parts[0]] = append(components[parts[0]], menu.Item{
			Label: com.Name,
			Cmd:   "object.component.add",
			Params: menu.Params{
				"Name": com.Name,
			},
		})
	}
	var componentMenu []menu.Item
	for k, v := range components {
		componentMenu = append(componentMenu, menu.Item{
			Label:   k,
			Submenu: v,
		})
	}

	menus.Register("inspector/add", func(ctx context.Context) []menu.Item {
		return componentMenu
	})

	menus.Register("explorer/context", func(ctx context.Context) []menu.Item {
		items := []menu.Item{
			{Label: "New", Submenu: []menu.Item{
				{Label: "Empty Object", Cmd: "object.new.empty"},
				{Label: "Dev Prefabs", Submenu: []menu.Item{
					{Label: "User Package", Cmd: "object.new.from", Params: menu.Params{"Prefab": "userpkg"}},
					{Label: "Worksite Source", Cmd: "object.new.from", Params: menu.Params{"Prefab": "worksitesrc"}},
					{Label: "Tractor Source", Cmd: "object.new.from", Params: menu.Params{"Prefab": "tractorsrc"}},
				}},
			}},
		}
		obj, ok := s.ContextObject(ctx)
		if !ok {
			return items
		}
		if i, ok := ObjectMenu(obj, "explorer/context"); ok {
			items = append(items, i...)
		}
		items = append(items, []menu.Item{
			{Separator: true},
			{Cmd: "copy-path", Label: "Copy Path", Params: menu.Params{
				"Path": obj.Path(),
			}},
			{Separator: true},
			{Label: "Add Component", Submenu: componentMenu},
			{Separator: true},
			{Cmd: "object.duplicate", Label: "Duplicate"},
			{Cmd: "object.rename", Label: "Rename"},
			{Cmd: "object.delete", Label: "Delete"},
		}...)
		return items
	})

	menus.Register("object/manage", func(ctx context.Context) []menu.Item {
		return []menu.Item{
			{Cmd: "object.refresh", Label: "Refresh", Icon: "sync"},
			{Cmd: "object.main.edit", Label: "Edit Main", Icon: "edit"},
			{Cmd: "editor.open", Label: "View Data", Icon: "file-alt"},
		}
	})

	menus.Register("component/manage", func(ctx context.Context) []menu.Item {
		obj, ok := s.ContextObject(ctx)
		if !ok {
			return []menu.Item{}
		}
		return []menu.Item{
			{Cmd: "object.component.reload", Label: "Reload", Icon: "sync"},
			{Cmd: "copy-path", Label: "Copy Path", Icon: "copy", Params: menu.Params{
				"Path": obj.Path(),
			}},
			{Cmd: "editor.open", Label: "View Source", Icon: "edit"},
			{Cmd: "object.component.remove", Label: "Remove", Icon: "trash"},
		}
	})

	menus.Register("component/reference", func(ctx context.Context) []menu.Item {
		return []menu.Item{
			{Cmd: "object.value.clipboard-path", Label: "Paste Path"},
		}
	})
}
