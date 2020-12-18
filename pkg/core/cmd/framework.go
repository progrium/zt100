package cmd

import (
	"time"

	L "github.com/progrium/zt100/pkg/misc/logging"
)

type Contributor interface {
	ContributeCommands(cmds *Registry)
}

type Observer interface {
	CommandExecuted(cmd Definition, params map[string]interface{})
}

type Framework struct {
	Contributors []Contributor
	Observers    []Observer
	Log          L.Logger

	cmds *Registry
	view *ViewState
}

func (f *Framework) ExecuteCommand(cmdID string, params map[string]interface{}) (result interface{}, err error) {
	start := time.Now()
	cmd := f.cmds.Get(cmdID)
	result, err = f.cmds.Execute(cmdID, params)
	if err != nil {
		L.Infof(f.Log, "cmd:err %s %+v %s", cmdID, params, err.Error())
		return
	}

	L.Infof(f.Log, "cmd: %s %+v %s", cmdID, params, time.Since(start).Truncate(time.Millisecond))

	for _, obs := range f.Observers {
		obs.CommandExecuted(cmd, params)
	}

	return
}
