package server

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
)

type line int

type file struct {
	filePath string // Full path to file
	packName string // Qualified name: go/ast
	funcs    map[string]line
}

func parseFunc(fset *token.FileSet, o *ast.Object) line {
	return line(fset.Position(o.Pos()).Line)
}

func parseFile(packName, filePath string) (*file, error) {
	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, filePath, nil, 0)
	if err != nil {
		return nil, fmt.Errorf("Error while parsing %v: %v\n", filePath, err)
	}
	packageName := packName

	ast.FileExports(f)
	funcs := make(map[string]line)
	for _, o := range f.Scope.Objects {
		if o.Kind == ast.Fun {
			funcs[o.Name] = parseFunc(fset, o)
		}
	}

	parsed := file{
		filePath: filePath,
		packName: packageName,
		funcs:    funcs}
	return &parsed, nil
}
