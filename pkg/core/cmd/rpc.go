package cmd

import (
	"github.com/manifold/qtalk/golang/rpc"
)

func (f *Framework) ContributeRPC(mux *rpc.RespondMux) {
	mux.Bind("cmd.exec", func(id string, params map[string]interface{}) (interface{}, error) {
		// TODO: pass along context
		return f.ExecuteCommand(id, params)
	})
}
