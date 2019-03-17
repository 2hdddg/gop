package client

import (
	"fmt"
	"log"
	"net/rpc"
	"strconv"

	"github.com/2hdddg/gop/config"
	"github.com/2hdddg/gop/search"
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

func invoke(client *rpc.Client, req *search.Request) {
	res := &search.Response{}
	err := client.Call("Search.Search", req, res)
	if err != nil {
		log.Fatalf("Failed to call server: %s", err)
	}

	for _, h := range res.Hits {
		fmt.Printf("%s:%d:%s\n", h.Path, h.Line, h.Descr)
	}
}

func Run(config *config.Config, port int, params *Params) {

	client, err := connectToServer(port)
	if err != nil {
		log.Fatalf("Failed to connect to server: %s", err)
	}

	req := &search.Request{
		//Config: config,
		Name: params.Name,
	}

	invoke(client, req)
}
