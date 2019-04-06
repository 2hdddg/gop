package main

import (
	"flag"
	"fmt"
	"log"
	"net/rpc"
	"os"
	"strconv"

	"github.com/2hdddg/gop/pkg/config"
	"github.com/2hdddg/gop/pkg/parser"
	indexservice "github.com/2hdddg/gop/pkg/service/index"
	searchservice "github.com/2hdddg/gop/pkg/service/search"
)

func connectToServer(port int) *rpc.Client {
	// Panics when server not running
	defer func() {
		if r := recover(); r != nil {
			log.Fatalf("Failed to connect to server: %v\n", r)
		}
	}()

	client, err := rpc.DialHTTP("tcp", ":"+strconv.Itoa(port))
	if err != nil {
		log.Fatalf("Failed to dial server: %v\n", err)
	}

	return client
}

func search() {
	var (
		port int
		name string
		path string
	)

	flags := flag.NewFlagSet("search", flag.ExitOnError)
	flags.IntVar(&port, "port", 8080, "Server port")
	flags.StringVar(&name, "name", "", "Find definition")
	flags.StringVar(&path, "file", "", "Go file")
	flags.Parse(os.Args[2:])

	config := config.NewConfig()

	var imports []string

	client := connectToServer(port)

	// If a file is specified, parse it to extract context for
	// improved search precision.
	// List of imported packages is extracted to limit result
	// to those packages.
	if path != "" {
		imports, err := parser.ParseImports(path)
		if err != nil {
			log.Fatalf("Unable to parse imports from: %s", path)
		}
		curr, ok := config.PackageFromPath(path)
		if ok {
			imports = append(imports, curr)
		}
	}

	req := &searchservice.Request{
		Config:  *config,
		Name:    name,
		Imports: imports,
	}
	res, err := searchservice.Search(client, req)
	if err != nil {
		log.Fatalf("Failed to call server: %v", err)
	}

	// Output in grep format for vim to pickup
	for _, h := range res.Hits {
		fmt.Printf("%s:%d:%s\n", h.Path, h.Line, h.Descr)
	}
}

func index() {
	var (
		port int
		path string
	)

	flags := flag.NewFlagSet("index", flag.ExitOnError)
	flags.IntVar(&port, "port", 8080, "Server port")
	flags.StringVar(&path, "path", "", "Source code path")
	flags.Parse(os.Args[2:])

	if path == "" {
		fmt.Println("Needs path")
		return
	}
	client := connectToServer(port)
	req := &indexservice.IndexRequest{
		Path: path,
	}
	_, err := indexservice.Index(client, req)
	if err != nil {
		log.Fatalf("Failed to call server: %v", err)
	}
}