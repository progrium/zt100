package file

import (
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/progrium/watcher"
	"github.com/progrium/zt100/pkg/manifold"
	"github.com/progrium/zt100/pkg/manifold/library"
	"github.com/progrium/zt100/pkg/manifold/object"
	"github.com/progrium/zt100/pkg/ui/menu"
)

type Explorer struct {
	Path *Path

	object  manifold.Object
	unwatch func()
}

func (c *Explorer) Mounted(obj manifold.Object) error {
	c.object = obj
	return nil
}

func (c *Explorer) ComponentEnable() {
	if !c.exists() {
		return
	}
	// var err error
	// c.unwatch, err = c.Path.Watch(".", c.fsWatch)
	// if err != nil {
	// 	log.Println(err)
	// }
}

func (c *Explorer) ComponentDisable() {
	if c.unwatch != nil {
		c.unwatch()
	}
}

func (c *Explorer) fsWatch(event watcher.Event) {
	switch event.Op {
	case watcher.Create:
		obj := object.New(event.Name())
		if event.IsDir() {
			obj.AppendComponent(library.NewComponent("file.Path", &Path{
				Filepath: event.Path,
			}, ""))
			obj.AppendComponent(library.NewComponent("file.Explorer", &Explorer{}, ""))
		} else {
			obj.AppendComponent(library.NewComponent("file.Reference", &Reference{
				Filepath: event.Path,
			}, ""))
		}
		c.object.AppendChild(obj)
		// log.Println("CREATE", event.Name())
	case watcher.Remove:
		obj := c.object.FindChild(event.Name())
		if obj != nil {
			c.object.RemoveChild(obj)
			// log.Println("REMOVE", event.Name())
		}
	case watcher.Rename:
		obj := c.object.FindChild(event.Name())
		if obj != nil {
			obj.SetName(filepath.Base(event.Path))
			if cc := obj.Component("file.Explorer"); cc != nil {
				cc.Pointer().(*Explorer).Path.Filepath = event.Path
			}
			// log.Println("RENAME", event.Name())
		}
	case watcher.Move:
		obj := c.object.FindChild(event.Name())
		if obj != nil {
			obj.SetName(filepath.Base(event.Path))
			if c := obj.Component("file.Explorer"); c != nil {
				c.Pointer().(*Explorer).Path.Filepath = event.Path
			}
			// log.Println("MOVE", event.Name())
		}
	}
}

func (c *Explorer) exists() bool {
	if c.Path == nil {
		return false
	}
	if c.Path.Filepath == "" {
		return false
	}
	if _, err := os.Stat(c.Path.Filepath); os.IsNotExist(err) {
		return false
	}
	return true
}

func (c *Explorer) LazyChildren() bool {
	return true
}

func (c *Explorer) ObjectMenu(menuID string) []menu.Item {
	switch menuID {
	case "explorer/context":
		return []menu.Item{
			{Cmd: "file.new", Label: "New File", Icon: "file"},
			{Cmd: "file.mkdir", Label: "New Folder", Icon: "folder"},
			{Cmd: "file.show", Label: "Open in Finder"},
		}
	default:
		return []menu.Item{}
	}
}

func (c *Explorer) ChildNodes() (objs []manifold.Object) {
	if !c.exists() {
		return
	}
	fi, err := ioutil.ReadDir(c.Path.Filepath)
	if err != nil {
		panic(err)
	}
	for _, f := range fi {
		if strings.HasPrefix(f.Name(), ".") {
			continue
		}
		obj := object.New(f.Name())
		if f.IsDir() {
			obj.AppendComponent(library.NewComponent("file.Path", &Path{
				Filepath: path.Join(c.Path.Filepath, f.Name()),
			}, ""))
			obj.AppendComponent(library.NewComponent("file.Explorer", &Explorer{}, ""))
		} else {
			obj.AppendComponent(library.NewComponent("file.Reference", &Reference{
				Filepath: path.Join(c.Path.Filepath, f.Name()),
			}, ""))
		}
		objs = append(objs, obj)
	}
	return
}
