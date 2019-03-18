package index

import (
	"path"

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
	rootPath  string
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

func (i *Index) add(p *tree.Package) {
	ip := &Package{
		Path: p.Path,
		Name: p.Name,
	}
	for _, f := range p.Files {
		// Add functions
		for _, s := range f.Syms.Functions {
			h := &Hit{
				Package:  ip,
				Filename: f.Name,
				Line:     s.Line,
			}
			key := s.Name
			funcs := i.functions[key]
			funcs = append(funcs, h)
			i.functions[key] = funcs
		}
		for _, s := range f.Syms.Methods {
			h := &Hit{
				Package:  ip,
				Filename: f.Name,
				Line:     s.Line,
				Extra:    s.Object,
			}
			key := s.Name
			methods := i.methods[key]
			methods = append(methods, h)
			i.methods[key] = methods
		}
	}
}

func (i *Index) traverse(p *tree.Package) {
	i.add(p)
	for _, sp := range p.Packs {
		i.traverse(sp)
	}
}

func Build(tree *tree.Tree) Index {
	i := Index{
		rootPath:  tree.Path,
		functions: map[string][]*Hit{},
		methods:   map[string][]*Hit{},
	}
	for _, p := range tree.Packs {
		i.traverse(p)
	}
	return i
}

func (i *Index) Query(q *Query) *Result {
	return &Result{
		Functions: i.functions[q.Name],
		Methods:   i.methods[q.Name],
	}
}
