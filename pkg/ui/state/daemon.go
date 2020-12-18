package state

import (
	"context"

	"github.com/manifold/qtalk/golang/rpc"
)

func (f *Framework) InitializeDaemon() (err error) {
	f.clients = make(map[rpc.Caller]struct{})
	f.state = make(map[string]interface{})

	for _, c := range f.Contributors {
		n, s := c.InitializeState()
		f.state[n] = s
		if err := c.UpdateState(); err != nil {
			return err
		}
	}

	return nil
}

func (f *Framework) TerminateDaemon(ctx context.Context) error {
	for client := range f.clients {
		client.Call("state.close", nil, nil)
	}
	return nil
}
