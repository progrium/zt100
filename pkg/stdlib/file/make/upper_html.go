package make

import (
	"path/filepath"
	"strings"

	"github.com/spf13/afero"
)

type UpperHTML struct{}

func (b *UpperHTML) Match(name string) (string, bool) {
	if ok, _ := filepath.Match("*_upper.html", filepath.Base(name)); ok {
		return strings.Replace(name, "_upper.html", ".html", 1), true
	}
	return "", false
}

func (b *UpperHTML) Build(fs afero.Fs, dst, src string) error {
	srcData, err := afero.ReadFile(fs, src)
	if err != nil {
		return err
	}
	built := []byte(strings.ToUpper(string(srcData)))
	return afero.WriteFile(fs, dst, built, 0644)
}
