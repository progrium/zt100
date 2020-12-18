module github.com/progrium/zt100

go 1.14

replace worksite => ./.tractor

require (
	github.com/EdlinOrg/prominentcolor v1.0.0
	github.com/dave/jennifer v1.4.0
	github.com/davecgh/go-spew v1.1.1
	github.com/dustin/gojson v0.0.0-20160307161227-2e71ec9dd5ad
	github.com/goji/httpauth v0.0.0-20160601135302-2da839ab0f4d
	github.com/gorilla/mux v1.8.0
	github.com/gorilla/sessions v1.2.1
	github.com/gorilla/websocket v1.4.2
	github.com/keybase/go-ps v0.0.0-20190827175125-91aafc93ba19
	github.com/lucasb-eyer/go-colorful v1.0.3 // indirect
	github.com/manifold/qtalk v0.0.0-00010101000000-000000000000
	github.com/miekg/dns v1.1.28
	github.com/mitchellh/mapstructure v1.4.0
	github.com/nfnt/resize v0.0.0-20180221191011-83c6a9932646 // indirect
	github.com/nickvanw/ircx/v2 v2.0.0
	github.com/okta/okta-jwt-verifier-golang v1.0.0
	github.com/oliamb/cutter v0.2.2 // indirect
	github.com/progrium/esbuild v0.0.0-20200327212623-fae14fb26173
	github.com/progrium/prototypes v0.0.0-20190807232325-d9b2b4ba3a4f
	github.com/progrium/watcher v1.0.8-0.20200403214642-88c0f931de38
	github.com/rs/xid v1.2.1
	github.com/skratchdot/open-golang v0.0.0-20200116055534-eef842397966
	github.com/spf13/afero v1.4.1
	github.com/stretchr/testify v1.5.1
	github.com/urfave/negroni v1.0.0
	go.uber.org/zap v1.13.0
	golang.org/x/net v0.0.0-20200707034311-ab3426394381
	gopkg.in/sorcix/irc.v2 v2.0.0-20190306112350-8d7a73540b90
)

replace github.com/manifold/qtalk => ./qtalk
