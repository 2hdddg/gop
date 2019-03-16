package tree

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/2hdddg/gop/parser"
)

type File struct {
	Name string
	Path string
	Syms *parser.Symbols
}

type Parse func(path string) (*parser.Symbols, error)

type Package struct {
	Name  string
	Path  string
	Files []*File
	Packs []*Package
}

type Tree struct {
	Path  string
	Packs []*Package
}

func NewTree(path string) (*Tree, error) {
	// Check that root is valid
	if _, err := os.Stat(path); os.IsNotExist(err) {
		err = fmt.Errorf("Tree root doesn't exist: %v", path)
		return nil, err
	}

	// Make the root absolute. Use absolute paths everywhere to
	// make comparisons simpler when matching client roots with
	// server roots.
	abspath, err := filepath.Abs(path)
	if err != nil {
		return nil, err
	}

	return &Tree{
		Path:  abspath,
		Packs: []*Package{},
	}, nil
}

func (t *Tree) AddPackage(name string) *Package {
	s := newPackage(name, t.Path)
	t.Packs = append(t.Packs, s)
	return s
}

func (p *Package) AddPackage(name string) *Package {
	s := newPackage(name, p.Path)
	s.Name = strings.Join([]string{p.Name, name}, "/")
	p.Packs = append(p.Packs, s)
	return s
}

func (p *Package) AddFile(name string, parse Parse) (*File, error) {
	filepath := filepath.Join(p.Path, name)
	syms, err := parse(filepath)
	if err != nil {
		return nil, err
	}
	file := &File{
		Name: name,
		Path: filepath,
		Syms: syms,
	}
	p.Files = append(p.Files, file)

	return file, nil
}

func newPackage(name, parent string) *Package {
	return &Package{
		Name:  name,
		Files: []*File{},
		Packs: []*Package{},
		Path:  filepath.Join(parent, name),
	}
}
