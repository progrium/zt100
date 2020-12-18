package state

import (
	"github.com/manifold/qtalk/golang/rpc"
	"github.com/progrium/zt100/pkg/misc/jsonpointer"
)

func (f *Framework) ContributeRPC(mux *rpc.RespondMux) {
	mux.Bind("state.sync", rpc.HandlerFunc(func(r rpc.Responder, c *rpc.Call) {
		f.clients[c.Caller] = struct{}{}
		r.Return(f.state)
	}))

	mux.Bind("state.update", func(state map[string]map[string]interface{}) {
		f.mu.Lock()
		defer f.mu.Unlock()
		for ns, kv := range state {
			s, ok := f.state[ns]
			if !ok {
				continue
			}
			for k, v := range kv {
				if err := jsonpointer.SetReflect(s, k, v); err != nil {
					panic(err)
				}
			}
		}
	})
}
