package index

import (
	"testing"

	"github.com/2hdddg/gop/parser"
	"github.com/2hdddg/gop/tree"
)

type ParserFake struct {
	syms *parser.Symbols
	err  error
}

func (p *ParserFake) Parse(path string) (*parser.Symbols, error) {
	return p.syms, p.err
}

func (e *Hit) assert(t *testing.T, a *Hit) {
	if e.Filename != a.Filename {
		t.Errorf("Expected Hit filename to be %v but was %v",
			e.Filename, a.Filename)
	}
	if e.Line != a.Line {
		t.Errorf("Expected Hit line to be %v but was %v",
			e.Line, a.Line)
	}
	if e.Path() != a.Path() {
		t.Errorf("Expeced Hit path to be %v but was %v",
			e.Path(), a.Path())
	}
}

// Covers base functionality by having a really simple tree
// and simple query.
func TestBuildAndQueryBaseline(t *testing.T) {
	// Initialize parser to return nothing
	pars := &ParserFake{
		syms: parser.NewSymbols(),
		err:  nil,
	}
	// Setup parser result
	s := pars.syms
	s.Functions = append(s.Functions, parser.Symbol{
		Name: "Func1",
		Line: 666,
	})

	// Build tree with parsed data
	tree, err := tree.NewTree(".")
	if err != nil {
		t.Fatalf("Failed to create tree: %v", err)
	}
	pack := tree.AddPackage("x/pack")
	_, _ = pack.AddFile("thefile", pars)

	// Build index from tree
	i := Build(tree)

	// Query index
	q := &Query{
		Name: "Func1",
	}
	res := i.Query(q)

	if len(res.Functions) != 1 {
		t.Fatalf("Query should return 1 func")
	}

	res.Functions[0].assert(t, &Hit{
		Package: &Package{
			Path: pack.Path,
			Name: "pack",
		},
		Filename: "thefile",
		Line:     666,
	})

	// Same data but with query that specifies imports, should limit
	// result to imported packages.
	q.Imported = []string{"another_pack"}
	res = i.Query(q)
	if len(res.Functions) != 0 {
		t.Errorf("Query with import scope should return 0 func")
	}
	// And now the "correct" package
	q.Imported = []string{"x/pack"}
	res = i.Query(q)
	if len(res.Functions) != 1 {
		t.Errorf("Query with import scope should return 1 func")
	}

	// Query for the last part of package name
	q.Imported = []string{} // Shouldn't matter
	q.Name = "pack"
	res = i.Query(q)
	if len(res.Packages) != 1 {
		t.Errorf("Query for package should return 1 pack")
	}
}
