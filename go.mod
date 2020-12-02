module github.com/progrium/zt100

go 1.14

replace worksite => ./.tractor

replace github.com/manifold/tractor => /Users/progrium/Source/github.com/manifold/tractor

replace github.com/manifold/qtalk => /Users/progrium/Source/github.com/manifold/tractor/qtalk

replace github.com/progrium/macdriver => /Users/progrium/Source/github.com/progrium/macdriver

require (
	github.com/EdlinOrg/prominentcolor v1.0.0 // indirect
	github.com/gorilla/mux v1.8.0
	github.com/lucasb-eyer/go-colorful v1.0.3 // indirect
	github.com/manifold/tractor v0.0.0-00010101000000-000000000000
	github.com/nfnt/resize v0.0.0-20180221191011-83c6a9932646 // indirect
	github.com/oliamb/cutter v0.2.2 // indirect
	github.com/progrium/watcher v1.0.8-0.20200403214642-88c0f931de38
	github.com/spf13/afero v1.4.1
)
