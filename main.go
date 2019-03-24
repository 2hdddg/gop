package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/2hdddg/gop/client"
	"github.com/2hdddg/gop/config"
	"github.com/2hdddg/gop/server"
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
	flags.StringVar(&params.Name, "name", "", "Find function")
	flags.StringVar(&params.FilePath, "file", "", "Go file")
	flags.Parse(os.Args[2:])

	config := config.NewConfig()
	log.SetFlags(log.Ltime | log.Lshortfile)
	client.Run(config, port, &params)
}

func main() {
	command := os.Args[1]
	switch command {
	case "serve":
		serve()
	case "search":
		search()
	default:
		fmt.Println("Illegal command")
		return
	}
}

/*
:lgetexpr system("~/code/go_parser/gop --func " . expand('<cword>') . " --file " . expand('%:p')) | lopen
*/
