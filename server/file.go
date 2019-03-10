package server

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
)

type line int

type file struct {
	path     string // Full path to file
	packName string // Qualified name: go/ast
	funcs    map[string]line
}

func parseFunc(fset *token.FileSet, o *ast.Object) line {
	return line(fset.Position(o.Pos()).Line)
}

func parseFile(packName, path string) (*file, error) {
	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, path, nil, 0)
	if err != nil {
		return nil, fmt.Errorf("Error while parsing %v: %v\n", path, err)
	}

	ast.FileExports(f)
	funcs := make(map[string]line)
	for _, o := range f.Scope.Objects {
		if o.Kind == ast.Fun {
			funcs[o.Name] = parseFunc(fset, o)
		}
	}

	parsed := file{
		path:     path,
		packName: packName,
		funcs:    funcs,
	}
	return &parsed, nil
}
