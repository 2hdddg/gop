package client

import (
	"fmt"
	"github.com/2hdddg/gop/server"
	"github.com/2hdddg/gop/shared"
	"log"
	"net/rpc"
	"strconv"
)

type Params struct {
	Name     string
	FilePath string
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

func invoke(client *rpc.Client, query *server.Query) {
	answer := &server.Answer{}
	err := client.Call("Search.Search", query, answer)
	if err != nil {
		log.Fatalf("Failed to call server: %s", err)
	}

	// Write to stdout in grep format
	for _, l := range answer.Locations {
		writeInGrepFormat(l.Path, "object", query.Name, l.Line)
	}
	for _, e := range answer.Errors {
		fmt.Println(e)
	}
}

func Run(port int, params *Params) {
	config := shared.NewConfig()

	client, err := connectToServer(port)
	if err != nil {
		log.Fatalf("Failed to connect to server: %s", err)
	}

	query := &server.Query{
		Config: config,
		Name:   params.Name,
	}

	if params.FilePath != "" {
		query.Packages, _ = parseFileImports(params.FilePath)
		filePackName, _ := parseFilePackage(config, params.FilePath)
		query.Packages = append(query.Packages, filePackName)
	}
	query.Object = server.Function
	invoke(client, query)

	/*
		if params.PackFilter != "" {
			query.Object = server.Package
			invoke(client, query)
		}
	*/
}
