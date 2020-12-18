package features

import (
	"bytes"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"path/filepath"

	"github.com/gorilla/sessions"
	verifier "github.com/okta/okta-jwt-verifier-golang"
	"github.com/progrium/zt100/pkg/feature"
	"github.com/progrium/zt100/pkg/zt100"
)

type LoginFeature struct {
	Nonce    string
	State    string
	Sessions sessions.Store `json:"-"`

	Server *zt100.Server `json:"-"`
}

type LoginRedirector interface {
	LoginRedirect(nextURL string, r *http.Request) string
}

func (f *LoginFeature) LoginRedirect(nextURL string, r *http.Request) string {
	ctx := zt100.FromContext(r.Context())
	if ctx.App.Page("profile") != nil {
		return filepath.Join(filepath.Dir(nextURL), "profile")
	}
	return nextURL
}

func (f *LoginFeature) Initialize() {
	f.Sessions = sessions.NewCookieStore([]byte("zt100-session-store"))
	f.State = "ApplicationState"
	f.Nonce = "NonceNotSetYet"
}

func (f *LoginFeature) Flag() feature.Flag {
	return feature.Flag{
		Name: "login",
		Desc: "Okta Login",
		Subflags: []feature.Flag{
			{
				Name: "login:custom",
				Desc: "Custom, non-hosted login",
			},
		},
	}
}

func (f *LoginFeature) HandlePages() []string {
	return []string{"login", "logout"}
}

func (f *LoginFeature) AppCreated() {

}

func (f *LoginFeature) ContributeContext(ctx *zt100.ContextData, r *http.Request) {
	ctx.Contrib["Auth"] = struct {
		Profile         map[string]string
		IsAuthenticated bool
	}{
		Profile:         f.getProfileData(r),
		IsAuthenticated: f.isAuthenticated(r),
	}
}

func (f *LoginFeature) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch filepath.Base(r.URL.Path) {
	case "login":
		f.login(w, r)
	case "logout":
		f.logout(w, r)
	case "callback":
		f.callback(w, r)
	default:
		http.NotFound(w, r)
	}
}

func (f *LoginFeature) login(w http.ResponseWriter, r *http.Request) {
	session, err := f.Sessions.Get(r, "zt100-session-store")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	nextURL := "/"
	v, ok := session.Values["after_login"]
	if f.isAuthenticated(r) && ok {
		url, ok := v.(string)
		if ok && url != "" {
			nextURL = url
		}
		delete(session.Values, "after_login")
		session.Save(r, w)
		for i := len(f.Server.Features) - 1; i >= 0; i-- {
			feat := f.Server.Features[i]
			lr, ok := feat.(LoginRedirector)
			if ok {
				nextURL = lr.LoginRedirect(nextURL, r)
			}
		}
		http.Redirect(w, r, nextURL, http.StatusTemporaryRedirect)
		return
	}

	nextURL = r.URL.Query().Get("then")
	if nextURL == "" {
		u, _ := url.Parse(r.Referer())
		nextURL = u.Path
	}
	session.Values["after_login"] = nextURL
	session.Values["after_callback"] = r.URL.Path
	session.Save(r, w)

	f.Nonce, _ = generateNonce()
	callbackUrl := fmt.Sprintf("http://%s/feature/login/callback", r.Host)

	ctx := zt100.FromContext(r.Context())

	if ctx.HasFeature("login:custom") {
		issuerParts, _ := url.Parse(os.Getenv("ISSUER"))
		baseUrl := issuerParts.Scheme + "://" + issuerParts.Hostname()
		f.Server.Template.ExecuteTemplate(w, "login.html", struct {
			Profile         map[string]string
			IsAuthenticated bool
			BaseUrl         string
			ClientId        string
			Issuer          string
			State           string
			Nonce           string
			RedirectUrl     string
		}{
			Profile:         f.getProfileData(r),
			IsAuthenticated: f.isAuthenticated(r),
			BaseUrl:         baseUrl,
			ClientId:        os.Getenv("CLIENT_ID"),
			Issuer:          os.Getenv("ISSUER"),
			State:           f.State,
			Nonce:           f.Nonce,
			RedirectUrl:     callbackUrl,
		})
		return
	}

	q := r.URL.Query()
	q.Add("client_id", os.Getenv("CLIENT_ID"))
	q.Add("response_type", "code")
	q.Add("response_mode", "query")
	q.Add("scope", "openid profile email")
	q.Add("redirect_uri", callbackUrl)
	q.Add("state", f.State)
	q.Add("nonce", f.Nonce)
	nextUrl := fmt.Sprintf("%s/v1/authorize?%s", os.Getenv("ISSUER"), q.Encode())

	http.Redirect(w, r, nextUrl, http.StatusTemporaryRedirect)
}

func (f *LoginFeature) logout(w http.ResponseWriter, r *http.Request) {
	session, err := f.Sessions.Get(r, "zt100-session-store")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	delete(session.Values, "id_token")
	delete(session.Values, "access_token")

	session.Save(r, w)

	then := r.URL.Query().Get("then")
	if then == "" {
		then = r.Referer()
	}

	http.Redirect(w, r, then, http.StatusTemporaryRedirect)
}

func (f *LoginFeature) callback(w http.ResponseWriter, r *http.Request) {
	// Check the state that was returned in the query string is the same as the above state
	if r.URL.Query().Get("state") != f.State {
		fmt.Fprintln(w, "The state was not as expected")
		return
	}
	// Make sure the code was provided
	if r.URL.Query().Get("code") == "" {
		fmt.Fprintln(w, "The code was not returned or is not accessible")
		return
	}

	exchange := exchangeCode(r.URL.Query().Get("code"), r)
	//log.Println(exchange)

	session, err := f.Sessions.Get(r, "zt100-session-store")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	_, verificationError := f.verifyToken(exchange.IdToken)

	if verificationError != nil {
		fmt.Println(verificationError)
	}

	if verificationError == nil {
		session.Values["id_token"] = exchange.IdToken
		session.Values["access_token"] = exchange.AccessToken

		session.Save(r, w)
	}

	nextURL, ok := session.Values["after_callback"].(string)
	if !ok || nextURL == "" {
		log.Println("No after_callback session value")
		nextURL = "/"
	}

	http.Redirect(w, r, nextURL, http.StatusTemporaryRedirect)
}

func (f *LoginFeature) isAuthenticated(r *http.Request) bool {
	session, err := f.Sessions.Get(r, "zt100-session-store")

	if err != nil || session.Values["id_token"] == nil || session.Values["id_token"] == "" {
		return false
	}

	return true
}

func (f *LoginFeature) getProfileData(r *http.Request) map[string]string {
	m := make(map[string]string)

	session, err := f.Sessions.Get(r, "zt100-session-store")

	if err != nil || session.Values["access_token"] == nil || session.Values["access_token"] == "" {
		return m
	}

	reqUrl := os.Getenv("ISSUER") + "/v1/userinfo"

	req, _ := http.NewRequest("GET", reqUrl, bytes.NewReader([]byte("")))
	h := req.Header
	h.Add("Authorization", "Bearer "+session.Values["access_token"].(string))
	h.Add("Accept", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	json.Unmarshal(body, &m)

	return m
}

func (f *LoginFeature) verifyToken(t string) (*verifier.Jwt, error) {
	tv := map[string]string{}
	tv["nonce"] = f.Nonce
	tv["aud"] = os.Getenv("CLIENT_ID")
	jv := verifier.JwtVerifier{
		Issuer:           os.Getenv("ISSUER"),
		ClaimsToValidate: tv,
	}

	result, err := jv.New().VerifyIdToken(t)

	if err != nil {
		return nil, fmt.Errorf("%s", err)
	}

	if result != nil {
		return result, nil
	}

	return nil, fmt.Errorf("token could not be verified: %s", "")
}

func generateNonce() (string, error) {
	nonceBytes := make([]byte, 32)
	_, err := rand.Read(nonceBytes)
	if err != nil {
		return "", fmt.Errorf("could not generate nonce")
	}

	return base64.URLEncoding.EncodeToString(nonceBytes), nil
}

func exchangeCode(code string, r *http.Request) Exchange {
	authHeader := base64.StdEncoding.EncodeToString(
		[]byte(os.Getenv("CLIENT_ID") + ":" + os.Getenv("CLIENT_SECRET")))

	q := r.URL.Query()
	q.Add("grant_type", "authorization_code")
	q.Add("code", code)
	q.Add("redirect_uri", fmt.Sprintf("http://%s/feature/login/callback", r.Host))

	url := os.Getenv("ISSUER") + "/v1/token?" + q.Encode()

	req, _ := http.NewRequest("POST", url, bytes.NewReader([]byte("")))
	h := req.Header
	h.Add("Authorization", "Basic "+authHeader)
	h.Add("Accept", "application/json")
	h.Add("Content-Type", "application/x-www-form-urlencoded")
	h.Add("Connection", "close")
	h.Add("Content-Length", "0")

	client := &http.Client{}
	resp, _ := client.Do(req)
	body, _ := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()
	var exchange Exchange
	json.Unmarshal(body, &exchange)

	return exchange

}

type Exchange struct {
	Error            string `json:"error,omitempty"`
	ErrorDescription string `json:"error_description,omitempty"`
	AccessToken      string `json:"access_token,omitempty"`
	TokenType        string `json:"token_type,omitempty"`
	ExpiresIn        int    `json:"expires_in,omitempty"`
	Scope            string `json:"scope,omitempty"`
	IdToken          string `json:"id_token,omitempty"`
}
