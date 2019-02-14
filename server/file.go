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

type Location struct {
	Line   int
	Column int
}

type File struct {
	Path    string
	Package string
	Funcs   map[string]Location
}

func (f *File) PackagePath() string {
	return path.Dir(f.Path)
}

func isExported(name string) bool {
	r, _ := utf8.DecodeRuneInString(name)
	return unicode.IsUpper(r)
}

func parseFunc(fset *token.FileSet, o *ast.Object) Location {
	position := fset.Position(o.Pos())
	return Location{Line: position.Line, Column: position.Column}
}

func parseFile(path string) *File {
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

	parsed := File{Path: path, Package: packageName, Funcs: funcs}
	return &parsed
}
