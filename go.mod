module github.com/progrium/zt100

go 1.14

replace worksite => ./.tractor

require (
	github.com/EdlinOrg/prominentcolor v1.0.0 // indirect
	github.com/gorilla/mux v1.8.0
	github.com/gorilla/sessions v1.2.1 // indirect
	github.com/lucasb-eyer/go-colorful v1.0.3 // indirect
	github.com/nfnt/resize v0.0.0-20180221191011-83c6a9932646 // indirect
	github.com/okta/okta-jwt-verifier-golang v1.0.0 // indirect
	github.com/oliamb/cutter v0.2.2 // indirect
	github.com/progrium/watcher v1.0.8-0.20200403214642-88c0f931de38
	github.com/spf13/afero v1.4.1
)

// These are only necessary when also developing/changing tractor
// (and/or qtalk+macdriver). They are left uncommented as this is
// the assumed default for users of this repository at this time.
// The directories they point to are not in this repository and
// can be created through manual checkouts or symlinks.
// They must also be kept in sync with the equivalent replace
// directives in .tractor/go.mod so if these are commented out
// they should also be commented out there.
replace github.com/manifold/tractor => ./tractor

replace github.com/manifold/qtalk => ./qtalk

replace github.com/progrium/macdriver => ./macdriver
