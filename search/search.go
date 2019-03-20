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

func (res *Response) add(hits []*index.Hit, descr string) {
	for _, h := range hits {
		res.Hits = append(res.Hits, Hit{
			Path:  h.Path(),
			Line:  h.Line,
			Descr: descr,
		})
	}
}

func (s *Service) service() {
	var i *index.Index

	indexChan := make(chan *index.Index)

	log.Println("Started search service")

	for {
		select {
		case m := <-s.treeChan:
			log.Println("Got new/updated tree")
			go func() {
				indexChan <- index.Build(m.tree)
				m.ackChan <- ackMsg{}
			}()
		case m := <-indexChan:
			// Received index built in go routine above.
			i = m
		case m := <-s.reqChan:
			// Serve search request
			if i != nil {
				// Use pointer to curr index in case of index is rebuilt
				// while go func is running.
				go func(ii *index.Index) {
					res := ii.Query(&index.Query{Name: m.clientReq.Name})
					m.clientRes.add(res.Functions, "Function")
					m.clientRes.add(res.Methods, "Method")
					m.ackChan <- ackMsg{}
				}(i)
			} else {
				// No index, yet?
				m.ackChan <- ackMsg{}
			}
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

// Implements tree Progress interface
func (s *Service) OnPackageParsed(t *tree.Tree, p *tree.Package) {
}

// Implements tree Progress interface
func (s *Service) OnTreeParsed(t *tree.Tree) {
	s.NewOrUpdatedTree(t)
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