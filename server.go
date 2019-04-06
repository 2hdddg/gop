package main

import (
	"flag"
	"log"
	"net"
	"net/http"
	"net/rpc"
	"os"
	"strconv"

	"github.com/2hdddg/gop/config"
	service "github.com/2hdddg/gop/service"
)

func serve() {
	var port int

	flags := flag.NewFlagSet("serve", flag.ExitOnError)
	flags.IntVar(&port, "port", 8080, "Server port")
	flags.Parse(os.Args[2:])

	config := config.NewConfig()

	srv := service.NewService()
	err := srv.Start()
	if err != nil {
		log.Fatalf("Failed to start search service: %s", err)
	}

	log.Printf("Starting server on port %d", port)
	rpc.HandleHTTP()
	l, err := net.Listen("tcp", ":"+strconv.Itoa(port))
	if err != nil {
		log.Fatalf("Failed to listen on tcp port %d: %s",
			port, err)
	}

	// Build default indexes
	for _, root := range config.Paths() {
		srv.Index(&service.IndexReq{Path: root}, &service.IndexRes{})
	}

	err = http.Serve(l, nil)
	if err != nil {
		log.Fatalf("Failed to start http server: %s", err)
	}
}
