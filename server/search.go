package server

import (
	"log"
)

type search struct {
	indexChan chan *index
	queryChan chan *query
	indexes   map[string]*index // Key: root
}

func newSearch() *search {
	return &search{
		indexChan: make(chan *index),
		queryChan: make(chan *query),
		indexes:   make(map[string]*index),
	}
}

type query struct {
	*Query
	answerChan chan *Answer
}

func (s *search) thread() {
	for {
		select {
		case index := <-s.indexChan:
			s.indexes[index.root] = index
			log.Printf("Updated index for %v", index.root)

		case q := <-s.queryChan:
			a := &Answer{}
			f := func(i *index) {
				switch q.Object {
				case Function:
					i.funcByQuery(q.Query, a)
				case Package:
					i.packByName(q.Name, a)
				}
			}

			index, exists := s.indexes[q.Config.SystemPath]
			if exists {
				f(index)
			} else {
				a.Errors = append(a.Errors, "Syspath mismatch")
			}
			index, exists = s.indexes[q.Config.WorkspacePath]
			if exists {
				f(index)
			} else {
				a.Errors = append(a.Errors, "Workpath mismatch")
			}
			q.answerChan <- a
		}
	}
}

func (s *search) Search(clientQuery *Query, answer *Answer) error {
	q := &query{clientQuery, make(chan *Answer)}
	s.queryChan <- q
	*answer = *<-q.answerChan
	return nil
}
