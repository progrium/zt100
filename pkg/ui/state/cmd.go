package state

import "github.com/progrium/zt100/pkg/core/cmd"

func (f *Framework) CommandExecuted(c cmd.Definition, params map[string]interface{}) {
	f.UpdateState() // TODO: handle error
}
