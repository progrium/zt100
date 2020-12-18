package object

import (
	"testing"

	"github.com/progrium/zt100/pkg/manifold"
	"github.com/progrium/zt100/pkg/misc/notify"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func testObj(name string) manifold.Object {
	obj := New(name)
	obj.(*object).notifyDebounce = nil
	return obj
}

func TestObserveObject(t *testing.T) {
	t.Run("::Parent", func(t *testing.T) {
		root := testObj("::root")
		parent := testObj("parent")
		root.AppendChild(parent)
		obj := testObj("test")
		root.AppendChild(obj)
		assertObserve(t, obj, obj, "::Parent", root, parent, func() {
			obj.SetParent(parent)
		})
	})
	t.Run("::SiblingIndex", func(t *testing.T) {
		root := testObj("root")
		c1 := testObj("c1")
		c2 := testObj("c2")
		root.AppendChild(c1)
		root.AppendChild(c2)
		assertObserve(t, c2, c2, "::SiblingIndex", 1, 0, func() {
			c2.SetSiblingIndex(0)
		})
	})
	t.Run("::Children/RemoveChildAt", func(t *testing.T) {
		obj := testObj("test")
		parent := testObj("parent")
		parent.AppendChild(obj)
		assertObserve(t, parent, parent, "::Children", obj, nil, func() {
			parent.RemoveChildAt(0)
		})
	})
	t.Run("::Children/InsertChildAt", func(t *testing.T) {
		obj := testObj("test")
		parent := testObj("parent")
		assertObserve(t, parent, parent, "::Children", nil, obj, func() {
			parent.InsertChildAt(0, obj)
		})
	})
	t.Run("::Children/RemoveChild", func(t *testing.T) {
		obj := testObj("test")
		parent := testObj("parent")
		parent.AppendChild(obj)
		assertObserve(t, parent, parent, "::Children", obj, nil, func() {
			parent.RemoveChild(obj)
		})
	})
	t.Run("::Children/AppendChild", func(t *testing.T) {
		obj := testObj("test")
		parent := testObj("parent")
		assertObserve(t, parent, parent, "::Children", nil, obj, func() {
			parent.AppendChild(obj)
		})
	})
}

func assertObserve(t *testing.T, subject manifold.Object, changed manifold.Object, subpath string, old, new interface{}, action func()) {
	notified := false
	t.Logf("assert oberserver %q %+v", t.Name(), changed.Subpath(subpath))
	notify.Observe(subject, notify.Func(func(e interface{}) error {
		var event manifold.ObjectChange
		var ok bool
		if event, ok = e.(manifold.ObjectChange); !ok {
			return nil
		}
		if event.Path != subpath {
			return nil
		}
		t.Logf("OnChange %q %+v %+v %+v", event.Path, event.Object, event.Old, event.New)
		assert.Equal(t, subject, event.Object)
		assert.Equal(t, subpath, event.Path)
		assert.Equal(t, old, event.Old)
		assert.Equal(t, new, event.New)
		notified = true
		return nil
	}))
	require.False(t, notified)
	action()
	require.True(t, notified)
}
