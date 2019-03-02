package server

import (
	"log"
)

// Keeps track of exported functions in package.
// Each function has a file and location in the file.
type pack struct {
	name  string
	funcs map[string]FileLocation
}

// Merges functions in parsed file to the list of files in
// the package. Handles case where file has been merged before.
func (p *pack) mergeFile(f *file) {
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

func (p *pack) findFunc(name string) *FileLocation {
	log.Printf("Looking for %v\n", name)
	l, ok := p.funcs[name]
	if ok {
		return &l
	}
	return nil
}

func newPackage(name string) *pack {
	funcs := make(map[string]FileLocation)
	return &pack{name: name, funcs: funcs}
}
