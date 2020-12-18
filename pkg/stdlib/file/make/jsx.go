package make

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/progrium/zt100/pkg/misc/esbuild"
	"github.com/spf13/afero"
)

type JSX struct {
	DstExt  string
	SrcExt  string
	Factory string
}

func (b *JSX) Initialize() {
	if b.DstExt == "" {
		b.DstExt = ".mjs"
	}
	if b.SrcExt == "" {
		b.SrcExt = ".jsx"
	}
	if b.Factory == "" {
		b.Factory = "h"
	}
}

func (b *JSX) Match(name string) (string, bool) {
	if ok, _ := filepath.Match(fmt.Sprintf("*%s", b.DstExt), filepath.Base(name)); ok {
		return strings.Replace(name, b.DstExt, b.SrcExt, 1), true
	}
	return "", false
}

func (b *JSX) Build(fs afero.Fs, dst, src string) error {
	esbuild.JsxFactory = b.Factory
	built, err := esbuild.BuildFile(fs, src)
	if err != nil {
		return err
	}
	return afero.WriteFile(fs, dst, built, 0644)
}
