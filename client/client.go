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
	GoFilePath string
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

func invoke(client *rpc.Client, filter string, object server.Object) {
	answer := &server.Answer{}
	query := &server.Query{Object: object, Name: filter}
	err := client.Call("Search.Search", query, answer)
	if err != nil {
		log.Fatalf("Failed to call server: %s", err)
	}

	// Write to stdout in grep format
	for _, l := range answer.Locations {
		writeInGrepFormat(l.Path, "object", filter, l.Line)
	}
}

func Run(port int, query *Query) {
	client, err := connectToServer(port)
	if err != nil {
		log.Fatalf("Failed to connect to server: %s", err)
	}

	if query.FuncFilter != "" {
		invoke(client, query.FuncFilter, server.Function)
	}

	if query.PackFilter != "" {
		invoke(client, query.PackFilter, server.Package)
	}
}
