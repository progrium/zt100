package http

import (
	"github.com/progrium/zt100/pkg/manifold"
	"github.com/progrium/zt100/pkg/manifold/prefab"
)

func init() {
	prefab.Register([]manifold.ObjectPrefab{
		{
			ID:   "prefab123",
			Name: "My first prefab",
			Components: []manifold.ComponentSnapshot{
				{
					Name:    "http.Server",
					Enabled: true,
					Value:   nil,
				},
				{
					Name:    "http.SingleUserBasicAuth",
					Enabled: true,
					Value: map[string]interface{}{
						"Username": "Progrium",
						"Password": "foobar",
					},
				},
			},
		},

		{
			ID:   "prefab211",
			Name: "ChildTest",
			Components: []manifold.ComponentSnapshot{
				{
					Name:    "http.Server",
					Enabled: true,
					Value:   nil,
				},
			},
			Children: []manifold.ObjectPrefab{
				{
					Name: "child1",
				},
				{
					Name: "child2",
					Components: []manifold.ComponentSnapshot{
						{
							Name:    "http.Server",
							Enabled: true,
							Value:   nil,
						},
					},
				},
			},
		},
	})
}
