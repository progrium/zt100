package http

import (
	"context"
	"errors"
	"fmt"
	"net"
	"net/http"
	"strings"

	"github.com/progrium/zt100/pkg/core/service"
	"github.com/progrium/zt100/pkg/ui"
	"github.com/skratchdot/open-golang/open"
)

type Middleware interface {
	Middleware() func(http.Handler) http.Handler
}

type Server struct {
	http.Server

	Listener net.Listener
	Handler  http.Handler

	Services *service.Framework `tractor:"hidden"`
}

func (c *Server) ComponentEnable() {
	if c.Listener == nil || c.Handler == nil {
		return
	}

	c.Server = http.Server{
		Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			c.Handler.ServeHTTP(w, r)
		}),
	}

	c.Services.Start(c)
}

func (c *Server) OpenInBrowser(path ...string) ui.Script {
	open.Run(fmt.Sprintf("http://%s/%s", c.Listener.Addr().String(), strings.Join(path, "/")))
	return ui.JS("")
}

func (c *Server) Serve(ctx context.Context) {
	if err := c.Server.Serve(c.Listener); err != nil {
		if !errors.Is(err, http.ErrServerClosed) && !IsErrClosed(err) {
			panic(err)
		}
	}
}

func (c *Server) Shutdown(ctx context.Context) error {
	return c.Server.Shutdown(ctx)
}

func (c *Server) ComponentDisable() {
	if c.Server.Handler != nil {
		if err := c.Services.Stop(c); err != nil {
			panic(err)
		}
	}
}

func IsErrClosed(err error) bool {
	// TODO go1.16: use net.ErrClosed
	return strings.Contains(err.Error(), "use of closed network connection")
}
