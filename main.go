package main

import (
	"flag"
	"fmt"
	"github.com/2hdddg/gop/server"
	"net/rpc"
	"strconv"
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

func connectToServer(port int) (client *rpc.Client, err error) {
	// Panics when server not running
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("Failed to connect to server: %s", r)
		}
	}()

	client, err = rpc.DialHTTP("tcp", ":"+strconv.Itoa(port))

	return client, err
}

func main() {
	setupParameters()

	if isServer {
		server.Run(port)
		return
	}

	client, err := connectToServer(port)
	if err != nil {
		fmt.Println("Fatal", err)
		return
	}

	if definition != "" {
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
