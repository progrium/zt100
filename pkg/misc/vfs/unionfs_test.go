package vfs

import (
	"testing"

	"github.com/spf13/afero"
)

func TestUnionFS(t *testing.T) {

	layer0 := &afero.MemMapFs{}
	layer1 := &afero.MemMapFs{}
	layer2 := &afero.MemMapFs{}

	if err := afero.WriteFile(layer0, "all.txt", []byte("layer0"), 0755); err != nil {
		t.Fatal(err)
	}
	if err := afero.WriteFile(layer1, "all.txt", []byte("layer1"), 0755); err != nil {
		t.Fatal(err)
	}
	if err := afero.WriteFile(layer2, "all.txt", []byte("layer2"), 0755); err != nil {
		t.Fatal(err)
	}
	if err := afero.WriteFile(layer0, "layers/layer0.txt", []byte("layer0"), 0755); err != nil {
		t.Fatal(err)
	}
	if err := afero.WriteFile(layer1, "layers/layer1.txt", []byte("layer1"), 0755); err != nil {
		t.Fatal(err)
	}
	if err := afero.WriteFile(layer2, "layers/layer2.txt", []byte("layer2"), 0755); err != nil {
		t.Fatal(err)
	}

	unionfs := NewUnionFS(layer0, layer1, layer2)

	t.Run("open-all.txt", func(t *testing.T) {
		b, err := afero.ReadFile(unionfs, "all.txt")
		if err != nil {
			t.Fatal(err)
		}
		if string(b) != "layer2" {
			t.Fatalf("unexpected contents for file; got: %s", string(b))
		}
	})

	t.Run("open-layer0.txt", func(t *testing.T) {
		b, err := afero.ReadFile(unionfs, "layers/layer0.txt")
		if err != nil {
			t.Fatal(err)
		}
		if string(b) != "layer0" {
			t.Fatalf("unexpected contents for file; got: %s", string(b))
		}
	})

	t.Run("open-layer1.txt", func(t *testing.T) {
		b, err := afero.ReadFile(unionfs, "layers/layer1.txt")
		if err != nil {
			t.Fatal(err)
		}
		if string(b) != "layer1" {
			t.Fatalf("unexpected contents for file; got: %s", string(b))
		}
	})

	t.Run("open-layer2.txt", func(t *testing.T) {
		b, err := afero.ReadFile(unionfs, "layers/layer2.txt")
		if err != nil {
			t.Fatal(err)
		}
		if string(b) != "layer2" {
			t.Fatalf("unexpected contents for file; got: %s", string(b))
		}
	})

	t.Run("stat-layer1.txt", func(t *testing.T) {
		_, err := unionfs.Stat("layers/layer1.txt")
		if err != nil {
			t.Fatal(err)
		}
	})

	t.Run("dir-layers", func(t *testing.T) {
		fi, err := afero.ReadDir(unionfs, "layers")
		if err != nil {
			t.Fatal(err)
		}
		if len(fi) != 3 {
			t.Fatalf("unexpected file count for dir; got: %v", fi)
		}
	})

}
