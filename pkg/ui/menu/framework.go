package menu

import (
	L "github.com/progrium/zt100/pkg/misc/logging"
)

type Contributor interface {
	ContributeMenus(menus *Registry)
}

type Framework struct {
	*Registry

	Contributors []Contributor
	Log          L.Logger
}
