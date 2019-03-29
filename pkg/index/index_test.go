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
	func1   = p.Symbol{"Func1", p.Function, 666, "", p.Undefined}
	func2   = p.Symbol{"Func2", p.Function, 766, "", p.Undefined}
	meth1   = p.Symbol{"Method1", p.Function, 100, "Struct1", p.Struct}
	meth2   = p.Symbol{"Method2", p.Function, 123, "Struct2", p.Struct}
	struct1 = p.Symbol{"Struct1", p.Struct, 50, "", p.Undefined}
	struct2 = p.Symbol{"Struct2", p.Struct, 60, "", p.Undefined}
	intf1   = p.Symbol{"Interface1", p.Interface, 75, "", p.Undefined}
	intf2   = p.Symbol{"Interface2", p.Interface, 75, "", p.Undefined}
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
