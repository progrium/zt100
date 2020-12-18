package ui

import "github.com/progrium/zt100/pkg/ui/state"

var State *state.Framework

func Update() error {
	if State != nil {
		return State.UpdateState()
	}
	return nil
}
