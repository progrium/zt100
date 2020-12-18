package menu

import "context"

func (f *Framework) InitializeDaemon() error {
	f.Registry = &Registry{
		menus: make(map[string]func(context.Context) []Item),
	}

	for _, c := range f.Contributors {
		c.ContributeMenus(f.Registry)
	}

	// TODO: better place for this?
	f.Registry.Register("context", func(ctx context.Context) []Item {
		return []Item{
			{Cmd: "copy", Label: "Copy", Icon: "copy"},
			{Cmd: "cut", Label: "Cut", Icon: "cut"},
			{Cmd: "paste", Label: "Paste", Icon: "paste"},
		}
	})

	return nil
}
