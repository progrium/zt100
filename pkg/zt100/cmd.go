package zt100

import (
	"fmt"
	"log"

	"github.com/manifold/tractor/pkg/core/cmd"
	"github.com/manifold/tractor/pkg/manifold/library"
	"github.com/manifold/tractor/pkg/manifold/object"
)

func (c *Server) ContributeCommands(cmds *cmd.Registry) {
	cmds.Register(cmd.Definition{
		ID:       "zt100.new-tenant",
		Label:    "New Tenant",
		Category: "zt100",
		Desc:     "",
		Run: func(params struct {
			ID   string
			Name string
		}) {
			if params.Name == "" {
				params.Name = "newtenant"
			}
			n := object.New(params.Name)

			tc := library.Lookup("zt100.Tenant").New()
			tc.SetEnabled(true)
			n.AppendComponent(tc)

			// thc := library.Lookup("zt100.Theme").New()
			// thc.SetEnabled(true)
			// n.AppendComponent(thc)

			p := c.object.Root().FindID(params.ID)
			if p == nil {
				p = c.object.Root()
			}
			p.AppendChild(n)
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
		},
	})

	cmds.Register(cmd.Definition{
		ID:       "zt100.new-section",
		Label:    "New Page",
		Category: "zt100",
		Desc:     "",
		Run: func(params struct {
			PageID  string
			BlockID string
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

		},
	})

	cmds.Register(cmd.Definition{
		ID:       "zt100.override-text",
		Label:    "Override Block Text",
		Category: "zt100",
		Desc:     "",
		Run: func(params struct {
			Text    string
			Tenant  string
			App     string
			Page    string
			Section string
			Key     string
		}) error {
			_, _, p, _ := c.Lookup(params.Tenant, params.App, params.Page, "")
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
