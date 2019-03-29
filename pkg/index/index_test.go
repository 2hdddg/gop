package index

import (
	"fmt"
	"testing"

	p "github.com/2hdddg/gop/pkg/parser"
	"github.com/2hdddg/gop/pkg/tree"
)

type ParserFake struct {
	syms *p.Symbols
	err  error
}

func (p *ParserFake) Parse(path string) (*p.Symbols, error) {
	return p.syms, p.err
}

var (
	func1 = p.Symbol{
		Name: "Func1", Type: p.Function, Line: 666,
		ContextName: "", ContextType: p.Undefined}
	func2 = p.Symbol{
		Name: "Func2", Type: p.Function, Line: 766,
		ContextName: "", ContextType: p.Undefined}
	meth1 = p.Symbol{
		Name: "Method1", Type: p.Function, Line: 100,
		ContextName: "Struct1", ContextType: p.Struct}
	meth2 = p.Symbol{
		Name: "Method2", Type: p.Function, Line: 123,
		ContextName: "Struct2", ContextType: p.Struct}
	struct1 = p.Symbol{
		Name: "Struct1", Type: p.Struct, Line: 50,
		ContextName: "", ContextType: p.Undefined}
	struct2 = p.Symbol{
		Name: "Struct2", Type: p.Struct, Line: 60,
		ContextName: "", ContextType: p.Undefined}
	intf1 = p.Symbol{
		Name: "Interface1", Type: p.Interface, Line: 75,
		ContextName: "", ContextType: p.Undefined}
	intf2 = p.Symbol{
		Name: "Interface2", Type: p.Interface, Line: 75,
		ContextName: "", ContextType: p.Undefined}
)

func symsOneOfEach() *p.Symbols {
	syms := p.NewSymbols()
	syms.List = []p.Symbol{
		func1,
		func2,
		meth1,
		meth2,
		struct1,
		struct2,
		intf1,
		intf2,
	}
	return syms
}

func TestIndexQuery(t *testing.T) {
	// Tree
	tree, err := tree.NewTree("/")
	if err != nil {
		t.Fatalf("Failed to create tree: %v", err)
	}

	// Add packages
	pack1 := tree.AddPackage("x/pack1")
	pack2 := tree.AddPackage("y/pack2")
	p1 := &Package{pack1.Path, pack1.Name}
	p2 := &Package{pack2.Path, pack2.Name}

	// Add files
	pars := &ParserFake{
		syms: symsOneOfEach(),
	}
	// Both packages gets the same set of symbols
	pack1.AddFile("xfile1", pars)
	pack2.AddFile("yfile2", pars)

	// Build index from tree
	i := Build(tree)

	cases := []struct {
		d string
		q Query
		h []Hit
		p []Package
	}{
		{
			"Function, no import scope",
			Query{
				Name: "Func1",
			},
			[]Hit{
				{p1, "xfile1", func1},
				{p2, "yfile2", func1},
			},
			[]Package{},
		},
		{
			"Function, import scope",
			Query{
				Name:     "Func1",
				Imported: []string{"y/pack2"},
			},
			[]Hit{
				{p2, "yfile2", func1},
			},
			[]Package{},
		},
		{
			"No hit",
			Query{
				Name: "xyz",
			},
			[]Hit{},
			[]Package{},
		},
		{
			"Package, no scope",
			Query{
				Name: "pack1",
			},
			[]Hit{},
			[]Package{*p1},
		},
		{
			"Package, import scope (should not matter)",
			Query{
				Name:     "pack1",
				Imported: []string{"y/pack2"},
			},
			[]Hit{},
			[]Package{*p1},
		},
	}
	for _, c := range cases {
		hits := make([]Hit, 0)
		phits := make([]Package, 0)
		i.Query(&c.q,
			func(h Hit) {
				hits = append(hits, h)
			},
			func(p Package) {
				phits = append(phits, p)
			})

		if len(hits) != len(c.h) {
			t.Errorf("%s: expected %v number of hits but was %v",
				c.d, len(c.h), len(hits))
		}
		if len(phits) != len(c.p) {
			t.Errorf("%s: expected %v number of package hits but was %v",
				c.d, len(c.p), len(phits))
		}
		for x, s := range c.h {
			hit := hits[x]
			xs := fmt.Sprintf("%s:%s",
				hit.Path(), hit.Symbol.ToString())
			as := fmt.Sprintf("%s:%s",
				s.Path(), hit.Symbol.ToString())
			if xs != as {
				t.Errorf("%s: expected hit %v to be %v but was %v",
					c.d, x, as, xs)
			}
		}
		for x, pa := range c.p {
			pe := phits[x]
			as := fmt.Sprintf("%s:%s", pa.Path, pa.Name)
			es := fmt.Sprintf("%s:%s", pe.Path, pe.Name)
			if as != es {
				t.Errorf("%s: expected hit %v to be %v but was %v",
					c.d, x, as, es)
			}
		}
	}
}
