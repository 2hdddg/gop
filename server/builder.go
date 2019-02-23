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
		p = newPackage(path, name)
		b.packs[key] = p
	}

	return p
}

func (b *builder) add(f *file) {
	p := b.ensurePackage(f)
	p.merge(f)
}

func (b *builder) build() index {
	packs := make(map[string]pack)
	for _, p := range b.packs {
		packs[p.name] = *p
	}
	return index{packs: packs}
}
