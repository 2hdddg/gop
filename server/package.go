package server

import (
	"fmt"
)

type pack struct {
	name  string
	funcs map[string]FileLocation
}

func (p *pack) merge(f *file) {
	// Simple but expensive merge by removing all locations
	// for the file and adding all symbols in the file again.
	for name, loc := range p.funcs {
		if loc.FilePath == f.path {
			delete(p.funcs, name)
		}
	}
	for name, loc := range f.funcs {
		p.funcs[name] = FileLocation{Location: loc, FilePath: f.path}
	}
}

func (p *pack) find(name string) *FileLocation {
	fmt.Printf("Looking for %v\n", name)
	l, ok := p.funcs[name]
	if ok {
		return &l
	}
	return nil
}

func newPackage(path, name string) *pack {
	funcs := make(map[string]FileLocation)
	return &pack{name: name, funcs: funcs}
}
