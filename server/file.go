package server

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"path"
	"unicode"
	"unicode/utf8"
)

type file struct {
	path     string
	packName string
	funcs    map[string]Location
}

func (f *file) packPath() string {
	return path.Dir(f.path)
}

func isExported(name string) bool {
	r, _ := utf8.DecodeRuneInString(name)
	return unicode.IsUpper(r)
}

func parseFunc(fset *token.FileSet, o *ast.Object) Location {
	position := fset.Position(o.Pos())
	return Location{Line: position.Line, Column: position.Column}
}

func parseFile(path string) *file {
	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, path, nil, 0)
	if err != nil {
		fmt.Printf("Error while parsing %v: %v\n", path, err)
		return nil
	}

	packageName := f.Name.Name
	funcs := make(map[string]Location)
	for _, o := range f.Scope.Objects {
		if !isExported(o.Name) {
			continue
		}
		if o.Kind == ast.Fun {
			funcs[o.Name] = parseFunc(fset, o)
		}
	}

	parsed := file{path: path, packName: packageName, funcs: funcs}
	return &parsed
}
