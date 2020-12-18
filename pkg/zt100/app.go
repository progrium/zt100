package zt100

import (
	"github.com/progrium/zt100/pkg/manifold"
	"github.com/progrium/zt100/pkg/ui/menu"
)

type App struct {
	Name   string `tractor:"hidden"`
	OID    string `tractor:"hidden"`
	object manifold.Object
}

type MenuItem struct {
	Title string
	Page  string
}

func (a *App) Mounted(obj manifold.Object) error {
	a.object = obj
	a.Name = a.object.Name()
	a.OID = a.object.ID()
	return nil
}

func (s *App) Page(name string) *Page {
	for _, b := range s.Pages() {
		if b.Name() == name {
			return b
		}
	}
	return nil
}

func (s *App) Pages() (pages []*Page) {
	for _, obj := range s.object.Children() {
		com := obj.Component("zt100.Page")
		if com == nil {
			continue
		}
		p := com.Pointer().(*Page)
		pages = append(pages, p)
	}
	return pages
}

func (a *App) PageMenu() (items []MenuItem) {
	for _, page := range a.Pages() {
		if page.HideInMenu {
			continue
		}
		items = append(items, MenuItem{
			Title: page.Title,
			Page:  page.Name(),
		})
	}
	return items
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
