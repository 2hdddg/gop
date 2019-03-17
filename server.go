package main

import (
	"log"
	"path"

	"github.com/2hdddg/gop/config"
	"github.com/2hdddg/gop/index"
	"github.com/2hdddg/gop/tree"
)

func run_server(config *config.Config, port int) {
	log.Println(config)
	builder, err := tree.NewBuilder(config.SystemPath)
	if err != nil {
		log.Println(err)
		return
	}
	tree, err := builder.Build()
	if err != nil {
		log.Println(err)
		return
	}

	i := index.Build(tree)
	query := index.NewQuery("Build")
	funcs, methods := i.Query(query)
	for _, h := range funcs {
		log.Printf("%s:%d Function definition",
			path.Join(h.Package.Path, h.Filename), h.Line)
	}
	for _, h := range methods {
		log.Printf("Method: %V on %s", h, h.Extra)
	}
}
