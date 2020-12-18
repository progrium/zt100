package object

import (
	"github.com/progrium/zt100/pkg/manifold"
	"github.com/progrium/zt100/pkg/misc/notify"
)

// observe component list changes

func (o *object) AppendComponent(com manifold.Component) {
	o.componentlist.AppendComponent(com)
	com.SetContainer(o)
	o.UpdateRegistry()
	o.registry.Populate(com.Pointer())
	notify.Send(o, manifold.ObjectChange{
		Object: o,
		Path:   "::Components",
		New:    com,
	})

}

func (o *object) RemoveComponent(com manifold.Component) {
	o.componentlist.RemoveComponent(com)
	o.UpdateRegistry()
	notify.Send(o, manifold.ObjectChange{
		Object: o,
		Path:   "::Components",
		Old:    com,
	})
	if o.main == com {
		o.main = nil
	}
}

func (o *object) MoveComponentAt(idx, newidx int) {
	o.componentlist.MoveComponentAt(idx, newidx)
	notify.Send(o, manifold.ObjectChange{
		Object: o,
		Path:   "::Components",
	})
}

func (o *object) InsertComponentAt(idx int, com manifold.Component) {
	o.componentlist.InsertComponentAt(idx, com)
	com.SetContainer(o)
	o.UpdateRegistry()
	o.registry.Populate(com.Pointer())
	notify.Send(o, manifold.ObjectChange{
		Object: o,
		Path:   "::Components",
		New:    com,
	})
}

func (o *object) RemoveComponentAt(idx int) manifold.Component {
	c := o.componentlist.RemoveComponentAt(idx)
	o.UpdateRegistry()
	notify.Send(o, manifold.ObjectChange{
		Object: o,
		Path:   "::Components",
		Old:    c,
	})
	if o.main == c {
		o.main = nil
	}
	return c
}
