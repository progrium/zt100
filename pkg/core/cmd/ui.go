package cmd

type ViewState struct {
	Commands []string
}

func (f *Framework) InitializeState() (name string, v interface{}) {
	f.view = &ViewState{}

	name = "cmd"
	v = f.view

	return
}

func (f *Framework) UpdateState() error {
	return nil
}
