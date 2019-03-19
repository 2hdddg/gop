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

func Run(config *config.Config, port int) {
	service := search.NewService()
	client := service.Start()
	err := rpc.RegisterName("Search", client)
	if err != nil {
		log.Fatalf("Failed to register search service: %s", err)
	}

	// Build takes a while, run this in a go routine so that the
	// server is started fast.
	go func() {
		builder, err := tree.NewBuilder(config.SystemPath)
		// Report build progress to search service, this makes
		// it possible for search service to build incomplete
		// indexes to at least serve some results.
		builder.Progress = service
		_, err = builder.Build()
		if err != nil {
			log.Fatalf("Failed to build tree: %s", err)
		}
	}()

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
