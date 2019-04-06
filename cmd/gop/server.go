package main

import (
	"flag"
	"log"
	"net"
	"net/http"
	"net/rpc"
	"os"
	"strconv"

	"github.com/2hdddg/gop/pkg/config"
	indexservice "github.com/2hdddg/gop/pkg/service/index"
	searchservice "github.com/2hdddg/gop/pkg/service/search"
)

func serve() {
	var port int

	flags := flag.NewFlagSet("serve", flag.ExitOnError)
	flags.IntVar(&port, "port", 8080, "Server port")
	flags.Parse(os.Args[2:])

	config := config.NewConfig()

	searchSrv := searchservice.NewService()
	err := searchSrv.Start()
	if err != nil {
		log.Fatalf("Failed to start search service: %s", err)
	}
	// Search service implements progress interface
	indexSrv := indexservice.NewService(searchSrv)
	err = indexSrv.Start()
	if err != nil {
		log.Fatalf("Failed to start index service: %s", err)
	}

	for _, root := range config.Paths() {
		log.Println("Adding ", root)
		indexSrv.Add(root)
	}

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
