package server

import (
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"net/rpc"
	"path"
	"strconv"
	"strings"
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
		name := i.Name()
		if mode.IsDir() {
			dirs = append(dirs, name)
		} else if mode.IsRegular() && strings.LastIndex(name, ".go") > 0 {
			files = append(files, name)
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
			parsed, err := parseFile(path.Join(p, f))
			if err == nil {
				b.fileChan <- parsed
			}
		}()
	}
	waitGroup.Wait()
}

func (b *build) addToBuilder(p string) {
	dirs, files := probe(p)

	log.Printf("Analyzing %s", p)
	b.parseFiles(p, files)

	for _, dir := range dirs {
		b.addToBuilder(path.Join(p, dir))
	}
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
	fileChan  chan *file
	indexChan chan *index
}

func (b *build) thread() {
	builder := newBuilder()
	for {
		file := <-b.fileChan
		builder.add(file)
		b.indexChan <- builder.build()
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

	/*
		path := "/usr/share/go/src/io"
		_, files := probe(path)
		build.parseFiles(path, files)
	*/
	build.addToBuilder(config.goRoot)

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
