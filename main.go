package main

import (
	"flag"
	"fmt"
	"github.com/2hdddg/gop/server"
	"net/rpc"
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

	if isServer {
		server.Run(port)
		return
	}

	/* Client */
	if definition != "" {
		// When server is not running
		defer func() {
			if r := recover(); r != nil {
				fmt.Println("Failed to connect to server")
			}
		}()

		client, err := rpc.DialHTTP("tcp", ":1234")
		if err != nil {
			fmt.Println("Fatal", err)
		}
		a := &server.LocationsAnswer{}
		err = client.Call("Search.FuncDefinition", &definition, a)
		if err != nil {
			fmt.Println("Fatal", err)
		}
		fmt.Println(a)
	}
}

/*
:cgetexpr system("grep -n -r " . expand('<cword>') . " *")
*/
