package server

import (
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"net/rpc"
	"path"
	"strconv"
	"sync"
)

func probe(path string) (dirs, files []string) {
	entries, err := ioutil.ReadDir(path)
	if err != nil {
		log.Printf("Error indexing package at %v: %v\n", path, err)
		return nil, nil
	}

	for _, i := range entries {
		mode := i.Mode()
		if mode.IsDir() {
			dirs = append(dirs, i.Name())
		} else if mode.IsRegular() {
			files = append(files, i.Name())
		}
	}
	return dirs, files
}

func (b *build) parseFiles(p string, files []string) {
	var waitGroup sync.WaitGroup

	for _, tmp := range files {
		waitGroup.Add(1)
		f := tmp
		go func() {
			defer waitGroup.Done()
			parsed := parseFile(path.Join(p, f))
			b.fileChan <- *parsed
		}()
	}
	waitGroup.Wait()
}

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
	fileChan  chan file
	indexChan chan index
}

func (b *build) thread() {
	builder := newBuilder()
	for {
		file := <-b.fileChan
		builder.add(&file)
		b.indexChan <- builder.build()
	}
}

type search struct {
	indexChan     chan index
	packQueryChan chan packagesQuery
	locQueryChan  chan locationsQuery
}

func (s *search) thread() {
	var currIndex index
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
	query := locationsQuery{name: *name, answerChan: answerChan}
	s.locQueryChan <- query
	*a = *<-answerChan

	return nil
}

func Run(port int) {
	indexChan := make(chan index)

	search := search{
		locQueryChan:  make(chan locationsQuery),
		packQueryChan: make(chan packagesQuery),
		indexChan:     indexChan}
	go search.thread()

	build := build{
		fileChan:  make(chan file),
		indexChan: indexChan}
	go build.thread()

	path := "/usr/share/go/src/io"
	_, files := probe(path)
	build.parseFiles(path, files)

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
