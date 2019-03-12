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
	assertName(t, f.base.Name, "Exported")
	assertLine(t, f.base.Line, 2)
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
	assertName(t, s.base.Name, "AStruct")
	assertLine(t, s.base.Line, 2)
}

func TestParsingOfMethods(t *testing.T) {
	o = NewSymbols()
	c = `package x
		type AStruct struct {
			s string
		}

		func (a AStruct) ExportedOnAStruct() {
		} `
	e = o.Parse(c)

	assertLen(t, len(o.Methods), 1)
	m := &o.Methods[0]
	assertName(t, m.base.Name, "ExportedOnAStruct")
	assertLine(t, 6, m.base.Line)
	assertName(t, m.Object, "AStruct")
}
