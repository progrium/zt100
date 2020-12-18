package ui

type Script struct {
	Src string `json:"$cmd$script,omitempty"`
}

func JS(src string) Script {
	return Script{Src: src}
}
