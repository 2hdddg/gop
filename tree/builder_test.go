package tree

import (
	"testing"
)

func TestNewBuilder(t *testing.T) {
	_, err := NewBuilder(".")

	assertNoError(t, err)
}

func TestBuild(t *testing.T) {
	b, _ := NewBuilder("/")
	filesys := &FakeDir{
		name: "",
		dirs: []*FakeDir{
			&FakeDir{
				name: "pack1",
			},
			&FakeDir{
				name: "pack2",
			},
		},
		files: []*FakeFile{},
	}

	b.reader = filesys
	b.parser = &FakeParser{}
	tree, err := b.Build()

	assertNoError(t, err)
	if len(tree.Packs) != 2 {
		t.Error("Packs != 2")
	}
}
