package server

import (
	"github.com/2hdddg/gop/shared"
	"log"
	"net"
	"net/http"
	"net/rpc"
	"strconv"
)

type Object int

const (
	Function Object = 0
	Package  Object = 1
)

type Query struct {
	Object   Object
	Name     string
	Packages []string
	Config   *shared.Config
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
				q.answerChan <- currIndex.funcByQuery(q.Query)
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

	config := shared.NewConfig()
	log.Printf("Server config: %+v", config)
	if !config.Valid() {
		log.Fatalln("Invalid config")
	}

	go monitor(config, indexChan)

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
