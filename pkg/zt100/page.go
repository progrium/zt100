package zt100

import (
	"encoding/json"
	"fmt"
	"image/color"
	"net/http"
	"strings"

	"github.com/manifold/tractor/pkg/manifold"
	"github.com/manifold/tractor/pkg/manifold/comutil"
	"github.com/manifold/tractor/pkg/misc/notify"
	httplib "github.com/manifold/tractor/pkg/stdlib/http"
)

type Page struct {
	Title      string
	HideInMenu bool

	OID        string `tractor:"hidden"`
	ObjectName string `json:"Name"`

	events *httplib.EventSource
	server *Server
	object manifold.Object
}

func (p *Page) Initialize() {
	p.events = &httplib.EventSource{}
}

func (t *Page) Name() string {
	return t.object.Name()
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
	p.OID = obj.ID()
	p.ObjectName = obj.Name()

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

func (s *Page) Block(id string) *Block {
	for _, b := range s.Blocks() {
		if b.ID == id {
			return b
		}
	}
	return nil
}

func (p *Page) Blocks() (blocks []*Block) {
	for _, com := range p.object.Components() {
		s, ok := com.Pointer().(*Block)
		if !ok {
			continue
		}
		blocks = append(blocks, s)
	}
	return blocks
}

func (p *Page) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if p.events.IsEventStream(r) {
		p.events.ServeHTTP(w, r)
		return
	}

	ctx := FromContext(r.Context())

	if ctx.Block != nil {
		ctx.Blocks = []*Block{ctx.Block}
	}

	data, err := json.Marshal(ctx)
	if err != nil {
		panic(err)
	}

	var rgb color.RGBA
	rgb, err = ParseHexColor(ctx.Demo.Color)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	p.server.Template.ExecuteTemplate(w, "page.html", map[string]interface{}{
		"Data":  string(data),
		"Color": rgb,
		"Live":  r.URL.Query().Get("live") != "0",
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
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
