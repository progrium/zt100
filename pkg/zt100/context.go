package zt100

import (
	"context"
	"net/http"

	"github.com/gorilla/mux"
)

type ContextData struct {
	Server *Server

	Demo *Demo
	App  *App
	Page *Page

	Block *Block

	Menu       []MenuItem
	Blocks     []*Block
	Apps       []*App
	Pages      map[string][]*Page
	BaseBlocks []*Block

	Contrib map[string]interface{}
}

func (d *ContextData) HasFeature(flag string) bool {
	if d.Demo == nil {
		return false
	}
	for _, f := range d.Demo.Features {
		if f == flag {
			return true
		}
	}
	return false
}

func FromContext(ctx context.Context) ContextData {
	return ctx.Value("data").(ContextData)
}

func LoadContext(s *Server, r *http.Request) ContextData {
	ctx := ContextData{
		Server:  s,
		Pages:   make(map[string][]*Page),
		Contrib: make(map[string]interface{}),
	}
	vars := mux.Vars(r)

	ctx.BaseBlocks = ctx.Server.Blocks()
	ctx.Demo = s.Demo(vars["demo"])

	if ctx.Demo != nil {
		app := vars["app"]
		if app == "" {
			app = "main"
		}
		ctx.App = ctx.Demo.App(app)
		ctx.Apps = ctx.Demo.Apps()
		for _, app := range ctx.Apps {
			ctx.Pages[app.Name] = app.Pages()
		}
	}

	if ctx.App != nil {
		page := vars["page"]
		if page == "" {
			page = "index"
		}
		ctx.Page = ctx.App.Page(page)
		ctx.Menu = ctx.App.PageMenu()
	}

	if ctx.Page != nil {
		ctx.Blocks = ctx.Page.Blocks()
	}

	block := r.URL.Query().Get("block")
	if block == "" {
		block = vars["block"]
	}
	if ctx.Page != nil && block != "" {
		ctx.Block = ctx.Page.Block(block)
	}

	for _, feat := range s.Features {
		cc, ok := feat.(ContextContributor)
		if ok {
			cc.ContributeContext(&ctx, r)
		}
	}

	return ctx
}

type ContextContributor interface {
	ContributeContext(ctx *ContextData, r *http.Request)
}
