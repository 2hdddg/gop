package tree

import (
	"path/filepath"
	"testing"
)

func TestNew(t *testing.T) {
	tree, err := NewTree("/")

	// Should be a an empty list of packages
	assertEmptyListOfPackages(t, tree.Packs)

	// Should make relative paths absolute
	tree, err = NewTree(".")
	abspath, _ := filepath.Abs(".")
	assertPath(t, abspath, tree.Path)

	// Should fail on non-existent path
	tree, err = NewTree("/hubba/bubba/xyz")
	if err == nil {
		t.Errorf("Should return error if path doesn't exist")
	}
}

func setup() (tree *Tree, root string) {
	tree, _ = NewTree(".")
	root, _ = filepath.Abs(".")
	return
}

func TestAddPackageToTree(t *testing.T) {
	tree, root := setup()

	pack := tree.AddPackage("pname")

	assertName(t, "pname", pack.Name)
	assertPath(t, filepath.Join(root, pack.Name), pack.Path)
	assertEmptyListOfPackages(t, pack.Packs)
	assertEmptyListOfFiles(t, pack.Files)
}

func TestAddPackageToPackage(t *testing.T) {
	tree, root := setup()

	parent := tree.AddPackage("parent")
	pack := parent.AddPackage("pname")

	assertName(t, "parent/pname", pack.Name)
	assertPath(t, filepath.Join(root, pack.Name), pack.Path)
	assertEmptyListOfPackages(t, pack.Packs)
	assertEmptyListOfFiles(t, pack.Files)
}

func TestAddFileToPackage(t *testing.T) {
	tree, root := setup()

	pack := tree.AddPackage("pname")
	file, err := pack.AddFile("fname", parseFake)

	assertName(t, "fname", file.Name)
	assertPath(t, filepath.Join(root, "pname", "fname"), file.Path)
	assertNoError(t, err)
}
