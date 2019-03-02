package server

import "testing"

func buildIndex() *index {
	b := newBuilder()

	// First file with funcs: func1 & func2
	f := file{path: "file1/path", packName: "pack1",
		funcs: make(map[string]Location)}
	f.funcs["func1"] = Location{Line: 1, Column: 1}
	f.funcs["func2"] = Location{Line: 30, Column: 1}
	b.add(&f)

	// Second file in another package with funcs: func2 & func3
	f = file{path: "file2/path", packName: "pack2",
		funcs: make(map[string]Location)}
	f.funcs["func1"] = Location{Line: 61, Column: 1}
	f.funcs["func2"] = Location{Line: 90, Column: 1}
	b.add(&f)

	return b.build()
}

func findFuncDef(a *LocationsAnswer, path string) *FileLocation {
	for _, l := range a.Locations {
		if l.FilePath == path {
			return &l
		}
	}
	return nil
}

// Verifies that an index can be built and used to find two definitons
// of a function in two different packages.
func TestFindFuncDefinition(t *testing.T) {
	i := buildIndex()
	a := i.funcDefinition("func2")

	if len(a.Locations) != 2 {
		t.Errorf("Should have found 2 definitions, found %d",
			len(a.Locations))
	}

	l := findFuncDef(a, "file1/path")
	if l.Line != 30 {
		t.Errorf("Should have func2 at line 30")
	}

	l = findFuncDef(a, "file2/path")
	if l.Line != 90 {
		t.Errorf("Should have func2 at line 90")
	}
}
