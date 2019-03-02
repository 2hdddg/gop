package server

import (
	"log"
	"net"
	"net/http"
	"net/rpc"
	"strconv"
	"time"
)

type object int

const (
	Function object = 0
	Package  object = 1
)

type query struct {
	object     object
	name       string
	answerChan chan *Answer
}

type Location struct {
	Path   string
	Line   int
	Column int
}

type Answer struct {
	Locations []Location
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
	indexChan chan *index
	queryChan chan *query
}

func (s *search) thread() {
	var currIndex *index
	for {
		select {
		case currIndex = <-s.indexChan:

		case q := <-s.queryChan:
			switch q.object {
			case Function:
				q.answerChan <- currIndex.funcByName(q.name)
			case Package:
				q.answerChan <- currIndex.packByName(q.name)
			}
		}
	}
}

func (s *search) Func(name *string, a *Answer) error {

	answerChan := make(chan *Answer)
	query := &query{
		object: Function, name: *name, answerChan: answerChan}
	s.queryChan <- query
	*a = *<-answerChan

	return nil
}

func (s *search) Pack(name *string, a *Answer) error {

	answerChan := make(chan *Answer)
	query := &query{
		object: Package, name: *name, answerChan: answerChan}
	s.queryChan <- query
	*a = *<-answerChan

	return nil
}

func Run(port int) {
	indexChan := make(chan *index)

	search := search{
		queryChan: make(chan *query),
		indexChan: indexChan}
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
