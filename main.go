package main

import (
	"flag"
	"log"

	"github.com/2hdddg/gop/client"
	"github.com/2hdddg/gop/config"
	"github.com/2hdddg/gop/server"
)

var (
	isServer bool
	port     int
	params   client.Params
)

func setupParameters() {
	flag.BoolVar(&isServer, "serve", false, "Run server")
	flag.IntVar(&port, "port", 8080, "Server port")

	flag.StringVar(&params.Name, "name", "", "Find function")
	flag.StringVar(&params.FilePath, "file", "", "Go file")

	flag.Parse()
}

func main() {
	setupParameters()
	config := config.NewConfig()

	// Configure standard logger
	log.SetFlags(log.Ltime | log.Lshortfile)

	if isServer {
		server.Run(config, port)
		return
	}

	client.Run(config, port, &params)
}

/*
:lgetexpr system("~/code/go_parser/gop --func " . expand('<cword>') . " --file " . expand('%:p')) | lopen
*/
