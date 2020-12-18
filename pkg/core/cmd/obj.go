package cmd

import "github.com/progrium/zt100/pkg/manifold"

func (f *Framework) MountedComponent(com manifold.Component) {
	c, ok := com.Pointer().(Contributor)
	if ok {
		c.ContributeCommands(f.cmds)
	}
}
