package server

import (
	"path"
)

// Set of packs used to build an index.
type packs map[string]*pack

func newPacks() packs {
	return packs(make(map[string]*pack))
}

// Ensures that package referenced by parsed file exists in
// the set of packs.
func (packs *packs) ensurePackage(f *file) *pack {
	packPath := path.Dir(f.filePath)
	key := path.Join(packPath, f.packName)

	p, exists := (*packs)[key]
	if !exists {
		p = newPackage(f.packName, packPath)
		(*packs)[key] = p
	}

	return p
}

// Adds a pack to the set indirectly by adding a file that belongs to
// a package.
func (s *packs) addFile(f *file) {
	p := s.ensurePackage(f)
	p.mergeFile(f)
}

func (s *packs) buildIndex() *index {
	packs := make(map[string]*pack)
	funcs := make(map[string][]*pack)

	for _, p := range *s {
		// Important to copy package since it will be sent
		// to channel handling search while it might be modified
		// by channel that indexes packages.
		pcopy := *p
		packs[p.packName] = &pcopy
		for n, _ := range pcopy.funcs {
			funcs[n] = append(funcs[n], &pcopy)
		}
	}
	return &index{packs: packs, funcs: funcs}
}
