package zt100

import (
	"fmt"
	"image"
	"log"
	"net/http"

	"github.com/EdlinOrg/prominentcolor"
	"github.com/manifold/tractor/pkg/manifold"
	"github.com/manifold/tractor/pkg/ui"
	"github.com/manifold/tractor/pkg/ui/menu"
)

type Tenant struct {
	Domain string
	Color  string `field:"color"`

	object manifold.Object
}

func (s *Tenant) Mounted(obj manifold.Object) error {
	s.object = obj
	return nil
}

func (t *Tenant) Name() string {
	return t.object.Name()
}

func (t *Tenant) GetColor() ui.Script {
	resp, err := http.Get(fmt.Sprintf("https://logo.clearbit.com/%s", t.Domain))
	if err != nil {
		log.Println("GET ERROR", err)
		return ui.Script{}
	}
	defer resp.Body.Close()
	img, _, derr := image.Decode(resp.Body)
	if derr != nil {
		log.Println("DECODE ERROR", derr)
		return ui.Script{}
	}
	cols, cerr := prominentcolor.KmeansWithArgs(prominentcolor.ArgumentNoCropping, img)
	if cerr != nil {
		log.Println("COLOR ERROR", cerr)
		return ui.Script{}
	}
	if len(cols) < 1 {
		return ui.Script{}
	}
	color := cols[0].Color
	t.Color = fmt.Sprintf("#%02X%02X%02X", color.R, color.G, color.B)
	return ui.Script{}
}

func (s *Tenant) ObjectMenu(menuID string) []menu.Item {
	switch menuID {
	case "explorer/context":
		return []menu.Item{
			{Cmd: "zt100.new-app", Label: "New App", Icon: "plus"},
		}
	default:
		return []menu.Item{}
	}
}
