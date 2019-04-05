package server

import (
	"github.com/2hdddg/gop/pkg/config"
	"github.com/2hdddg/gop/pkg/service/index"
	"github.com/2hdddg/gop/pkg/service/search"

	"log"
	"net"
	"net/http"
	"net/rpc"
	"strconv"
)

func Run(config *config.Config, port int) {
	searchSrv := search.NewService()
	err := searchSrv.Start()
	if err != nil {
		log.Fatalf("Failed to start search service: %s", err)
	}
	// Search service implements progress interface
	indexSrv := index.NewService(searchSrv)
	err = indexSrv.Start()
	if err != nil {
		log.Fatalf("Failed to start index service: %s", err)
	}

	for _, root := range config.Paths() {
		log.Println("Adding ", root)
		indexSrv.Add(root)
	}

	log.Printf("Starting server on port %d", port)
	rpc.HandleHTTP()
	l, err := net.Listen("tcp", ":"+strconv.Itoa(port))
	if err != nil {
		log.Fatalf("Failed to listen on tcp port %d: %s",
			port, err)
	}
	err = http.Serve(l, nil)
	if err != nil {
		log.Fatalf("Failed to start http server: %s", err)
	}
}
