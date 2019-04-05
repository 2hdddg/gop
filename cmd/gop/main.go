package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/2hdddg/gop/cmd/gop/client"
	"github.com/2hdddg/gop/cmd/gop/server"
	"github.com/2hdddg/gop/pkg/config"
)

func serve() {
	var port int

	flags := flag.NewFlagSet("serve", flag.ExitOnError)
	flags.IntVar(&port, "port", 8080, "Server port")
	flags.Parse(os.Args[2:])

	config := config.NewConfig()
	log.SetFlags(log.Ltime | log.Lshortfile)
	server.Run(config, port)
}

func search() {
	var (
		port   int
		params client.Params
	)

	flags := flag.NewFlagSet("search", flag.ExitOnError)
	flags.IntVar(&port, "port", 8080, "Server port")
	flags.StringVar(&params.Name, "name", "", "Find definition")
	flags.StringVar(&params.FilePath, "file", "", "Go file")
	flags.Parse(os.Args[2:])

	config := config.NewConfig()
	log.SetFlags(log.Ltime | log.Lshortfile)
	client.Run(config, port, &params)
}

type command struct {
	name    string
	descr   string
	handler func()
}

var commands = []command{
	command{name: "serve", descr: "Start server", handler: serve},
	command{name: "search", descr: "Search", handler: search},
}

func showRootUsage() {
	fmt.Println("Usage:")
	for _, c := range commands {
		fmt.Printf("%v %v\n", c.name, c.descr)
	}
}

func main() {
	if len(os.Args) <= 1 {
		showRootUsage()
		return
	}

	name := os.Args[1]
	for _, c := range commands {
		if c.name == name {
			c.handler()
			return
		}
	}

	fmt.Println("Illegal command")
	showRootUsage()
}

/*
:lgetexpr system("~/code/go_parser/gop --func " . expand('<cword>') . " --file " . expand('%:p')) | lopen
*/
