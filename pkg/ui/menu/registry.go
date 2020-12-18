package menu

import (
	"context"
	"sync"
)

type Item struct {
	Label     string
	Icon      string `json:",omitempty"`
	Cmd       string `json:",omitempty"`
	Params    Params `json:",omitempty"`
	Submenu   []Item `json:",omitempty"`
	Separator bool   `json:",omitempty"`
}

type Params map[string]string

type Registry struct {
	menus map[string]func(context.Context) []Item
	sync.Mutex
}

func (r *Registry) Get(menuID string, ctx context.Context) []Item {
	r.Lock()
	defer r.Unlock()
	return r.menus[menuID](ctx)
}

func (r *Registry) Register(menuID string, fn func(context.Context) []Item) {
	r.Lock()
	defer r.Unlock()
	r.menus[menuID] = fn
}
