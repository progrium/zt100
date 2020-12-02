package zt100

import (
	"github.com/manifold/tractor/pkg/manifold"
	"github.com/manifold/tractor/pkg/manifold/comutil"
	httplib "github.com/manifold/tractor/pkg/stdlib/http"
	"github.com/manifold/tractor/pkg/ui"
	"github.com/manifold/tractor/pkg/ui/menu"
)

type App struct {
	object manifold.Object
}

func (a *App) Mounted(obj manifold.Object) error {
	a.object = obj
	return nil
}

func (a *App) PageMenu() (items []MenuItem) {
	for _, obj := range a.object.Children() {
		pagecom := obj.Component("zt100.Page")
		if pagecom == nil {
			continue
		}
		page := pagecom.Pointer().(*Page)
		items = append(items, MenuItem{
			Title: page.Title,
			Page:  obj.Name(),
		})
	}
	return items
}

func (a *App) OpenInBrowser(path ...string) ui.Script {

	var tnt Tenant
	tobj := comutil.AncestorValue(a.object, &tnt)

	var srv httplib.Server
	comutil.AncestorValue(a.object, &srv)
	return srv.OpenInBrowser("t", tobj.Name(), a.object.Name())
}

func (s *App) ObjectMenu(menuID string) []menu.Item {
	switch menuID {
	case "explorer/context":
		return []menu.Item{
			{Cmd: "zt100.new-page", Label: "New Page", Icon: "plus"},
			{Cmd: "file.open", Label: "Open in Browser", Icon: "browser"},
		}
	default:
		return []menu.Item{}
	}
}