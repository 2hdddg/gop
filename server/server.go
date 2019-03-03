package server

import (
	"log"
	"net"
	"net/http"
	"net/rpc"
	"strconv"
	"time"
)

type Object int

const (
	Function Object = 0
	Package  Object = 1
)

type Query struct {
	Object        Object
	Name          string
	Packages      []string
	WorkspacePath string
}

// Internal query
type query struct {
	*Query
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
	packs := newPacks()
	sendIndex := func() {
		log.Printf("Building and sending new index")
		b.indexChan <- packs.buildIndex()
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
			packs.addFile(file)
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
			switch q.Object {
			case Function:
				q.answerChan <- currIndex.funcByName(q.Name)
			case Package:
				q.answerChan <- currIndex.packByName(q.Name)
			}
		}
	}
}

func (s *search) Search(clientQuery *Query, answer *Answer) error {
	q := &query{clientQuery, make(chan *Answer)}
	s.queryChan <- q
	*answer = *<-q.answerChan
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

	monitor := newMonitor(config, build.fileChan)
	monitor.start()

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
