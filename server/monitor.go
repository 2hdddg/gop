package server

import (
	"github.com/2hdddg/gop/shared"
	"log"
	"time"
)

func monitor(config *shared.Config, indexChan chan *index) {
	systemTree := newTree(config.SystemPath)
	systemFileChan := make(chan *file)
	workspaceFileChan := make(chan *file)
	var workspaceTree *tree

	go newTraverser(config.SystemPath, systemFileChan).traverse()
	if config.WorkspacePath != "" {
		workspaceTree = newTree(config.WorkspacePath)
		go newTraverser(
			config.WorkspacePath,
			workspaceFileChan).traverse()
	}

	sendIndex := func(t *tree) {
		log.Printf("Building and sending new index")
		indexChan <- t.buildIndex()
	}

	for {
		select {
		// Timeout
		case <-time.After(2 * time.Second):
			if systemTree.dirty {
				sendIndex(systemTree)
			}
			if workspaceTree != nil && workspaceTree.dirty {
				sendIndex(workspaceTree)
			}
		// Progress/update on system tree
		case file := <-systemFileChan:
			systemTree.addFile(file)
		case file := <-workspaceFileChan:
			workspaceTree.addFile(file)
		}
	}
}
