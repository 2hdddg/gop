package main

import (
	"fmt"
)

type FileLocation struct {
	Location
	FilePath string
}

type Package struct {
	Name  string
	Funcs map[string]FileLocation
}

func (p *Package) Merge(f *File) {
	// Simple but expensive merge by removing all locations
	// for the file and adding all symbols in the file again.
	for name, loc := range p.Funcs {
		if loc.FilePath == f.Path {
			delete(p.Funcs, name)
		}
	}
	for name, loc := range f.Funcs {
		p.Funcs[name] = FileLocation{Location: loc, FilePath: f.Path}
	}
}

func (p *Package) Find(name string) *FileLocation {
	fmt.Printf("Looking for %v\n", name)
	l, ok := p.Funcs[name]
	if ok {
		return &l
	}
	return nil
}

func NewPackage(path, name string) *Package {
	funcs := make(map[string]FileLocation)
	return &Package{Name: name, Funcs: funcs}
}
