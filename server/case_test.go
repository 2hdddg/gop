package server

import "testing"

func buildIndex() *index {
	t := newTree("root/path")

	// First file with funcs: func1 & func2
	f := newFile("pack1", "pack1/file1")
	f.funcs["func1"] = append(f.funcs["func1"], funcDescr{line: 1})
	f.funcs["func2"] = append(f.funcs["func2"], funcDescr{line: 30})
	t.addFile(f)

	// Second file in another package with funcs: func2 & func3
	f = newFile("pack2", "pack2/file2")
	f.funcs["func1"] = append(f.funcs["func1"], funcDescr{line: 61})
	f.funcs["func2"] = append(f.funcs["func2"], funcDescr{line: 90})
	t.addFile(f)

	return t.buildIndex()
}

func find(a *Answer, path string) *Location {
	for _, l := range a.Locations {
		if l.Path == path {
			return &l
		}
	}
	return nil
}

// Verifies that functions can be found
func TestFindFunc(t *testing.T) {
	i := buildIndex()
	q := Query{Object: Function, Name: "func2"}
	a := &Answer{}
	i.funcByQuery(&q, a)

	if len(a.Locations) != 2 {
		t.Errorf("Should have found 2 definitions, found %d",
			len(a.Locations))
	}

	l := find(a, "pack1/file1")
	if l.Line != 30 {
		t.Errorf("Should have func2 at line 30")
	}

	l = find(a, "pack2/file2")
	if l.Line != 90 {
		t.Errorf("Should have func2 at line 90")
	}
}

// Verifies that packages can be found
func TestFindPack(t *testing.T) {
	i := buildIndex()
	a := &Answer{}
	i.packByName("pack1", a)

	p := find(a, "root/path/pack1")
	if p == nil {
		t.Error("Didn't find pack1")
	}
}
