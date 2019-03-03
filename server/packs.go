package server

/*
type builder struct {
	packs map[string]*pack
}
*/

// Set of packs used to build an index.

type packs map[string]*pack

func newPacks() packs {
	return packs(make(map[string]*pack))
}

// Ensures that package referenced by parsed file exists in
// the set of packs.
func (s *packs) ensurePackage(f *file) *pack {
	path := f.packPath()
	name := f.packName
	key := path + name

	p, exists := (*s)[key]
	if !exists {
		p = newPackage(name, path)
		(*s)[key] = p
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
		packs[p.name] = &pcopy
		for n, _ := range pcopy.funcs {
			funcs[n] = append(funcs[n], &pcopy)
		}
	}
	return &index{packs: packs, funcs: funcs}
}
