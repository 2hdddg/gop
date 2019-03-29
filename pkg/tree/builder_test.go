package tree

import (
	"testing"
)

func TestNewBuilder(t *testing.T) {
	_, err := NewBuilder(".")

	assertNoError(t, err)
}

func buildFlatFakeFilesystem() *FakeDir {
	return &FakeDir{
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
}

func TestBuild(t *testing.T) {
	b, _ := NewBuilder("/")
	filesys := buildFlatFakeFilesystem()

	b.reader = filesys
	b.parser = &FakeParser{}
	tree, err := b.Build()

	assertNoError(t, err)
	if len(tree.Packs) != 2 {
		t.Error("Packs != 2")
	}
}

func TestBuildWithProgress(t *testing.T) {
	b, _ := NewBuilder("/")
	filesys := buildFlatFakeFilesystem()
	r := &ProgressRecorder{}
	b.Progress = r

	b.reader = filesys
	b.parser = &FakeParser{}
	_, _ = b.Build()

	if r.NumPackageParsed != 2 {
		t.Error("Should have recorded 2 package callbacks")
	}
	if r.NumTreesParsed != 1 {
		t.Error("Should have recorded 1 tree callback")
	}
}
