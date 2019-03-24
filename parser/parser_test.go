package parser

import (
	"testing"
)

func assertLen(t *testing.T, expected, actual int) {
	if expected != actual {
		t.Fatalf("Expected len %d but was %d", expected, actual)
	}
}

func TestParse(t *testing.T) {
	cases := []struct {
		desc string
		code string
		syms Symbols
		err  error
	}{
		{
			"Exported func",
			`package x
				func Exported() {
				}`,
			Symbols{
				Functions: []Symbol{
					{"Exported", 2, "", ""},
				},
			},
			nil,
		},
		{
			"Exported struct",
			`package x
				type AStruct struct {
					s string
				}`,
			Symbols{
				Structs: []Symbol{
					{"AStruct", 2, "", ""},
				},
			},
			nil,
		},
		{
			"Methods on struct",
			`package x
				type AStruct struct {
					s string
				}

				func (a AStruct) ExportedOnAStruct() {
				}

				func (a *AStruct) ExportedOnAStructPtr() {
				}`,
			Symbols{
				Functions: []Symbol{
					{"ExportedOnAStruct", 6, "AStruct", ""},
					{"ExportedOnAStructPtr", 9, "*AStruct", ""},
				},
				Structs: []Symbol{
					{"AStruct", 2, "", ""},
				},
			},
			nil,
		},
		{
			"Exported interface",
			`package x
			type AInterface interface {
				Meth(x string)
			}`,
			Symbols{
				Interfaces: []Symbol{
					{"AInterface", 2, "", ""},
				},
			},
			nil,
		},
	}

	assertSymbol := func(desc string, a, e *Symbol) {
		if a.Name != e.Name {
			t.Errorf("%s: Expected symbol name %v but was %v (%v)",
				desc, e.Name, a.Name, a)
		}
		if a.Line != e.Line {
			t.Errorf("%s: Expected symbol line %v but was %v (%v)",
				desc, e.Line, a.Line, a)
		}
		if a.Parent != e.Parent {
			t.Errorf("%s: Expected symbol parent %v but was %v (%v)",
				desc, e.Parent, a.Parent, a)
		}
	}

	assertSymbols := func(desc, objtype string, a, e []Symbol) {
		if len(a) != len(e) {
			t.Errorf("%s: Expected %v number of %v but was %v",
				desc, len(e), objtype, len(a))
			return
		}
		for i, m := range a {
			assertSymbol(desc, &m, &e[i])
		}
	}

	for _, c := range cases {
		syms := NewSymbols()
		err := syms.Parse(c.code)

		if err != c.err {
			t.Errorf("%s: Expected error to be %v but was %v",
				c.desc, c.err, err)
		}
		assertSymbols(c.desc, "methods",
			syms.Functions, c.syms.Functions)
		assertSymbols(c.desc, "funcs",
			syms.Functions, c.syms.Functions)
		assertSymbols(c.desc, "structs",
			syms.Structs, c.syms.Structs)
		assertSymbols(c.desc, "interfaces",
			syms.Interfaces, c.syms.Interfaces)
	}
}

func TestParsingOfImports(t *testing.T) {
	c := `package x
		  import (
			"main/sub"
			"another"
		  )
		}`
	imports, _ := parseImports(c)
	assertLen(t, 2, len(imports))
	expected := []string{"main/sub", "another"}
	for i, v := range imports {
		if v != expected[i] {
			t.Errorf("Expected %v but was %v", expected, i)
			break
		}
	}
}
