module tractor/zt100

go 1.14

replace worksite => ./

replace github.com/progrium/zt100 => ../

require (
	github.com/manifold/tractor v0.0.0
	github.com/progrium/zt100 v0.0.0-20201127154403-6f4f1a173508
	worksite v0.0.0-00010101000000-000000000000
)


replace github.com/manifold/tractor => ../tractor
replace github.com/manifold/qtalk => ../tractor/qtalk
replace github.com/progrium/macdriver => ../macdriver


