package zt100

import (
	"fmt"
	"image"
	"log"
	"net/http"

	_ "image/png"

	"github.com/EdlinOrg/prominentcolor"
	"github.com/progrium/zt100/pkg/manifold"
	"github.com/progrium/zt100/pkg/manifold/comutil"
	"github.com/progrium/zt100/pkg/ui"
	"github.com/progrium/zt100/pkg/ui/menu"
)

type Demo struct {
	Domain     string
	Color      string `field:"color"`
	OID        string `tractor:"hidden"`
	Vertical   string
	Features   []string
	ObjectName string `json:"Name"`

	object manifold.Object
}

func (s *Demo) Mounted(obj manifold.Object) error {
	s.object = obj
	s.OID = obj.ID()
	s.ObjectName = obj.Name()
	return nil
}

func (s *Demo) App(name string) *App {
	for _, b := range s.Apps() {
		if b.Name == name {
			return b
		}
	}
	return nil
}

func (s *Demo) Apps() (apps []*App) {
	for _, obj := range s.object.Children() {
		for _, com := range comutil.Enabled(obj) {
			if com.Name() != "zt100.App" {
				continue
			}
			t, ok := com.Pointer().(*App)
			if ok {
				apps = append(apps, t)
			}
		}
	}
	return apps
}

func (t *Demo) Name() string {
	return t.object.Name()
}

func (t *Demo) GetColor() ui.Script {
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

func (s *Demo) ObjectMenu(menuID string) []menu.Item {
	switch menuID {
	case "explorer/context":
		return []menu.Item{
			{Cmd: "zt100.new-app", Label: "New App", Icon: "plus"},
		}
	default:
		return []menu.Item{}
	}
}
