package zt100

import (
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/gorilla/mux"

	"github.com/manifold/tractor/pkg/core/cmd"
	"github.com/manifold/tractor/pkg/core/obj"
	"github.com/manifold/tractor/pkg/manifold"
	"github.com/manifold/tractor/pkg/manifold/comutil"
	mk "github.com/manifold/tractor/pkg/stdlib/file/make"
	"github.com/manifold/tractor/pkg/ui/menu"
)

var Message = "Hello zt100"

type Server struct {
	Cmds    *cmd.Framework `tractor:"hidden"`
	Objects *obj.Service   `tractor:"hidden"`

	StaticHandler http.Handler
	Builder       *mk.JSX
	Blocks        string
	object        manifold.Object
}

func (s *Server) ObjectMenu(menuID string) []menu.Item {
	switch menuID {
	case "explorer/context":
		return []menu.Item{
			{Cmd: "zt100.new-prospect", Label: "New Prospect", Icon: "plus-square"},
		}
	default:
		return []menu.Item{}
	}
}

func (s *Server) Mounted(obj manifold.Object) error {
	s.object = obj
	return nil
}

func (s *Server) Prospects() (prospects []*Prospect) {
	for _, obj := range s.object.Children() {
		for _, com := range comutil.Enabled(obj) {
			if com.Name() != "zt100.Prospect" {
				continue
			}
			t, ok := com.Pointer().(*Prospect)
			if ok {
				prospects = append(prospects, t)
			}
		}
	}
	return prospects
}

func (s *Server) Lookup(prospectName, appName, pageName, sectionKey string) (prospect manifold.Object, app manifold.Object, page manifold.Object, sections []*Section) {
	//vars := mux.Vars(r)
	for _, c := range s.object.Children() {
		if c.Name() == prospectName {
			prospect = c
			break
		}
	}
	if prospect == nil {
		return nil, nil, nil, nil
	}
	for _, a := range prospect.Children() {
		if a.Name() == appName {
			app = a
			break
		}
	}
	if app == nil {
		return prospect, nil, nil, nil
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
		return prospect, app, nil, nil
	}
	for _, c := range comutil.Enabled(page) {
		if c.Name() != "zt100.Section" {
			continue
		}
		if sectionKey != "" && c.ID() != sectionKey {
			continue
		}
		b, ok := c.Pointer().(*Section)
		if ok {
			sections = append(sections, b)
		}
	}
	return prospect, app, page, sections
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	log.Println(r.Method, r.URL.Path)
	m := mux.NewRouter()

	m.HandleFunc("/c/{cmd}", func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		var v map[string]interface{}
		defer r.Body.Close()

		if vars["cmd"] == "zt100.new-section" && r.URL.Query().Get("upload") == "1" {
			if err := r.ParseMultipartForm(100 << 20); err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
			mf, mfh, _ := r.FormFile("file")
			// if err != nil {
			// 	http.Error(w, err.Error(), http.StatusBadRequest)
			// 	return
			// }
			v = map[string]interface{}{
				"PageID":    r.FormValue("PageID"),
				"BlockID":   r.FormValue("BlockID"),
				"Image":     mf,
				"ImageSize": mfh.Size,
			}
		} else {
			dec := json.NewDecoder(r.Body)
			if err := dec.Decode(&v); err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
		}

		_, err := s.Cmds.ExecuteCommand(vars["cmd"], v)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	})

	m.HandleFunc("/t/{prospect}/{app}", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, fmt.Sprintf("%s/index", r.URL.Path), 302)
	})
	m.HandleFunc("/t/{prospect}/{app}/{page}", func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		_, _, page, _ := s.Lookup(vars["prospect"], vars["app"], vars["page"], "")
		if page == nil {
			http.Error(w, "page not found", http.StatusNotFound)
			return
		}
		page.Component("zt100.Page").Pointer().(*Page).ServeHTTP(w, r)
	})

	m.HandleFunc("/e/{prospect}.json", func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		var prospect *Prospect
		for _, t := range s.Prospects() {
			if t.Name() == vars["prospect"] {
				prospect = t
			}
		}
		if prospect == nil {
			http.NotFound(w, r)
			return
		}
		w.Header().Set("content-type", "text/javascript")
		subject := strings.TrimPrefix(r.URL.RawQuery, "?")
		if subject == "" {
			subject = "/"
			apps := prospect.Apps()
			if len(apps) > 0 {
				subject = fmt.Sprintf("%s/index", apps[0].Name)
			}
		}
		parts := strings.Split(subject, "/")
		appname := parts[0]
		pagename := parts[1]

		var sections []*Section
		var apps []*App
		var pageOID string
		pages := make(map[string][]*Page)
		for _, a := range prospect.Apps() {
			apps = append(apps, a)
			pages[a.Name] = a.Pages()
			if a.Name == appname {
				for _, p := range a.Pages() {
					if p.Name == pagename {
						sections = p.Sections()
						pageOID = p.OID
						break
					}
				}
			}
		}

		var blocks []*Block
		for _, c := range comutil.MustLookup(s.Blocks).Children() {
			b, ok := c.Component("zt100.Block").Pointer().(*Block)
			if !ok {
				continue
			}
			blocks = append(blocks, b)
		}

		data := map[string]interface{}{
			"ProspectName":   prospect.Name(),
			"ProspectDomain": prospect.Domain,
			"ProspectOID":    prospect.OID,
			"AppName":        appname,
			"PageName":       pagename,
			"PageOID":        pageOID,
			"Apps":           apps,
			"Pages":          pages,
			"Sections":       sections,
			"Blocks":         blocks,
		}
		enc := json.NewEncoder(w)
		if err := enc.Encode(data); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	})
	m.HandleFunc("/e/{prospect}", func(w http.ResponseWriter, r *http.Request) {
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
			"Prospects": s.Prospects(),
		})
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	})
	m.PathPrefix("/blocks").Handler(http.StripPrefix("/blocks/", http.FileServer(http.Dir(localpath("../../blocks")))))
	m.PathPrefix("/static").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir(localpath("../../static")))))
	m.PathPrefix("/uploads").Handler(http.StripPrefix("/uploads/", http.FileServer(http.Dir(localpath("../../local/uploads")))))
	//m.PathPrefix("/").Handler(s.StaticHandler)
	m.ServeHTTP(w, r)
}

func localpath(subpath string) string {
	_, filename, _, _ := runtime.Caller(1)
	dir, _ := filepath.Abs(path.Join(path.Dir(filename), subpath))
	return dir
}

func fileExists(path string) bool {
	info, err := os.Stat(path)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}
