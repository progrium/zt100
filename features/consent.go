package features

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"path/filepath"

	"github.com/gorilla/sessions"
	"github.com/progrium/zt100/pkg/feature"
	"github.com/progrium/zt100/pkg/zt100"
)

type ConsentFeature struct {
	Sessions sessions.Store `json:"-"`
	Server   *zt100.Server  `json:"-"`
}

func (f *ConsentFeature) Initialize() {
	f.Sessions = sessions.NewCookieStore([]byte("zt100-session-store"))
}

func (f *ConsentFeature) Flag() feature.Flag {
	return feature.Flag{
		Name: "login:consent",
		Desc: "Require consent first login",
	}
}

func (f *ConsentFeature) LoginRedirect(nextURL string, r *http.Request) string {
	ctx := zt100.FromContext(r.Context())
	if ctx.HasFeature("login:consent") {
		return filepath.Join(filepath.Dir(nextURL), fmt.Sprintf("consent?then=%s", nextURL))
	}
	return nextURL
}

func (f *ConsentFeature) HandlePages() []string {
	return []string{"consent"}
}

func (f *ConsentFeature) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch filepath.Base(r.URL.Path) {
	case "consent":
		f.consent(w, r)
	case "reset":
		f.reset(w, r)
	default:
		http.NotFound(w, r)
	}
}

func (f *ConsentFeature) reset(w http.ResponseWriter, r *http.Request) {
	session, err := f.Sessions.Get(r, "zt100-session-store")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	delete(session.Values, "consented")
	session.Save(r, w)
	io.WriteString(w, "consent reset for session")
}

func (f *ConsentFeature) consent(w http.ResponseWriter, r *http.Request) {
	session, err := f.Sessions.Get(r, "zt100-session-store")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if r.Method == "POST" {
		session.Values["consented"] = "true"
		session.Save(r, w)
	}

	nextURL := r.URL.Query().Get("then")
	if nextURL == "" {
		u, _ := url.Parse(r.Referer())
		nextURL = u.Path
	}

	_, ok := session.Values["consented"]
	if ok {
		http.Redirect(w, r, nextURL, http.StatusTemporaryRedirect)
		return
	}

	ctx := zt100.FromContext(r.Context())
	if ctx.Page != nil {
		ctx.Page.ServeHTTP(w, r)
		return
	}
	http.NotFound(w, r)
}
