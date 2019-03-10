package client

import (
	"fmt"
	"github.com/2hdddg/gop/shared"
	"go/parser"
	"go/token"
	"path/filepath"
	"strconv"
	"strings"
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

func parseFilePackage(config *shared.Config,
	path string) (pack string, e error) {
	p, _ := filepath.Abs(path)
	p = filepath.Dir(p)
	p = strings.TrimPrefix(p, config.WorkspacePath)
	pack = strings.TrimPrefix(p, "/")
	return
}
