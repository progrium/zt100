package zt100

import (
	"net/http"

	"github.com/progrium/zt100/pkg/manifold"
)

type Block struct {
	Text     map[string]string
	BaseName string
	Source   []byte `tractor:"hidden"`
	Image    []byte `tractor:"hidden"`

	Name string `tractor:"hidden"`
	OID  string `tractor:"hidden"`
	ID   string `tractor:"hidden"`

	object manifold.Object
}

func (b *Block) Initialize() {
	if b.Text == nil {
		b.Text = make(map[string]string)
	}
}

func (c *Block) Mounted(obj manifold.Object) error {
	c.object = obj
	_, com := obj.FindComponent(c)
	c.OID = obj.ID()
	c.ID = com.ID()
	if c.Name == "" {
		c.Name = obj.Name()
	}
	return nil
}

func (b *Block) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := FromContext(r.Context())
	w.Header().Set("content-type", "text/javascript")
	if len(b.Source) > 0 {
		w.Write(b.Source)
		return
	}
	w.Write(ctx.Server.Block(b.BaseName).Source)
}
