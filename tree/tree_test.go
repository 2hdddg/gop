package tree

import (
	"path/filepath"
	"testing"

	"github.com/2hdddg/gop/parser"
)

var (
	tree *Tree
	err  error
	root string
	pack *Package
	file *File
)

func TestNew(t *testing.T) {
	tree, err = NewTree("/")

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

func setup() {
	tree, _ = NewTree(".")
	root, _ = filepath.Abs(".")
}

func assertPath(t *testing.T, expected, actual string) {
	if expected != actual {
		t.Errorf("Expected path %v but was %v", expected, actual)
	}
}

func assertName(t *testing.T, expected, actual string) {
	if expected != actual {
		t.Errorf("Expected name %v but was %v", expected, actual)
	}
}

func assertEmptyListOfPackages(t *testing.T, packs []*Package) {
	if packs == nil || len(packs) != 0 {
		t.Errorf("Expected non nil empty list of packages")
	}
}

func assertEmptyListOfFiles(t *testing.T, files []*File) {
	if files == nil || len(files) != 0 {
		t.Errorf("Expected non nil empty list of files")
	}
}

func TestAddPackageToTree(t *testing.T) {
	setup()

	pack = tree.AddPackage("pname")

	assertName(t, "pname", pack.Name)
	assertPath(t, filepath.Join(root, pack.Name), pack.Path)
	assertEmptyListOfPackages(t, pack.Packs)
	assertEmptyListOfFiles(t, pack.Files)
}

func TestAddPackageToPackage(t *testing.T) {
	setup()

	parent := tree.AddPackage("parent")
	pack = parent.AddPackage("pname")

	assertName(t, "parent/pname", pack.Name)
	assertPath(t, filepath.Join(root, pack.Name), pack.Path)
	assertEmptyListOfPackages(t, pack.Packs)
	assertEmptyListOfFiles(t, pack.Files)
}

func parseFake(path string) (*parser.Symbols, error) {
	return parser.NewSymbols(), nil
}

func TestAddFileToPackage(t *testing.T) {
	setup()

	pack = tree.AddPackage("pname")
	file, err = pack.AddFile("fname", parseFake)

	assertName(t, "fname", file.Name)
	assertPath(t, filepath.Join(root, "pname", "fname"), file.Path)
}
