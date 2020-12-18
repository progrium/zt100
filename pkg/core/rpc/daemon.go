package rpc

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"os"

	"github.com/manifold/qtalk/golang/rpc"
)

func (f *Framework) InitializeDaemon() error {
	port := os.Getenv("WORKSITE_PORT")
	if port == "" {
		port = "0"
	}

	var err error
	f.l, err = net.Listen("tcp", fmt.Sprintf("127.0.0.1:%s", port))
	if err != nil {
		return err
	}
	f.s = &http.Server{
		Handler: f,
	}

	f.mux = rpc.NewRespondMux(rpc.JSONCodec{})

	for _, c := range f.Contributors {
		c.ContributeRPC(f.mux)
	}

	return nil
}

func (f *Framework) TerminateDaemon(ctx context.Context) error {
	return f.s.Shutdown(ctx)
}
