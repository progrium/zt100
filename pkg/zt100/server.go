package zt100

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"path"
	"path/filepath"
	"runtime"

	"github.com/gorilla/mux"
	"github.com/spf13/afero"

	"github.com/manifold/tractor/pkg/core/cmd"
	"github.com/manifold/tractor/pkg/manifold"
	"github.com/manifold/tractor/pkg/manifold/comutil"
	"github.com/manifold/tractor/pkg/stdlib/file"
	mk "github.com/manifold/tractor/pkg/stdlib/file/make"
	"github.com/manifold/tractor/pkg/ui/menu"
)

var Message = "Hello zt100"

type Server struct {
	Cmds *cmd.Framework `tractor:"hidden"`

	StaticHandler http.Handler
	Builder       *mk.JSX
	Blocks        string
	object        manifold.Object
}

func (s *Server) ObjectMenu(menuID string) []menu.Item {
	switch menuID {
	case "explorer/context":
		return []menu.Item{
			{Cmd: "zt100.new-tenant", Label: "New Tenant", Icon: "plus-square"},
		}
	default:
		return []menu.Item{}
	}
}

func (s *Server) Mounted(obj manifold.Object) error {
	s.object = obj
	return nil
}

func (s *Server) Tenants() (tenants []*Tenant) {
	for _, obj := range s.object.Children() {
		for _, com := range comutil.Enabled(obj) {
			if com.Name() != "zt100.Tenant" {
				continue
			}
			t, ok := com.Pointer().(*Tenant)
			if ok {
				tenants = append(tenants, t)
			}
		}
	}
	return tenants
}

func (s *Server) Lookup(tenantName, appName, pageName, sectionKey string) (tenant manifold.Object, app manifold.Object, page manifold.Object, sections []*Section) {
	//vars := mux.Vars(r)
	for _, c := range s.object.Children() {
		if c.Name() == tenantName {
			tenant = c
			break
		}
	}
	if tenant == nil {
		return nil, nil, nil, nil
	}
	for _, a := range tenant.Children() {
		if a.Name() == appName {
			app = a
			break
		}
	}
	if app == nil {
		return tenant, nil, nil, nil
	}
	if pageName == "" {
		pageName = "index"
	}
	for _, s := range app.Children() {
		if s.Name() == pageName {
			page = s
			break
		}
	}
	if page == nil {
		return tenant, app, nil, nil
	}
	for _, c := range comutil.Enabled(page) {
		if c.Name() != "zt100.Section" {
			continue
		}
		if sectionKey != "" && c.Key() != sectionKey {
			continue
		}
		b, ok := c.Pointer().(*Section)
		if ok {
			sections = append(sections, b)
		}
	}
	return tenant, app, page, sections
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	log.Println(r.Method, r.URL.Path)
	m := mux.NewRouter()
	m.HandleFunc("/blocks/{block}.js", func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		var block manifold.Object
		for _, c := range comutil.MustLookup(s.Blocks).Children() {
			if c.Name() == fmt.Sprintf("%s.jsx", vars["block"]) {
				block = c
				break
			}
		}
		if block == nil {
			http.Error(w, "block not found", http.StatusNotFound)
			return
		}
		//b := block.Component("zt100.Block").Pointer().(*Block)
		fr := block.Component("file.Reference").Pointer().(*file.Reference)
		b, berr := ioutil.ReadFile(fr.Filepath)
		if berr != nil {
			http.Error(w, berr.Error(), http.StatusServiceUnavailable)
		}

		ext := filepath.Ext(block.Name())
		name := block.Name()[:len(block.Name())-len(ext)]

		fs := afero.NewMemMapFs()
		afero.WriteFile(fs, fmt.Sprintf("%s.jsx", name), b, 0644)
		if err := s.Builder.Build(fs, fmt.Sprintf("%s.js", name), fmt.Sprintf("%s.jsx", name)); err != nil {
			http.Error(w, err.Error(), http.StatusServiceUnavailable)
		}
		d, err := afero.ReadFile(fs, fmt.Sprintf("%s.js", name))
		if err != nil {
			http.Error(w, err.Error(), http.StatusServiceUnavailable)
		}
		w.Header().Set("content-type", "text/javascript")
		w.Write(d)
	})
	m.HandleFunc("/c/{cmd}", func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		defer r.Body.Close()
		dec := json.NewDecoder(r.Body)
		var v map[string]interface{}
		if err := dec.Decode(&v); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		_, err := s.Cmds.ExecuteCommand(vars["cmd"], v)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	})
	m.HandleFunc("/t/{tenant}/{app}", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, fmt.Sprintf("%s/index", r.URL.Path), 302)
	})
	m.HandleFunc("/t/{tenant}/{app}/{page}", func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		_, _, page, _ := s.Lookup(vars["tenant"], vars["app"], vars["page"], "")
		if page == nil {
			http.Error(w, "page not found", http.StatusNotFound)
			return
		}
		page.Component("zt100.Page").Pointer().(*Page).ServeHTTP(w, r)
	})
	m.HandleFunc("/e/{tenant}", func(w http.ResponseWriter, r *http.Request) {
		ts, err := template.ParseFiles(
			"./tmpl/editor.page.html",
			"./tmpl/base.layout.html",
		)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		err = ts.Execute(w, nil)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	})
	m.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		ts, err := template.ParseFiles(
			"./tmpl/index.page.html",
			"./tmpl/base.layout.html",
		)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		err = ts.Execute(w, map[string]interface{}{
			"Tenants": s.Tenants(),
		})
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	})
	m.PathPrefix("/static").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir(localpath("../../static")))))
	m.PathPrefix("/").Handler(s.StaticHandler)
	m.ServeHTTP(w, r)
}

func localpath(subpath string) string {
	_, filename, _, _ := runtime.Caller(1)
	dir, _ := filepath.Abs(path.Join(path.Dir(filename), subpath))
	return dir
}
