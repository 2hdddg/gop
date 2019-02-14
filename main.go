package main

import (
	"flag"
	"github.com/2hdddg/gop/server"
)

var (
	isServer bool
	port     int
)

func main() {
	flag.BoolVar(&isServer, "serve", false, "Run server")
	flag.IntVar(&port, "port", 8080, "Server port")
	flag.Parse()

	if isServer {
		server.Run(port)
	}
}
