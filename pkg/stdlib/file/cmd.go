package file

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"reflect"

	"github.com/progrium/zt100/pkg/core/cmd"
)

func (c *Explorer) ContributeCommands(cmds *cmd.Registry) {
	cmds.Register(cmd.Definition{
		ID:       "file.new",
		Label:    "New File",
		Category: "Files",
		Desc:     "Creates a new file",
		Run: func(params struct {
			ID   string
			Name string
		}) {
			if params.ID == "" || params.Name == "" {
				return
			}
			obj := c.object.Root().FindID(params.ID)
			if obj == nil {
				return
			}
			path := &Path{}
			rv := reflect.ValueOf(path)
			obj.ValueTo(rv)
			if path.Filepath == "" {
				return
			}
			f, err := os.Create(fmt.Sprintf("%s/%s", path.Filepath, params.Name))
			if err != nil {
				log.Println(err)
				return
			}
			f.Close()
		},
	})

	cmds.Register(cmd.Definition{
		ID:       "file.mkdir",
		Label:    "New Folder",
		Category: "Files",
		Desc:     "Creates a new folder",
		Run: func(params struct {
			ID   string
			Name string
		}) {
			if params.ID == "" || params.Name == "" {
				return
			}
			obj := c.object.Root().FindID(params.ID)
			if obj == nil {
				return
			}
			path := &Path{}
			rv := reflect.ValueOf(path)
			obj.ValueTo(rv)
			if path.Filepath == "" {
				return
			}
			err := os.Mkdir(fmt.Sprintf("%s/%s", path.Filepath, params.Name), 0755)
			if err != nil {
				log.Println(err)
				return
			}
		},
	})

	cmds.Register(cmd.Definition{
		ID:       "file.show",
		Label:    "Show File",
		Category: "Files",
		Desc:     "Shows a file or folder in Finder",
		Run: func(params struct {
			ID string
		}) {
			if params.ID == "" {
				return
			}
			obj := c.object.Root().FindID(params.ID)
			if obj == nil {
				return
			}

			path := &Path{}
			rv := reflect.ValueOf(path)
			obj.ValueTo(rv)
			if path.Filepath != "" {
				c := exec.Command("open", path.Filepath)
				c.Start()
				return
			}

			ref := &Reference{}
			rv = reflect.ValueOf(ref)
			obj.ValueTo(rv)
			if ref.Filepath != "" {
				c := exec.Command("open", ref.Filepath)
				c.Start()
				return
			}
		},
	})
}
