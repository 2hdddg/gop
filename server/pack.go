package server

import (
	"log"
)

// Keeps track of exported functions in package.
// Each function has a file and location in the file.
type pack struct {
	packPath string              // Full path to package
	packName string              // Qualified name: go/ast
	funcs    map[string]Location // Functions names lookup
}

// Merges functions in parsed file to the list of files in
// the package. Handles case where file has been merged before.
func (p *pack) mergeFile(f *file) {
	// Simple but expensive merge by removing all locations
	// for the file and adding all symbols in the file again.
	for funcName, location := range p.funcs {
		if location.Path == f.filePath {
			delete(p.funcs, funcName)
		}
	}
	for funcName, line := range f.funcs {
		p.funcs[funcName] = Location{
			Path: f.filePath,
			Line: int(line)}
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
	return &pack{packName: name, packPath: path, funcs: funcs}
}
