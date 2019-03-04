package client

import (
	"fmt"
	"go/parser"
	"go/token"
	"strconv"
)

func parseFileImports(path string) (packs []string, e error) {
	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, path, nil, parser.ImportsOnly)
	if err != nil {
		return nil, fmt.Errorf("Error while parsing %v: %v\n", path, err)
	}

	for _, i := range f.Imports {
		name, _ := strconv.Unquote(i.Path.Value)
		packs = append(packs, name)
	}

	return packs, nil
}
