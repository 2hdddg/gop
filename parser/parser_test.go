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
				List: []Symbol{
					{"Exported", Function, 2, "", Undefined},
				},
			},
			nil,
		},
		{
			"Exported struct",
			`package x
				type AStruct struct {
					S string
				}`,
			Symbols{
				List: []Symbol{
					{"AStruct", Struct, 2, "", Undefined},
					{"S", Field, 3, "AStruct", Struct},
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
				List: []Symbol{
					{"AStruct", Struct, 2, "", Undefined},
					{"s", Field, 3, "AStruct", Struct},
					{"ExportedOnAStruct", Method, 6, "AStruct", Struct},
					{"ExportedOnAStructPtr", Method, 9, "*AStruct", Struct},
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
				List: []Symbol{
					{"AInterface", Interface, 2, "", Undefined},
					{"Meth", Method, 3, "AInterface", Interface},
				},
			},
			nil,
		},
	}

	assertSymbol := func(desc string, a, e *Symbol) {
		if a.Type != e.Type {
			t.Errorf("%s: Expected symbol type %v but was %v (%v)",
				desc, e.Type, a.Type, a)
		}
		if a.Name != e.Name {
			t.Errorf("%s: Expected symbol name %v but was %v (%v)",
				desc, e.Name, a.Name, a)
		}
		if a.Line != e.Line {
			t.Errorf("%s: Expected symbol line %v but was %v (%v)",
				desc, e.Line, a.Line, a)
		}
		if a.ContextName != e.ContextName {
			t.Errorf("%s: Expected symbol contextname %v but was %v (%v)",
				desc, e.ContextName, a.ContextName, a)
		}

	}

	assertSymbols := func(desc string, a, e []Symbol) {
		if len(a) != len(e) {
			t.Errorf("%s: Expected %v number but was %v",
				desc, len(e), len(a))
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
		assertSymbols(c.desc, syms.List, c.syms.List)
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
