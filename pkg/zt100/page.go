package zt100

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"path/filepath"
	"strings"

	"github.com/gorilla/mux"
	"github.com/manifold/tractor/pkg/manifold"
	"github.com/manifold/tractor/pkg/manifold/comutil"
	"github.com/manifold/tractor/pkg/misc/notify"
	httplib "github.com/manifold/tractor/pkg/stdlib/http"
)

type Page struct {
	Title string

	events *httplib.EventSource
	server *Server
	object manifold.Object
}

func (p *Page) Initialize() {
	p.events = &httplib.EventSource{}
}

func (p *Page) Mounted(obj manifold.Object) error {
	p.object = obj
	if p.Title == "" {
		name := obj.Name()
		if name == "index" {
			name = "home"
		}
		p.Title = strings.Title(name)
	}

	var server Server
	so := comutil.AncestorValue(obj, &server)
	p.server = &server
	notify.Observe(so, notify.Func(func(event interface{}) error {
		if p == nil {
			return notify.Stop
		}
		p.events.Broadcast([]byte(fmt.Sprintf("%#v", event)))
		return nil
	}))

	return nil
}

func (p *Page) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if p.events.IsEventStream(r) {
		p.events.ServeHTTP(w, r)
		return
	}
	vars := mux.Vars(r)
	tenant, app, page, sections := p.server.Lookup(vars["tenant"], vars["app"], vars["page"], r.URL.Query().Get("section"))
	if tenant == nil || app == nil || page == nil {
		http.Error(w, "page not found", http.StatusNotFound)
		return
	}
	var el []string
	var overrides []map[string]string
	idx := 0
	for _, s := range sections {
		if s.Block == nil {
			continue
		}
		overrides = append(overrides, s.Overrides)
		bo := comutil.Object(s.Block)
		ext := filepath.Ext(bo.Name())
		name := bo.Name()[:len(bo.Name())-len(ext)]
		_, com := page.FindComponent(s)
		el = append(el, fmt.Sprintf(`h((await import("/blocks/%s.js")).default, {overrides: config.Overrides[%d], section: "%s"}, "Hello")`, name, idx, com.Key()))
		idx++
	}
	t := tenant.Component("zt100.Tenant").Pointer().(*Tenant)
	a := app.Component("zt100.App").Pointer().(*App)
	config, err := json.Marshal(map[string]interface{}{
		"Tenant":       tenant.Name(),
		"TenantColor":  t.Color,
		"TenantDomain": t.Domain,
		"PageMenu":     a.PageMenu(),
		"Page":         page.Name(),
		"App":          app.Name(),
		"Overrides":    overrides,
	})
	if err != nil {
		panic(err)
	}
	io.WriteString(w, fmt.Sprintf(Template, string(config), t.Color, strings.Join(el, ",\n")))
}
