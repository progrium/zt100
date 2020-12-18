package menu

import (
	"context"

	"github.com/manifold/qtalk/golang/rpc"
)

func (f *Framework) ContributeRPC(mux *rpc.RespondMux) {
	mux.Bind("menu.get", func(menuID string, objID string) []Item {
		ctx := context.WithValue(context.Background(), "obj", objID)
		return f.Registry.Get(menuID, ctx)
	})

}
