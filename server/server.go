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

	builder, err := tree.NewBuilder(config.SystemPath)
	tree, err := builder.Build()
	service.NewOrUpdatedTree(tree)

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
