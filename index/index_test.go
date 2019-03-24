package index

import (
	"fmt"
	"testing"

	p "github.com/2hdddg/gop/parser"
	"github.com/2hdddg/gop/tree"
)

type ParserFake struct {
	syms *p.Symbols
	err  error
}

func (p *ParserFake) Parse(path string) (*p.Symbols, error) {
	return p.syms, p.err
}

func symsOneOfEach() *p.Symbols {
	syms := p.NewSymbols()
	syms.Functions = []p.Symbol{
		{"Func1", 666, "", ""},
		{"Func2", 766, "", ""},
		{"Method1", 100, "Struct1", "struct"},
		{"Method2", 123, "Struct2", "struct"},
	}
	syms.Structs = []p.Symbol{
		{"Struct1", 50, "", ""},
		{"Struct2", 60, "", ""},
	}
	syms.Interfaces = []p.Symbol{
		{"Interface1", 75, "", ""},
		{"Interface2", 75, "", ""},
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
		s []string // filepath:line:extra
	}{
		{
			"Function, no import scope",
			Query{
				Name: "Func1",
			},
			[]string{
				"/x/pack1/xfile1:666: func",
				"/y/pack2/yfile2:666: func",
			},
		},
		{
			"Function, import scope",
			Query{
				Name:     "Func1",
				Imported: []string{"y/pack2"},
			},
			[]string{
				"/y/pack2/yfile2:666: func",
			},
		},
		{
			"No hit",
			Query{
				Name: "xyz",
			},
			[]string{},
		},
		{
			"Package, no scope",
			Query{
				Name: "pack1",
			},
			[]string{"/x/pack1:0: package"},
		},
		{
			"Package, import scope (should not matter)",
			Query{
				Name:     "pack1",
				Imported: []string{"y/pack2"},
			},
			[]string{"/x/pack1:0: package"},
		},
	}
	for _, c := range cases {
		hits := i.Query(&c.q)

		if len(hits) != len(c.s) {
			t.Errorf("%s: expected %v number of hits but was %v",
				c.d, len(c.s), len(hits))
		}

		for x, s := range c.s {
			hit := hits[x]
			xs := fmt.Sprintf("%s:%v:%s",
				hit.Path(), hit.Line, hit.Extra)
			if xs != s {
				t.Errorf("%s: expected hit %v to be %v but was %v",
					c.d, x, s, xs)
			}
		}
	}
}
