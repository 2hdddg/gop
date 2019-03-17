package main

import (
	"flag"
	"github.com/2hdddg/gop/client"
	//"github.com/2hdddg/gop/server"
	"log"
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

	// Configure standard logger
	log.SetFlags(log.Ltime | log.Lshortfile)

	if isServer {
		//server.Run(port)
		run_server(port)
		return
	}

	client.Run(port, &params)
}

/*
:lgetexpr system("~/code/go_parser/gop --func " . expand('<cword>') . " --file " . expand('%:p')) | lopen
*/
