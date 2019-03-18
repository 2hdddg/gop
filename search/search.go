package search

// Search service

import (
	"log"

	"github.com/2hdddg/gop/index"
	"github.com/2hdddg/gop/tree"
)

type Service struct {
	treeChan chan *treeMsg
	reqChan  chan *requestMsg
}

// RPC client
type Client struct {
	reqChan chan *requestMsg
}

// Exposed over RPC
type Request struct {
	Name string
}

// Exposed over RPC
type Response struct {
	Hits []Hit
}

type Hit struct {
	Path  string
	Line  int
	Descr string
}

type ackMsg struct{}

type requestMsg struct {
	ackChan   chan ackMsg
	clientReq *Request
	clientRes *Response
}

type treeMsg struct {
	ackChan chan ackMsg
	tree    *tree.Tree
}

func NewService() *Service {
	return &Service{
		treeChan: make(chan *treeMsg),
		reqChan:  make(chan *requestMsg),
	}
}

func (s *Service) service() {
	var i index.Index

	log.Println("Started search service")

	for {
		select {
		case m := <-s.treeChan:
			log.Println("Got new/updated tree")
			i = index.Build(m.tree)
			m.ackChan <- ackMsg{}
		case m := <-s.reqChan:
			res := i.Query(&index.Query{Name: m.clientReq.Name})

			for _, h := range res.Functions {
				m.clientRes.Hits = append(m.clientRes.Hits, Hit{
					Path:  h.Path(),
					Line:  h.Line,
					Descr: "Func def",
				})
			}

			m.ackChan <- ackMsg{}
		}
	}
}

func (s *Service) Start() *Client {
	go s.service()
	return &Client{reqChan: s.reqChan}
}

func (s *Service) NewOrUpdatedTree(t *tree.Tree) {
	ackChan := make(chan ackMsg)
	s.treeChan <- &treeMsg{
		tree:    t,
		ackChan: ackChan,
	}
	<-ackChan
}

func (c *Client) Search(req *Request, res *Response) error {
	ackChan := make(chan ackMsg)
	c.reqChan <- &requestMsg{
		clientReq: req,
		clientRes: res,
		ackChan:   ackChan,
	}
	<-ackChan

	return nil
}
