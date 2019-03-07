package server

import ()

// Keeps track of exported functions in package.
// Each function has a file and location in the file.
type pack struct {
	name  string              // Qualified name: go/ast
	funcs map[string]Location // Functions names lookup
}

// Merges functions in parsed file to the list of files in
// the package. Handles case where file has been merged before.
func (p *pack) mergeFile(f *file) {
	// Simple but expensive merge by removing all locations
	// for the file and adding all symbols in the file again.
	for funcName, location := range p.funcs {
		if location.Path == f.path {
			delete(p.funcs, funcName)
		}
	}
	for funcName, line := range f.funcs {
		p.funcs[funcName] = Location{
			Path: f.path,
			Line: int(line)}
	}
}

// Find function in package
func (p *pack) findFunc(name string) *Location {
	l, ok := p.funcs[name]
	if ok {
		return &l
	}
	return nil
}

func newPack(name string) *pack {
	return &pack{
		name:  name,
		funcs: make(map[string]Location),
	}
}
