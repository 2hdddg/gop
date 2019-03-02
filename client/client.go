package client

import (
	"fmt"
	"github.com/2hdddg/gop/server"
	"log"
	"net/rpc"
	"strconv"
)

type Query struct {
	FuncFilter string
	PackFilter string
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

func writeInGrepFormat(path, what, filter string, line int) {
	fmt.Printf("%s:%d:%s matching '%s'\n", path, line, what, filter)
}

func Run(port int, query *Query) {
	client, err := connectToServer(port)
	if err != nil {
		log.Fatalf("Failed to connect to server: %s", err)
	}

	if query.FuncFilter != "" {
		a := &server.Answer{}
		err = client.Call("Search.Func", &query.FuncFilter, a)
		if err != nil {
			log.Fatalf("Failed to call server: %s", err)
		}

		// Write to stdout in grep format
		for _, l := range a.Locations {
			writeInGrepFormat(l.Path, "Func", query.FuncFilter, l.Line)
		}
	}

	if query.PackFilter != "" {
		a := &server.Answer{}
		err = client.Call("Search.Pack", &query.PackFilter, a)
		if err != nil {
			log.Fatalf("Failed to call server: %s", err)
		}

		// Write to stdout in grep format
		for _, l := range a.Locations {
			writeInGrepFormat(l.Path, "Pack", query.PackFilter, l.Line)
		}
	}
}
