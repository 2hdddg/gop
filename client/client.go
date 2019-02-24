package client

import (
	"fmt"
	"github.com/2hdddg/gop/server"
	"log"
	"net/rpc"
	"strconv"
)

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

func Run(port int, def string) {
	client, err := connectToServer(port)
	if err != nil {
		log.Fatalf("Failed to connect to server: %s", err)
	}

	if def != "" {
		a := &server.LocationsAnswer{}
		err = client.Call("Search.FuncDefinition", &def, a)
		if err != nil {
			log.Fatalf("Failed to call server: %s", err)
		}

		// Write to stdout in grep format
		for _, l := range a.Locations {
			fmt.Printf("%s:%d:Definition of %s\n",
				l.FilePath, l.Line, def)
		}
	}
}
