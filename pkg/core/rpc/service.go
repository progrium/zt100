package rpc

import (
	"context"
	"net/http"

	L "github.com/progrium/zt100/pkg/misc/logging"

	"github.com/manifold/qtalk/golang/mux"
	"github.com/manifold/qtalk/golang/rpc"
	"golang.org/x/net/websocket"
)

func (f *Framework) Serve(ctx context.Context) {
	L.Debugf(f.Log, "listening at ws://%s", f.l.Addr().String())
	if err := f.s.Serve(f.l); err != nil {
		if err == http.ErrServerClosed {
			return
		}
		L.Error(f.Log, err)
	}
}

func (f *Framework) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if f.h != nil && r.Header.Get("Upgrade") != "websocket" {
		f.h.ServeHTTP(w, r)
		return
	}
	websocket.Handler(func(conn *websocket.Conn) {
		conn.PayloadType = websocket.BinaryFrame
		sess := mux.NewSession(r.Context(), conn)
		L.Debug(f.Log, "new session")
		srv := &rpc.Server{
			Mux: f.mux,
		}
		go srv.Respond(sess)
		sess.Wait()
	}).ServeHTTP(w, r)
}
