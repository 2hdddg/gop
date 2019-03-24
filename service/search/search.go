package search

// Search service

import (
	"log"
	"net/rpc"

	"github.com/2hdddg/gop/config"
	"github.com/2hdddg/gop/index"
	"github.com/2hdddg/gop/tree"
)

type Service struct {
	treeChan chan *treeMsg
	reqChan  chan *requestMsg
}

// RPC instance
type Instance struct {
	reqChan chan *requestMsg
}

// Exposed over RPC
type Request struct {
	Name    string
	Imports []string
	Config  config.Config
}

// Exposed over RPC
type Response struct {
	Hits []Hit
}

// Exposed over RPC
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

func (res *Response) add(hits []*index.Hit) {
	for _, h := range hits {
		res.Hits = append(res.Hits, Hit{
			Path:  h.Path(),
			Line:  h.Line,
			Descr: h.Extra,
		})
	}
}

func search(req *Request, res *Response, indexmap map[string]*index.Index) {
	// Copy the needed indexes from map to list to make sure that
	// we can return to search fast and continue the search in a
	// go routine.
	roots := req.Config.Paths()
	indexes := make([]*index.Index, 0, len(roots))
	for _, r := range roots {
		i, exists := indexmap[r]
		if exists {
			indexes = append(indexes, i)
		} else {
			log.Printf("No index for %v", r)
		}
	}

	// From this point it's ok to return and continue in go routine

	q := index.NewQuery(req.Name)
	q.Imported = req.Imports
	for _, i := range indexes {
		result := i.Query(q)
		res.add(result)
	}
}

func (s *Service) service() {
	indexes := map[string]*index.Index{}
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
			indexes[m.RootPath] = m
		case m := <-s.reqChan:
			search(m.clientReq, m.clientRes, indexes)
			m.ackChan <- ackMsg{}
		}
	}
}

func (s *Service) Start() error {
	i := &Instance{reqChan: s.reqChan}
	err := rpc.RegisterName("Search", i)
	if err != nil {
		return err
	}
	go s.service()
	return nil
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

func (i *Instance) Search(req *Request, res *Response) error {
	ackChan := make(chan ackMsg)
	i.reqChan <- &requestMsg{
		clientReq: req,
		clientRes: res,
		ackChan:   ackChan,
	}
	<-ackChan

	return nil
}

// Called from RPC client. Wrapper.
func Search(c *rpc.Client, req *Request) (*Response, error) {
	res := &Response{}
	if err := c.Call("Search.Search", req, res); err != nil {
		return nil, err
	}
	return res, nil
}
