package index

// Index service

import (
	"log"

	"github.com/2hdddg/gop/pkg/tree"
)

type Service struct {
	progress      tree.Progress
	indexPathChan chan *indexPathMsg
}

func NewService(progress tree.Progress) *Service {
	return &Service{
		progress:      progress,
		indexPathChan: make(chan *indexPathMsg),
	}
}

type indexPathMsg struct {
	path string
}

func (s *Service) service() {
	for {
		select {
		case i := <-s.indexPathChan:
			go func() {
				builder, err := tree.NewBuilder(i.path)
				if err != nil {
					log.Printf("Failed to init tree %v\n", err)
					return
				}
				builder.Progress = s.progress
				_, err = builder.Build()
				if err != nil {
					log.Printf("Failed to build tree: %s\n", err)
					return
				}
			}()
		}
	}
}

func (s *Service) Start() error {
	if err := registerRpcInstance(s.indexPathChan); err != nil {
		return err
	}
	go s.service()
	return nil
}

func (s *Service) Add(path string) {
	s.indexPathChan <- &indexPathMsg{path: path}
}
