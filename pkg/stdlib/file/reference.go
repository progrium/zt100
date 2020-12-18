package file

import (
	"log"
	"os"
	"path"

	"github.com/progrium/zt100/pkg/manifold"
	"github.com/progrium/zt100/pkg/misc/notify"
)

type Reference struct {
	Filepath   string
	SyncName   bool
	SyncDelete bool

	object manifold.Object
}

func (c *Reference) Mounted(obj manifold.Object) error {
	c.object = obj
	c.SyncName = true
	notify.Observe(obj, notify.Func(func(event interface{}) error {
		change, ok := event.(manifold.ObjectChange)
		if !ok || change.Path != "::Name" || obj.ID() != change.Object.ID() {
			return nil
		}
		newName := change.New.(string)
		base := path.Base(c.Filepath)
		dir := path.Dir(c.Filepath)
		if c.SyncName && base != newName {
			newPath := path.Join(dir, newName)
			if err := os.Rename(c.Filepath, newPath); err != nil {
				log.Println(err)
				return nil
			}
			c.Filepath = newPath
		}
		return nil
	}))
	return nil
}
