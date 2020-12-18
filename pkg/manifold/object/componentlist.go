package object

import (
	"strconv"
	"strings"

	"github.com/progrium/zt100/pkg/manifold"
)

type componentlist struct {
	components []manifold.Component
}

func (l *componentlist) Components() []manifold.Component {
	c := make([]manifold.Component, len(l.components))
	copy(c, l.components)
	return c
}

func (l *componentlist) AppendComponent(com manifold.Component) {
	l.components = append(l.components, com)
}

func (l *componentlist) RemoveComponent(com manifold.Component) {
	defer com.SetContainer(nil)
	for idx, c := range l.components {
		if c == com {
			l.RemoveComponentAt(idx)
			return
		}
	}
}

func (l *componentlist) MoveComponentAt(idx, newidx int) {
	if idx < 0 || idx > len(l.components) || newidx < 0 || newidx > len(l.components) {
		panic("index out of range")
	}
	c := l.components[idx]
	l.components = append(l.components[:idx], l.components[idx+1:]...)
	if newidx == len(l.components) {
		l.AppendComponent(c)
	} else {
		l.InsertComponentAt(newidx, c)
	}
}

func (l *componentlist) InsertComponentAt(idx int, com manifold.Component) {
	l.components = append(l.components[:idx], append([]manifold.Component{com}, l.components[idx:]...)...)
}

func (l *componentlist) RemoveComponentAt(idx int) manifold.Component {
	c := l.components[idx]
	l.components = append(l.components[:idx], l.components[idx+1:]...)
	c.SetContainer(nil)
	return c
}

func (l *componentlist) HasComponent(com manifold.Component) bool {
	for _, c := range l.components {
		if c == com {
			return true
		}
	}
	return false
}

func (l *componentlist) Component(name string) manifold.Component {
	// support taking a full data path for convenience?
	for _, c := range l.components {
		if c.Name() == name || c.ID() == name {
			return c
		}
	}
	path := strings.Split(name, "/")
	parts := strings.Split(path[0], ":")
	idx, err := strconv.Atoi(parts[0])
	if err != nil {
		return nil
	}
	return l.components[idx]

}
