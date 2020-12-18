package zt100

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"path"
	"path/filepath"
	"runtime"
	"text/template"

	"github.com/gorilla/mux"
	"github.com/progrium/zt100/pkg/feature"
	"github.com/spf13/afero"

	"github.com/progrium/zt100/pkg/core/cmd"
	"github.com/progrium/zt100/pkg/core/obj"
	"github.com/progrium/zt100/pkg/manifold"
	"github.com/progrium/zt100/pkg/manifold/comutil"
	"github.com/progrium/zt100/pkg/ui/menu"
)

type PageHandler interface {
	HandlePages() []string
	ServeHTTP(w http.ResponseWriter, r *http.Request)
}

type Server struct {
	Cmds     *cmd.Framework     `tractor:"hidden" json:"-"`
	Objects  *obj.Service       `tractor:"hidden" json:"-"`
	Template *template.Template `json:"-"`

	Features []feature.Feature `json:"-"`

	object manifold.Object
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	log.Println(r.Method, r.URL.Path)

	m := mux.NewRouter()
	m.HandleFunc("/", s.index)
	m.HandleFunc("/cmd/{cmd}", s.cmd)
	m.HandleFunc("/edit/{demo}/{app}/{page}", s.edit)
	m.HandleFunc("/data/{demo}/{app}/{page}.json", s.data)
	m.HandleFunc("/preview/{demo}/{app}/{page}", s.preview)
	m.HandleFunc("/preview/{demo}/{app}/{page}/{block}.js", s.block)
	m.PathPrefix("/feature/{feature}").HandlerFunc(s.feature)
	m.PathPrefix("/static").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir(localpath("../../static")))))
	//m.PathPrefix("/uploads").Handler(http.StripPrefix("/uploads/", http.FileServer(http.Dir(localpath("../../local/uploads")))))
	m.Handle("/preview/{demo}/{app}", http.RedirectHandler(fmt.Sprintf("%s/index", r.URL.Path), 302))
	m.ServeHTTP(w, r)
}

func (s *Server) Initialize() {
	s.Template = template.Must(template.ParseGlob("./tmpl/*"))
	feature.Load(s)
}

func (s *Server) Mounted(obj manifold.Object) error {
	s.object = obj
	return nil
}

func (s *Server) Demo(name string) *Demo {
	for _, b := range s.Demos() {
		if b.Name() == name {
			return b
		}
	}
	return nil
}

func (s *Server) Demos() (demos []*Demo) {
	for _, obj := range s.object.Children() {
		for _, com := range comutil.Enabled(obj) {
			if com.Name() != "zt100.Demo" {
				continue
			}
			t, ok := com.Pointer().(*Demo)
			if ok {
				demos = append(demos, t)
			}
		}
	}
	return demos
}

func (s *Server) Block(name string) *Block {
	for _, b := range s.Blocks() {
		if b.Name == name {
			return b
		}
	}
	return nil
}

func (s *Server) Blocks() (blocks []*Block) {
	fs := afero.Afero{Fs: afero.NewOsFs()}
	path := localpath("../../blocks")
	dir, err := fs.ReadDir(path)
	if err != nil {
		panic(err)
	}
	for _, fi := range dir {
		if fi.IsDir() {
			continue
		}
		if fi.Name() == "_template.js" {
			continue
		}
		b, err := fs.ReadFile(filepath.Join(path, fi.Name()))
		if err != nil {
			log.Println(err)
			continue
		}
		ext := filepath.Ext(fi.Name())
		blocks = append(blocks, &Block{
			Name:   fi.Name()[:len(fi.Name())-len(ext)],
			Source: b,
		})
	}
	return blocks
}

func (s *Server) Feature(flag string) feature.Feature {
	for _, f := range s.Features {
		if f.Flag().Name == flag {
			return f
		}
	}
	return nil
}

func (s *Server) feature(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	for _, f := range s.Features {
		h, ok := f.(http.Handler)
		if ok && f.Flag().Name == vars["feature"] {
			h.ServeHTTP(w, r)
			return
		}
	}
	http.NotFound(w, r)
}

func (s *Server) cmd(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	var params map[string]interface{}

	if vars["cmd"] == "zt100.block.attach-image" && r.URL.Query().Get("upload") == "1" {
		// special case for this cmd that uses multipart form
		if err := r.ParseMultipartForm(100 << 20); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		mf, mfh, _ := r.FormFile("file")
		params = map[string]interface{}{
			"PageID":    r.FormValue("PageID"),
			"BlockID":   r.FormValue("BlockID"),
			"Image":     mf,
			"ImageSize": mfh.Size,
		}
	} else {
		// default case
		dec := json.NewDecoder(r.Body)
		defer r.Body.Close()
		if err := dec.Decode(&params); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
	}

	_, err := s.Cmds.ExecuteCommand(vars["cmd"], params)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (s *Server) data(w http.ResponseWriter, r *http.Request) {
	ctx := LoadContext(s, r)

	if ctx.Demo == nil {
		fmt.Printf("%#v\n", ctx)
		http.NotFound(w, r)
		return
	}

	w.Header().Set("content-type", "application/json")
	enc := json.NewEncoder(w)
	if err := enc.Encode(ctx); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (s *Server) preview(w http.ResponseWriter, r *http.Request) {
	ctx := LoadContext(s, r)
	r = r.WithContext(context.WithValue(r.Context(), "data", ctx))

	vars := mux.Vars(r)
	for _, f := range s.Features {
		ph, ok := f.(PageHandler)
		flag := f.Flag().Name
		if ok && ctx.HasFeature(flag) {
			for _, page := range ph.HandlePages() {
				if page == vars["page"] {
					ph.ServeHTTP(w, r)
					return
				}
			}
		}
	}

	if ctx.Page == nil {
		http.NotFound(w, r)
		return
	}

	ctx.Page.ServeHTTP(w, r)
}

func (s *Server) block(w http.ResponseWriter, r *http.Request) {
	ctx := LoadContext(s, r)
	if ctx.Block == nil {
		http.NotFound(w, r)
		return
	}
	ctx.Block.ServeHTTP(w, r.WithContext(context.WithValue(r.Context(), "data", ctx)))
}

func (s *Server) edit(w http.ResponseWriter, r *http.Request) {
	err := s.Template.ExecuteTemplate(w, "edit.html", nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (s *Server) index(w http.ResponseWriter, r *http.Request) {
	//b, _ := json.MarshalIndent(getProfileData(r), "", "  ")
	err := s.Template.ExecuteTemplate(w, "index.html", map[string]interface{}{
		"Demos":    s.Demos(),
		"OID":      s.object.ID(),
		"Features": feature.FlattenFlags(s.Features),
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (s *Server) ObjectMenu(menuID string) []menu.Item {
	switch menuID {
	case "explorer/context":
		return []menu.Item{
			{Cmd: "zt100.new-demo", Label: "New Demo", Icon: "plus-square"},
		}
	default:
		return []menu.Item{}
	}
}

func localpath(subpath string) string {
	_, filename, _, _ := runtime.Caller(1)
	dir, _ := filepath.Abs(path.Join(path.Dir(filename), subpath))
	return dir
}
