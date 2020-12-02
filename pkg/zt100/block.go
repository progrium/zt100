package zt100

import (
	"fmt"
	"io/ioutil"
	"log"

	"github.com/manifold/tractor/pkg/core/fswatch"
	"github.com/manifold/tractor/pkg/manifold"
	"github.com/manifold/tractor/pkg/stdlib/file"
	"github.com/manifold/tractor/pkg/ui"
	"github.com/progrium/watcher"
)

type Block struct {
	source []byte `tractor:"hidden"`

	Watcher *fswatch.Service `tractor:"hidden"`
	watch   *fswatch.Watch
	fileref *file.Reference
	object  manifold.Object
}

func (c *Block) Mounted(obj manifold.Object) error {
	c.object = obj
	c.fileref = c.object.Component("file.Reference").Pointer().(*file.Reference)
	return nil
}

func (c *Block) Initialize() {
	if len(c.source) == 0 {
		c.source = []byte(`
export default function({attrs}) {
	return (
		<div></div>
	)
}
`)
	}
}

func (c *Block) Source() []byte {
	b, err := ioutil.ReadFile(c.fileref.Filepath)
	if err != nil {
		panic(err)
	}
	return b
}

func (c *Block) fileUpdate(event watcher.Event) {
	if event.Op == watcher.Write {
		var err error
		c.source, err = ioutil.ReadFile(c.watch.Path)
		if err != nil {
			log.Println(err)
		}

	}
}

func (c *Block) EditSource() ui.Script {
	if c.object.Name() == "" {
		return ui.Script{}
	}
	return ui.Script{Src: fmt.Sprintf(`window.T.exec("editor.open", {Filename: "%s"})`, c.fileref.Filepath)}
}

func (c *Block) ComponentDisable() {
	if c.watch != nil {
		c.Watcher.Unwatch(c.watch)
	}
}
