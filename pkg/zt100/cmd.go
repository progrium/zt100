package zt100

import (
	"fmt"
	"io"
	"log"

	"github.com/progrium/zt100/pkg/core/cmd"
	"github.com/progrium/zt100/pkg/manifold/library"
	"github.com/progrium/zt100/pkg/manifold/object"
)

func (c *Server) ContributeCommands(cmds *cmd.Registry) {
	cmds.Register(cmd.Definition{
		ID:       "zt100.new-demo",
		Label:    "New Demo",
		Category: "zt100",
		Desc:     "",
		Run: func(params struct {
			ID       string
			Name     string
			Domain   string
			Color    string
			Vertical string
			Features []interface{}
		}) {
			if params.Name == "" {
				params.Name = "newdemo"
			}
			n := object.New(params.Name)

			var flags []string
			for _, feat := range params.Features {
				flags = append(flags, feat.(string))
			}

			tc := library.Lookup("zt100.Demo").New()
			tc.SetEnabled(true)
			tc.SetField("Domain", params.Domain)
			tc.SetField("Vertical", params.Vertical)
			tc.SetField("Features", flags)
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
		ID:       "zt100.new-block",
		Label:    "New Block",
		Category: "zt100",
		Desc:     "",
		Run: func(params struct {
			PageID   string
			BaseName string
		}) {
			p := c.object.Root().FindID(params.PageID)
			if p == nil {
				return
			}

			pc := library.Lookup("zt100.Block").New()
			pc.SetEnabled(true)
			pc.SetField("BaseName", params.BaseName)
			pc.SetField("Name", params.BaseName)
			p.AppendComponent(pc)

			if err := c.Objects.MountedComponent(pc, p); err != nil {
				log.Println(err)
			}
		},
	})

	cmds.Register(cmd.Definition{
		ID:       "zt100.block.attach-image",
		Label:    "Attach Block Image",
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

			com := p.Component(params.BlockID)
			if com != nil && params.ImageSize > 0 {
				d := make([]byte, params.ImageSize)
				defer params.Image.Close()
				if _, err := params.Image.Read(d); err == nil {
					if err := com.SetField("Image", d); err != nil {
						log.Println(err)
					}
					// os.Mkdir("local/uploads", 0755)
					// if err := ioutil.WriteFile(fmt.Sprintf("local/uploads/%s.png", pc.ID()), d, 0644); err != nil {
					// 	log.Println(err)
					// }
				}
			}
		},
	})

	cmds.Register(cmd.Definition{
		ID:       "zt100.block.set-text",
		Label:    "Set Block Text",
		Category: "zt100",
		Desc:     "",
		Run: func(params struct {
			PageID  string
			BlockID string
			Key     string
			Text    string
		}) error {
			p := c.object.Root().FindID(params.PageID)
			if p == nil {
				return nil
			}

			com := p.Component(params.BlockID)
			if com != nil {
				s := com.Pointer().(*Block)
				s.Text[params.Key] = params.Text
				return com.SetField(fmt.Sprintf("Text/%s", params.Key), params.Text)
			}

			return nil
		},
	})
}
