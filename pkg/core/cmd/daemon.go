package cmd

import (
	"log"
)

func (f *Framework) InitializeDaemon() error {
	f.cmds = &Registry{
		cmds: make(map[string]Definition),
	}

	for _, c := range f.Contributors {
		c.ContributeCommands(f.cmds)
	}

	f.cmds.Register(Definition{
		ID:       "debug.print",
		Label:    "Print",
		Category: "Debug",
		Desc:     "Print params",
		Run: func(params interface{}) {
			log.Println(params)
		},
	})

	return nil
}
