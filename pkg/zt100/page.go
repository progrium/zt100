package zt100

import (
	"encoding/json"
	"fmt"
	"image/color"
	"net/http"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/gorilla/mux"
	"github.com/manifold/tractor/pkg/manifold"
	"github.com/manifold/tractor/pkg/manifold/comutil"
	"github.com/manifold/tractor/pkg/misc/notify"
	httplib "github.com/manifold/tractor/pkg/stdlib/http"
)

type Page struct {
	Title string
	Name  string `tractor:"hidden"`
	OID   string `tractor:"hidden"`

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
	p.Name = obj.Name()
	p.OID = obj.ID()

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

func (p *Page) Sections() (sections []*Section) {
	for _, com := range p.object.Components() {
		s, ok := com.Pointer().(*Section)
		if !ok {
			continue
		}
		sections = append(sections, s)
	}
	return sections
}

func (p *Page) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if p.events.IsEventStream(r) {
		p.events.ServeHTTP(w, r)
		return
	}
	vars := mux.Vars(r)
	prospect, app, page, sections := p.server.Lookup(vars["prospect"], vars["app"], vars["page"], r.URL.Query().Get("section"))
	if prospect == nil || app == nil || page == nil {
		http.Error(w, "page not found", http.StatusNotFound)
		return
	}
	var sectionData []interface{}
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
		//el = append(el, fmt.Sprintf(`h((await import("/blocks/%s.js")).default, {overrides: config.Overrides[%d], section: "%s"}, "Hello")`, name, idx, com.Key()))
		sectionData = append(sectionData, map[string]interface{}{
			"Block": name,
			"Idx":   idx,
			"Key":   com.ID(),
		})
		idx++
	}
	t := prospect.Component("zt100.Prospect").Pointer().(*Prospect)
	a := app.Component("zt100.App").Pointer().(*App)
	config, err := json.Marshal(map[string]interface{}{
		"Prospect":       prospect.Name(),
		"ProspectColor":  t.Color,
		"ProspectDomain": t.Domain,
		"PageMenu":       a.PageMenu(),
		"Page":           page.Name(),
		"PageOID":        page.ID(),
		"App":            app.Name(),
		"Overrides":      overrides,
	})
	if err != nil {
		panic(err)
	}

	ts, err := template.ParseFiles(
		"./tmpl/app.page.html",
	)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var rgb color.RGBA
	rgb, err = ParseHexColor(t.Color)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	err = ts.Execute(w, map[string]interface{}{
		"Config":   string(config),
		"Color":    rgb,
		"Sections": sectionData,
		"Live":     r.URL.Query().Get("live") != "0",
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	//io.WriteString(w, fmt.Sprintf(Template, string(config), t.Color, strings.Join(el, ",\n")))
}

func ParseHexColor(s string) (c color.RGBA, err error) {
	c.A = 0xff
	switch len(s) {
	case 7:
		_, err = fmt.Sscanf(s, "#%02x%02x%02x", &c.R, &c.G, &c.B)
	case 4:
		_, err = fmt.Sscanf(s, "#%1x%1x%1x", &c.R, &c.G, &c.B)
		// Double the hex digits:
		c.R *= 17
		c.G *= 17
		c.B *= 17
	default:
		err = fmt.Errorf("invalid length, must be 7 or 4")

	}
	return
}
