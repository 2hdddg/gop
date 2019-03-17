package tree

import (
	"testing"
)

func TestNewBuilder(t *testing.T) {
	_, err := NewBuilder(onParsedFake, ".")

	assertNoError(t, err)
}

func TestBuild(t *testing.T) {
	b, _ := NewBuilder(onParsedFake, "/")
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

	b.readDir = MakeReadDir(filesys)
	b.parse = parseFake
	tree, err := b.Build()

	assertNoError(t, err)
	if len(tree.Packs) != 2 {
		t.Error("Packs != 2")
	}
}
