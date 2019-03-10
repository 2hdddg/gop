package server

import "testing"

func TestFindFuncInPack(t *testing.T) {
	p := newPack("p1")
	p.funcs["f1"] = Location{Path: "f1path"}

	// Existing
	found := p.findFunc("f1")
	if found.Path != "f1path" {
		t.Error("Should have found f1")
	}

	// Not existing
	found = p.findFunc("!f1")
	if found != nil {
		t.Error("Should not have found !f1")
	}
}

func TestMergePack(t *testing.T) {
	p := newPack("pack1")
	f := file{
		path:     "file1",
		packName: "pack1",
		funcs:    make(map[string][]funcDescr),
	}
	f.funcs["func1"] = append(f.funcs["func1"], funcDescr{line: 1})

	// Merge to package containing no functions
	p.mergeFile(&f)
	if p.findFunc("func1") == nil {
		t.Error("Didn't find func1")
	}

	// Merge same thing to package again
	p.mergeFile(&f)
	if p.findFunc("func1") == nil {
		t.Error("Didn't find func1 second time")
	}

	// Same file but different function, should remove previous
	delete(f.funcs, "func1")
	f.funcs["func2"] = append(f.funcs["func2"], funcDescr{line: 10})
	p.mergeFile(&f)
	if p.findFunc("func1") != nil {
		t.Error("Shouldn't find func1")
	}

	// Same file but different location, should update location
	f.funcs["func2"] = append(f.funcs["xxx"], funcDescr{line: 20})
	p.mergeFile(&f)
	if p.findFunc("func2").Line != 20 {
		t.Error("Should update line")
	}

	// Another file with another func, all funcs should exist
	f = *newFile("pack1", "file2")
	f.funcs["func3"] = append(f.funcs["func3"], funcDescr{line: 7})
	p.mergeFile(&f)
	if p.findFunc("func2") == nil || p.findFunc("func3") == nil {
		t.Error("Should find funcs")
	}
}
