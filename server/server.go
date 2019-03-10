package server

import (
	"github.com/2hdddg/gop/shared"
	"log"
	"net"
	"net/http"
	"net/rpc"
	"strconv"
)

type Object int

const (
	Function Object = 0
	Package  Object = 1
)

type Query struct {
	Object   Object
	Name     string
	Packages []string
	Config   *shared.Config
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

func Run(port int) {
	config := shared.NewConfig()
	if !config.Valid() {
		log.Fatalln("Invalid config")
	}

	search := newSearch()
	go search.thread()
	err := rpc.RegisterName("Search", search)
	if err != nil {
		log.Fatalf("Failed to register search service: %s", err)
	}

	go monitor(config, search.indexChan)

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
