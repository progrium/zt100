package rpc

import (
	"net"
	"net/http"

	"github.com/manifold/qtalk/golang/rpc"
	L "github.com/progrium/zt100/pkg/misc/logging"
)

type Contributor interface {
	ContributeRPC(mux *rpc.RespondMux)
}

type Framework struct {
	Contributors []Contributor
	Log          L.Logger

	mux *rpc.RespondMux
	l   net.Listener
	s   *http.Server
	h   http.Handler
}

func (f *Framework) Addr() string {
	return f.l.Addr().String()
}

func (f *Framework) SetHandler(h http.Handler) {
	f.h = h
}
