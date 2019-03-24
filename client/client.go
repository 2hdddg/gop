package client

import (
	"fmt"
	"log"
	"net/rpc"
	"strconv"

	"github.com/2hdddg/gop/config"
	"github.com/2hdddg/gop/parser"
	"github.com/2hdddg/gop/service/search"
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

func Run(config *config.Config, port int, params *Params) {
	var imports []string

	client, err := connectToServer(port)
	if err != nil {
		log.Fatalf("Failed to connect to server: %s", err)
	}

	// If a file is specified, parse it to extract context for
	// improved search precision.
	// List of imported packages is extracted to limit result
	// to those packages.
	if params.FilePath != "" {
		imports, err = parser.ParseImports(params.FilePath)
		if err != nil {
			log.Fatalf("Unable to parse imports from: %s",
				params.FilePath)
		}
		curr, ok := config.PackageFromPath(params.FilePath)
		if ok {
			imports = append(imports, curr)
		}
	}

	req := &search.Request{
		Config:  *config,
		Name:    params.Name,
		Imports: imports,
	}
	res, err := search.Search(client, req)
	if err != nil {
		log.Fatalf("Failed to call server: %v", err)
	}

	// Output in grep format for vim to pickup
	for _, h := range res.Hits {
		fmt.Printf("%s:%d:%s\n", h.Path, h.Line, h.Descr)
	}
}
