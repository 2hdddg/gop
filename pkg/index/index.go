package index

import (
	"path"
	"strings"

	"github.com/2hdddg/gop/pkg/parser"
	"github.com/2hdddg/gop/pkg/tree"
)

type Package struct {
	Path string
	Name string
}

type Hit struct {
	Package  *Package
	Filename string
	Symbol   parser.Symbol
}

func (h *Hit) Path() string {
	return path.Join(h.Package.Path, h.Filename)
}

type Index struct {
	RootPath  string
	functions map[string][]*Hit
	structs   map[string][]*Hit
	interfs   map[string][]*Hit
	fields    map[string][]*Hit
	packs     map[string][]*Package
}

type Query struct {
	Name     string
	Imported []string
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
		Symbol:   *s,
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
		for _, s := range f.Syms.List {
			h := toHit(ip, f, &s)
			// Put in different maps for reduced size and for
			// simpler impl of queries filtered by type.
			switch s.Type {
			case parser.Method:
			case parser.Function:
				appendToMap(s.Name, h, i.functions)
			case parser.Struct:
				appendToMap(s.Name, h, i.structs)
			case parser.Interface:
				appendToMap(s.Name, h, i.interfs)
			case parser.Field:
				appendToMap(s.Name, h, i.fields)
			}
		}
	}

	// Put last part of package name in index: x/y/z -> z
	// Last part is usually what's needed from code as path.Split
	parts := strings.Split(p.Name, "/")
	if len(parts) >= 1 {
		key := parts[len(parts)-1]
		l := i.packs[key]
		l = append(l, ip)
		i.packs[key] = l
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
		structs:   map[string][]*Hit{},
		interfs:   map[string][]*Hit{},
		fields:    map[string][]*Hit{},
		packs:     map[string][]*Package{},
	}
	for _, p := range tree.Packs {
		i.traverse(p)
	}
	return &i
}

func filterAndAdd(hits []*Hit, imported []string, total []*Hit) {
	for _, h := range hits {
		for _, i := range imported {
			if h.Package.Name == i {
				total = append(total, h)
				continue
			}
		}
	}
}

func add(hits []*Hit, total []*Hit) {
	for _, h := range hits {
		total = append(total, h)
	}
}

type OnHit func(h Hit)
type OnPackage func(p Package)

func (i *Index) Query(q *Query, onHit OnHit, onPack OnPackage) {
	num := 0

	noFilter := func(hits []*Hit) {
		for _, h := range hits {
			onHit(*h)
			num++
		}
	}
	importFilter := func(hits []*Hit) {
		for _, h := range hits {
			for _, i := range q.Imported {
				if h.Package.Name == i {
					onHit(*h)
					num++
					continue
				}
			}
		}
	}

	filter := noFilter
	if len(q.Imported) > 0 {
		filter = importFilter
	}

	filter(i.functions[q.Name])
	filter(i.structs[q.Name])
	filter(i.interfs[q.Name])
	filter(i.fields[q.Name])

	for _, p := range i.packs[q.Name] {
		onPack(*p)
		num++
	}
}
