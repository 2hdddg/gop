package parser

import "testing"

var (
	o *Symbols
	e error
	c string
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
}

func assertName(t *testing.T, expected, actual string) {
	if expected != actual {
		t.Errorf("Expected name %v but was %v", actual, expected)
	}
}

func assertLine(t *testing.T, expected, actual int) {
	if expected != actual {
		t.Errorf("Expected line %v but was %v", expected, actual)
	}
}

func TestParsingOfFunctions(t *testing.T) {
	o = NewSymbols()
	c = `package x
		func Exported() {
		}`
	e = o.Parse(c)

	assertLen(t, len(o.Functions), 1)
	f := &o.Functions[0]
	assertName(t, f.Name, "Exported")
	assertLine(t, f.Line, 2)
}

func TestParsingOfStructs(t *testing.T) {
	o = NewSymbols()
	c = `package x
		type AStruct struct {
			s string
		}`
	e = o.Parse(c)

	assertLen(t, len(o.Structs), 1)
	s := &o.Structs[0]
	assertName(t, s.Name, "AStruct")
	assertLine(t, s.Line, 2)
}

func TestParsingOfMethods(t *testing.T) {
	o = NewSymbols()
	c = `package x
		type AStruct struct {
			s string
		}

		func (a AStruct) ExportedOnAStruct() {
		} 

		func (a *AStruct) ExportedOnAStructPtr() {
		}`
	e = o.Parse(c)

	assertLen(t, 2, len(o.Methods))
	m := &o.Methods[0]
	m.assert(t, &Symbol{
		Name: "ExportedOnAStruct",
		Line: 6,
	})
	assertName(t, m.Object, "AStruct")
}
