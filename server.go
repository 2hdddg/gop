package main

import (
	"log"

	"github.com/2hdddg/gop/index"
	"github.com/2hdddg/gop/tree"
)

func run_server(port int) {
	builder, err := tree.NewBuilder("/home/peter/code/go_parser")
	if err != nil {
		log.Println(err)
		return
	}
	tree, err := builder.Build()
	if err != nil {
		log.Println(err)
		return
	}
	log.Println(tree)

	i := index.Build(tree)
	log.Println(&i)

	query := index.NewQuery("Build")
	hits := i.Query(query)
	for _, h := range hits {
		log.Printf("%V", h)
	}
}
