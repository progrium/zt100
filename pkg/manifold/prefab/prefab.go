package prefab

import (
	"github.com/progrium/zt100/pkg/manifold"
	"github.com/progrium/zt100/pkg/manifold/library"
	"github.com/progrium/zt100/pkg/manifold/object"
	"github.com/rs/xid"
)

var (
	registered []*RegisteredPrefab
)

type RegisteredPrefab struct {
	Snapshot manifold.ObjectPrefab
	Name     string
	ID       string
}

func (rp *RegisteredPrefab) New() manifold.Object {
	return loadPrefab(rp.Snapshot)
}

func loadPrefab(p manifold.ObjectPrefab) manifold.Object {
	snapshot := manifold.ObjectSnapshot{
		ID:         xid.New().String(),
		Name:       p.Name,
		Attrs:      p.Attrs,
		Main:       p.Main,
		Components: p.Components,
	}
	obj := object.FromSnapshot(snapshot)
	library.LoadComponents(obj, snapshot)
	for _, cp := range p.Children {
		child := loadPrefab(cp)
		obj.AppendChild(child)
	}
	return obj
}

func Register(prefabs []manifold.ObjectPrefab) {
	for _, p := range prefabs {
		registered = append(registered, &RegisteredPrefab{
			Snapshot: p,
			Name:     p.Name,
			ID:       p.ID,
		})
	}
}

func Registered() []*RegisteredPrefab {
	r := make([]*RegisteredPrefab, len(registered))
	copy(r, registered)
	return r
}

func LookupID(id string) *RegisteredPrefab {
	for _, rc := range registered {
		if rc.ID == id {
			return rc
		}
	}
	return nil
}
