package server

import (
	"log"
)

// Keeps track of exported functions in package.
// Each function has a file and location in the file.
type pack struct {
	name  string
	path  string
	funcs map[string]Location
}

// Merges functions in parsed file to the list of files in
// the package. Handles case where file has been merged before.
func (p *pack) mergeFile(f *file) {
	// Simple but expensive merge by removing all locations
	// for the file and adding all symbols in the file again.
	for name, loc := range p.funcs {
		if loc.Path == f.path {
			delete(p.funcs, name)
		}
	}
	for name, line := range f.funcs {
		p.funcs[name] = Location{Path: f.path, Line: int(line)}
	}
}

func (p *pack) findFunc(name string) *Location {
	log.Printf("Looking for %v\n", name)
	l, ok := p.funcs[name]
	if ok {
		return &l
	}
	return nil
}

func newPackage(name, path string) *pack {
	funcs := make(map[string]Location)
	return &pack{name: name, path: path, funcs: funcs}
}
