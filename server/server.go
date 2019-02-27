package server

import (
	"log"
	"net"
	"net/http"
	"net/rpc"
	"strconv"
	"time"
)

type packagesQuery struct {
	answerChan chan *PackagesAnswer
}

type locationsQuery struct {
	name       string
	answerChan chan *LocationsAnswer
}

type PackagesAnswer struct {
	Packages []string
}

type Location struct {
	Line   int
	Column int
}

type FileLocation struct {
	Location
	FilePath string
}

type LocationsAnswer struct {
	Locations []FileLocation
}

type build struct {
	fileChan  chan *file
	indexChan chan *index
}

func (b *build) thread() {
	count := 0
	builder := newBuilder()
	sendIndex := func() {
		log.Printf("Building and sending new index")
		b.indexChan <- builder.build()
		count = 0
	}

	for {
		select {
		case <-time.After(2 * time.Second):
			if count > 0 {
				sendIndex()
			}
		case file := <-b.fileChan:
			count += 1
			builder.add(file)
			if count > 100 {
				sendIndex()
			}
		}
	}
}

type search struct {
	indexChan     chan *index
	packQueryChan chan *packagesQuery
	locQueryChan  chan *locationsQuery
}

func (s *search) thread() {
	var currIndex *index
	for {
		select {
		case currIndex = <-s.indexChan:

		case q := <-s.packQueryChan:
			q.answerChan <- currIndex.allPackages()

		case q := <-s.locQueryChan:
			q.answerChan <- currIndex.funcDefinition(q.name)
		}
	}
}

func (s *search) FuncDefinition(
	name *string, a *LocationsAnswer) error {

	answerChan := make(chan *LocationsAnswer)
	query := &locationsQuery{name: *name, answerChan: answerChan}
	s.locQueryChan <- query
	*a = *<-answerChan

	return nil
}

func Run(port int) {
	indexChan := make(chan *index)

	search := search{
		locQueryChan:  make(chan *locationsQuery),
		packQueryChan: make(chan *packagesQuery),
		indexChan:     indexChan}
	go search.thread()

	build := build{
		fileChan:  make(chan *file),
		indexChan: indexChan}
	go build.thread()

	config := newConfig()
	log.Printf("Server config: %+v", config)
	if !config.valid() {
		log.Fatalln("Invalid config")
	}

	// In go thread to be responsive to client requests, could
	// of course cause in-complete answers but better than hanging
	// or failing.
	go config.system.build(build.fileChan)

	log.Printf("Starting server on port %d", port)
	err := rpc.RegisterName("Search", &search)
	if err != nil {
		log.Fatalf("Failed to register search service: %s", err)
	}
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
