package server

import (
	"github.com/2hdddg/gop/shared"
	"log"
	"time"
)

func monitor(config *shared.Config, indexChan chan *index) {
	count := 0
	systemTree := newTree(config.SystemPath)
	systemFileChan := make(chan *file)

	go newTraverser(config.SystemPath, systemFileChan).traverse()

	sendIndex := func() {
		log.Printf("Building and sending new index")
		indexChan <- systemTree.buildIndex()
		count = 0
	}

	for {
		select {
		// Timeout
		case <-time.After(2 * time.Second):
			if count > 0 {
				sendIndex()
			}
		// Progress/update on system tree
		case file := <-systemFileChan:
			count += 1
			systemTree.addFile(file)
			if count > 100 {
				sendIndex()
			}
		}
	}
}
