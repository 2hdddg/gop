package server

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"path"
	"sync"
)

type QueryRequest struct {
	answerChan chan Answer
	query      Query
}

/*
:cgetexpr system("grep -n -r " . expand('<cword>') . " *")
*/

func probe(path string) (dirs, files []string) {
	entries, err := ioutil.ReadDir(path)
	if err != nil {
		fmt.Printf("Error indexing package at %v: %v\n", path, err)
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

func parseFiles(p string, files []string, fileChan chan File) {
	var waitGroup sync.WaitGroup

	for _, tmp := range files {
		waitGroup.Add(1)
		f := tmp
		go func() {
			defer waitGroup.Done()
			parsed := parseFile(path.Join(p, f))
			fileChan <- *parsed
		}()
	}
	waitGroup.Wait()
}

func builder(fileChan chan File, indexChan chan Index) {
	builder := NewBuilder()
	for {
		file := <-fileChan
		builder.Add(&file)
		indexChan <- builder.Build()
	}
}

func searcher(indexChan chan Index, queryReqChan chan QueryRequest) {
	var index Index

	for {
		select {
		case index = <-indexChan:
			fmt.Println("Received new search index")
		case queryReq := <-queryReqChan:
			fmt.Printf("Received query: %T\n", queryReq.query)
			queryReq.answerChan <- queryReq.query.Process(&index)
		}
	}
}

func queryRequest(queryReqChan chan QueryRequest, query Query) Answer {
	answerChan := make(chan Answer)
	queryReq := QueryRequest{answerChan: answerChan, query: query}
	queryReqChan <- queryReq
	answer := <-answerChan
	return answer
}

func Run(port int) {
	fileChan := make(chan File)
	indexChan := make(chan Index)
	queryReqChan := make(chan QueryRequest)

	// Builds index of received files and sends index on channel.
	go builder(fileChan, indexChan)
	// Receives index used to process queries
	go searcher(indexChan, queryReqChan)

	path := "/usr/share/go/src/io"
	dirs, files := probe(path)
	parseFiles(path, files, fileChan)
	fmt.Println(dirs)

	http.HandleFunc("/packages",
		func(w http.ResponseWriter, r *http.Request) {
			query := &PackagesQuery{}
			answer := queryRequest(queryReqChan, query)
			answer.Response(w)
		})
	http.HandleFunc("/definition",
		func(w http.ResponseWriter, r *http.Request) {
			params := r.URL.Query()
			query := &DefinitionQuery{name: params["name"][0]}
			answer := queryRequest(queryReqChan, query)
			answer.Response(w)
		})

	http.ListenAndServe(":8080", nil)
}
