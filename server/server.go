package server

import (
	"github.com/2hdddg/gop/config"
	"github.com/2hdddg/gop/search"
	"github.com/2hdddg/gop/tree"

	"log"
	"net"
	"net/http"
	"net/rpc"
	"strconv"
)

/*
type Object int

const (
	Function Object = 0
	Package  Object = 1
)

type Query struct {
	Object   Object
	Name     string
	Packages []string
	Config   *config.Config
}

type Location struct {
	Path   string
	Line   int
	Column int
}

type Answer struct {
	Locations []Location
	Errors    []string
}
*/

func Run(config *config.Config, port int) {
	service := search.NewService()
	client := service.Start()
	err := rpc.RegisterName("Search", client)
	if err != nil {
		log.Fatalf("Failed to register search service: %s", err)
	}

	builder, err := tree.NewBuilder(config.SystemPath)
	tree, err := builder.Build()

	service.NewOrUpdatedTree(tree)

	//go monitor(config, search.indexChan)

	log.Printf("Starting server on port %d", port)
	rpc.HandleHTTP()
	l, err := net.Listen("tcp", ":"+strconv.Itoa(port))
	if err != nil {
		log.Fatalf("Failed to listen on tcp port %d: %s",
			port, err)
	}
	err = http.Serve(l, nil)
	if err != nil {
		log.Fatalf("Failed to start http server: %s", err)
	}
}
