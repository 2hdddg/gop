package tree

import (
	"path"

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

func NewTree(path string) *Tree {
	return &Tree{
		Path:  path,
		Packs: make([]*Package, 3),
	}
}

func (t *Tree) AddPackage(name string) *Package {
	s := newPackage(name, t.Path)
	t.Packs = append(t.Packs, s)
	return s
}

func (p *Package) AddPackage(name string) *Package {
	s := newPackage(name, p.Path)
	p.Packs = append(p.Packs, s)
	return s
}

func (p *Package) AddFile(name string, parse Parse) (*File, error) {
	filepath := path.Join(p.Path, name)
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
		Files: make([]*File, 3),
		Packs: make([]*Package, 3),
		Path:  path.Join(parent, name),
	}
}
