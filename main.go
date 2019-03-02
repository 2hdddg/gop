package main

import (
	"flag"
	"github.com/2hdddg/gop/client"
	"github.com/2hdddg/gop/server"
	"log"
)

var (
	isServer bool
	port     int
	query    client.Query
)

func setupParameters() {
	flag.BoolVar(&isServer, "serve", false, "Run server")
	flag.IntVar(&port, "port", 8080, "Server port")

	flag.StringVar(&query.FuncFilter, "func", "", "Find function")
	flag.StringVar(&query.PackFilter, "pack", "", "Find package")
	flag.StringVar(&query.GoFilePath, "file", "", "Go file")

	flag.Parse()
}

func main() {
	setupParameters()

	// Configure standard logger
	log.SetFlags(log.Ltime | log.Lshortfile)

	if isServer {
		server.Run(port)
		return
	}

	client.Run(port, &query)
}

/*
:cgetexpr system("grep -n -r " . expand('<cword>') . " *")
:cgetexpr system("./gop --def " . expand('<cword>')) | copen
*/
