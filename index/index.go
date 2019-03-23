package index

import (
	"path"

	"github.com/2hdddg/gop/parser"
	"github.com/2hdddg/gop/tree"
)

type Package struct {
	Path string
	Name string
}

type Hit struct {
	Package  *Package
	Filename string
	Line     int
	Extra    string // Depends on type of hit
}

func (h *Hit) Path() string {
	return path.Join(h.Package.Path, h.Filename)
}

type Index struct {
	RootPath  string
	functions map[string][]*Hit
	methods   map[string][]*Hit
}

type Query struct {
	Name     string
	Imported []string
}

type Result struct {
	Functions []*Hit
	Methods   []*Hit
}

func NewQuery(name string) *Query {
	return &Query{
		Name:     name,
		Imported: []string{},
	}
}

func toHit(p *Package, f *tree.File, s *parser.Symbol) *Hit {
	return &Hit{
		Package:  p,
		Filename: f.Name,
		Line:     s.Line,
		Extra:    s.Object,
	}
}

func appendToMap(key string, h *Hit, m map[string][]*Hit) {
	l := m[key]
	l = append(l, h)
	m[key] = l
}

func (i *Index) add(p *tree.Package) {
	ip := &Package{
		Path: p.Path,
		Name: p.Name,
	}
	for _, f := range p.Files {
		for _, s := range f.Syms.Functions {
			appendToMap(s.Name, toHit(ip, f, &s), i.functions)
		}
		for _, s := range f.Syms.Methods {
			appendToMap(s.Name, toHit(ip, f, &s), i.methods)
		}
	}
}

func (i *Index) traverse(p *tree.Package) {
	i.add(p)
	for _, sp := range p.Packs {
		i.traverse(sp)
	}
}

func Build(tree *tree.Tree) *Index {
	i := Index{
		RootPath:  tree.Path,
		functions: map[string][]*Hit{},
		methods:   map[string][]*Hit{},
	}
	for _, p := range tree.Packs {
		i.traverse(p)
	}
	return &i
}

func importFilter(hits []*Hit, imported []string) []*Hit {
	filtered := make([]*Hit, 0, len(hits))
	for _, h := range hits {
		for _, i := range imported {
			if h.Package.Name == i {
				filtered = append(filtered, h)
				continue
			}
		}
	}
	return filtered
}

func (i *Index) Query(q *Query) *Result {
	funcs := i.functions[q.Name]
	meths := i.methods[q.Name]

	if len(q.Imported) > 0 {
		funcs = importFilter(funcs, q.Imported)
		meths = importFilter(meths, q.Imported)
	}

	return &Result{
		Functions: funcs,
		Methods:   meths,
	}
}
