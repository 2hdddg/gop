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

	definition string
)

func setupParameters() {
	flag.BoolVar(&isServer, "serve", false, "Run server")
	flag.IntVar(&port, "port", 8080, "Server port")

	flag.StringVar(&definition, "def", "", "Find definition")
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

	client.Run(port, definition)
}

/*
:cgetexpr system("grep -n -r " . expand('<cword>') . " *")
:cgetexpr system("./gop --def " . expand('<cword>')) | copen
*/
