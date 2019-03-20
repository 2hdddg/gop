package parser

import (
	"testing"
)

func assertLen(t *testing.T, expected, actual int) {
	if expected != actual {
		t.Fatalf("Expected len %d but was %d", expected, actual)
	}
}

func (a *Symbol) assert(t *testing.T, e *Symbol) {
	if e.Name != a.Name {
		t.Errorf("Expected name %v but was %v", e.Name, a.Name)
	}
	if e.Line != a.Line {
		t.Errorf("Expected line %v but was %v", e.Line, a.Line)
	}
	if e.Object != a.Object {
		t.Errorf("Expected object %v but was %v", e.Object, a.Object)
	}
}

func TestParsingOfFunctions(t *testing.T) {
	o := NewSymbols()
	c := `package x
		func Exported() {
		}`
	_ = o.Parse(c)

	assertLen(t, len(o.Functions), 1)
	o.Functions[0].assert(t, &Symbol{
		Name: "Exported",
		Line: 2,
	})
}

func TestParsingOfStructs(t *testing.T) {
	o := NewSymbols()
	c := `package x
		type AStruct struct {
			s string
		}`
	_ = o.Parse(c)

	assertLen(t, len(o.Structs), 1)
	o.Structs[0].assert(t, &Symbol{
		Name: "AStruct",
		Line: 2,
	})
}

func TestParsingOfMethods(t *testing.T) {
	o := NewSymbols()
	c := `package x
		type AStruct struct {
			s string
		}

		func (a AStruct) ExportedOnAStruct() {
		} 

		func (a *AStruct) ExportedOnAStructPtr() {
		}`
	_ = o.Parse(c)

	assertLen(t, 2, len(o.Methods))
	o.Methods[0].assert(t, &Symbol{
		Name:   "ExportedOnAStruct",
		Line:   6,
		Object: "AStruct",
	})
	o.Methods[1].assert(t, &Symbol{
		Name:   "ExportedOnAStructPtr",
		Line:   9,
		Object: "AStruct",
	})
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
