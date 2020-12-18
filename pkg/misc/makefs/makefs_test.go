package makefs

import (
	"bytes"
	"testing"

	"github.com/progrium/zt100/pkg/misc/esbuild"
	"github.com/spf13/afero"
)

func TestMakefs(t *testing.T) {
	existFile := []byte("foo")
	srcFile := []byte("<html></html>\n")
	dstFile := []byte("m(\"html\", null);\n")
	f := afero.NewMemMapFs()
	if err := afero.WriteFile(f, "exists.js", existFile, 0644); err != nil {
		t.Fatal(err)
	}
	if err := afero.WriteFile(f, "html.jsx", srcFile, 0644); err != nil {
		t.Fatal(err)
	}
	mfs := New(f, f)
	mfs.Register(".js", ".jsx", func(fs afero.Fs, dst, src string) ([]byte, error) {
		return esbuild.BuildFile(fs, src)
	})

	t.Run("made file", func(t *testing.T) {
		got, err := afero.ReadFile(mfs, "html.js")
		if err != nil {
			t.Fatal(err)
		}
		if !bytes.Equal(got, dstFile) {
			t.Errorf("got %q, want %q", got, dstFile)
		}
	})

	t.Run("existing file", func(t *testing.T) {
		got, err := afero.ReadFile(mfs, "exists.js")
		if err != nil {
			t.Fatal(err)
		}
		if !bytes.Equal(got, existFile) {
			t.Errorf("got %q, want %q", got, existFile)
		}
	})
}
