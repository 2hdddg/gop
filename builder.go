package main

type Builder struct {
	packages map[string]*Package
}

func NewBuilder() *Builder {
	packages := make(map[string]*Package)
	return &Builder{packages: packages}
}

func (i *Builder) ensurePackage(f *File) *Package {
	path := f.PackagePath()
	name := f.Package
	key := path + name

	p, exists := i.packages[key]
	if !exists {
		p = NewPackage(path, name)
		i.packages[key] = p
	}

	return p
}

func (i *Builder) Add(f *File) {
	p := i.ensurePackage(f)
	p.Merge(f)
}

func (i *Builder) Build() Index {
	packages := make(map[string]Package)
	for _, p := range i.packages {
		packages[p.Name] = *p
	}
	return Index{packages: packages}
}
