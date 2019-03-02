package server

type builder struct {
	packs map[string]*pack
}

func newBuilder() *builder {
	packs := make(map[string]*pack)
	return &builder{packs: packs}
}

func (b *builder) ensurePackage(f *file) *pack {
	path := f.packPath()
	name := f.packName
	key := path + name

	p, exists := b.packs[key]
	if !exists {
		p = newPackage(path)
		b.packs[key] = p
	}

	return p
}

func (b *builder) add(f *file) {
	p := b.ensurePackage(f)
	p.mergeFile(f)
}

func (b *builder) build() *index {
	packs := make(map[string]*pack)
	funcs := make(map[string][]*pack)

	for _, p := range b.packs {
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
