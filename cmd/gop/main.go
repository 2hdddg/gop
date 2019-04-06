package main

import (
	"fmt"
	"log"
	"os"
)

type command struct {
	name    string
	descr   string
	handler func()
}

var commands = []command{
	command{name: "serve", descr: "Start server", handler: serve},
	command{name: "search", descr: "Search", handler: search},
	command{name: "index", descr: "Index path", handler: index},
}

func showRootUsage() {
	fmt.Println("Usage:")
	for _, c := range commands {
		fmt.Printf("%v %v\n", c.name, c.descr)
	}
}

func main() {
	log.SetFlags(log.Ltime | log.Lshortfile)
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
