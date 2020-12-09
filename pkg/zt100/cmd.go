package zt100

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"

	"github.com/manifold/tractor/pkg/core/cmd"
	"github.com/manifold/tractor/pkg/manifold/library"
	"github.com/manifold/tractor/pkg/manifold/object"
)

func (c *Server) ContributeCommands(cmds *cmd.Registry) {
	cmds.Register(cmd.Definition{
		ID:       "zt100.new-prospect",
		Label:    "New Prospect",
		Category: "zt100",
		Desc:     "",
		Run: func(params struct {
			ID       string
			Name     string
			Domain   string
			Color    string
			Vertical string
		}) {
			if params.Name == "" {
				params.Name = "newprospect"
			}
			n := object.New(params.Name)

			tc := library.Lookup("zt100.Prospect").New()
			tc.SetEnabled(true)
			tc.SetField("Domain", params.Domain)
			tc.SetField("Vertical", params.Vertical)
			if params.Color != "" {
				tc.SetField("Color", params.Color)
			} else {
				tc.CallMethod("GetColor", nil, nil)
			}
			n.AppendComponent(tc)

			p := c.object.Root().FindID(params.ID)
			if p == nil {
				p = c.object.Root()
			}
			p.AppendChild(n)

			cmds.Execute("zt100.new-app", map[string]interface{}{
				"ID":   n.ID(),
				"Name": "main",
			})
		},
	})

	cmds.Register(cmd.Definition{
		ID:       "zt100.new-app",
		Label:    "New App",
		Category: "zt100",
		Desc:     "",
		Run: func(params struct {
			ID   string
			Name string
		}) {
			if params.Name == "" {
				params.Name = "newapp"
			}
			n := object.New(params.Name)

			ac := library.Lookup("zt100.App").New()
			ac.SetEnabled(true)
			n.AppendComponent(ac)

			p := c.object.Root().FindID(params.ID)
			if p == nil {
				p = c.object.Root()
			}
			p.AppendChild(n)

			cmds.Execute("zt100.new-page", map[string]interface{}{
				"ID":   n.ID(),
				"Name": "index",
			})
		},
	})

	cmds.Register(cmd.Definition{
		ID:       "zt100.new-section",
		Label:    "New Page",
		Category: "zt100",
		Desc:     "",
		Run: func(params struct {
			PageID    string
			BlockID   string
			Image     io.ReadCloser
			ImageSize int64
		}) {
			p := c.object.Root().FindID(params.PageID)
			if p == nil {
				return
			}

			bo := c.object.Root().FindID(params.BlockID)
			if bo == nil {
				return
			}
			b := bo.Component("zt100.Block").Pointer().(*Block)

			pc := library.Lookup("zt100.Section").New()
			pc.SetEnabled(true)
			pc.SetField("Block", b)
			p.AppendComponent(pc)

			if params.ImageSize > 0 {
				d := make([]byte, params.ImageSize)
				defer params.Image.Close()
				if _, err := params.Image.Read(d); err == nil {
					os.Mkdir("local/uploads", 0755)
					if err := ioutil.WriteFile(fmt.Sprintf("local/uploads/%s.png", pc.ID()), d, 0644); err != nil {
						log.Println(err)
					}
				}
			}

			if err := c.Objects.MountedComponent(pc, p); err != nil {
				log.Println(err)
			}
		},
	})

	cmds.Register(cmd.Definition{
		ID:       "zt100.new-page",
		Label:    "New Page",
		Category: "zt100",
		Desc:     "",
		Run: func(params struct {
			ID   string
			Name string
		}) {
			if params.Name == "" {
				params.Name = "newpage"
			}
			n := object.New(params.Name)

			pc := library.Lookup("zt100.Page").New()
			pc.SetEnabled(true)
			n.AppendComponent(pc)

			p := c.object.Root().FindID(params.ID)
			if p == nil {
				p = c.object.Root()
			}
			p.AppendChild(n)

			cmds.Execute("zt100.new-section", map[string]interface{}{
				"PageID":  n.ID(),
				"BlockID": "bv3plpmg10l0lrrsjsfg",
			})
			cmds.Execute("zt100.new-section", map[string]interface{}{
				"PageID":  n.ID(),
				"BlockID": "bv3pmq6g10l0lrrsjsgg",
			})
		},
	})

	cmds.Register(cmd.Definition{
		ID:       "zt100.override-text",
		Label:    "Override Block Text",
		Category: "zt100",
		Desc:     "",
		Run: func(params struct {
			Text     string
			Prospect string
			App      string
			Page     string
			Section  string
			Key      string
		}) error {
			_, _, p, _ := c.Lookup(params.Prospect, params.App, params.Page, "")
			section := p.Component(params.Section)
			if section != nil {
				s := section.Pointer().(*Section)
				s.Overrides[params.Key] = params.Text
				return section.SetField(fmt.Sprintf("Overrides/%s", params.Key), params.Text)
			}
			return nil
		},
	})
}
