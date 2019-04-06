package service

import (
	"net/rpc"

	"fmt"

	"github.com/2hdddg/gop/pkg/config"
	"github.com/2hdddg/gop/pkg/index"
	"github.com/2hdddg/gop/pkg/tree"
)

type SearchReq struct {
	Name    string
	Imports []string
	Config  config.Config
}

type SearchRes struct {
	Hits []Hit
}

type IndexReq struct {
	Path string
}

type IndexRes struct {
}

type Hit struct {
	Path  string
	Line  int
	Descr string
}

type Service struct {
	treeChan   chan *treeMsg
	searchChan chan *searchMsg
	indexChan  chan *indexMsg
}

type RpcService struct {
	service *Service
}

type ackMsg struct {
	err error
}

type searchMsg struct {
	ackChan   chan ackMsg
	clientReq *SearchReq
	clientRes *SearchRes
}

type treeMsg struct {
	ackChan chan ackMsg
	tree    *tree.Tree
}

type indexMsg struct {
	ackChan   chan ackMsg
	clientReq *IndexReq
	clientRes *IndexRes
}

func NewService() *Service {
	return &Service{
		treeChan:   make(chan *treeMsg),
		searchChan: make(chan *searchMsg),
		indexChan:  make(chan *indexMsg),
	}
}

func search(msg *searchMsg, indexmap map[string]*index.Index) {
	req := msg.clientReq
	res := msg.clientRes

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
			fmt.Printf("No index for %v\n", r)
		}
	}

	// From this point it's ok to return and continue in go routine

	q := index.NewQuery(req.Name)
	q.Imported = req.Imports
	for _, i := range indexes {
		i.Query(q,
			func(h index.Hit) {
				res.Hits = append(res.Hits, Hit{
					Path:  h.Path(),
					Line:  h.Symbol.Line,
					Descr: h.Symbol.ToString(),
				})
			},
			func(p index.Package) {
				res.Hits = append(res.Hits, Hit{
					Path:  p.Path,
					Line:  0,
					Descr: " Package",
				})
			})
	}
	fmt.Printf("Search for '%v' resulted in %v hits from %v indexes\n",
		req.Name, len(res.Hits), len(indexes))
}

func (s *Service) service() {
	indexes := map[string]*index.Index{}
	indexChan := make(chan *index.Index)

	for {
		select {
		// Start building a tree
		case i := <-s.indexChan:
			go func() {
				builder, err := tree.NewBuilder(i.clientReq.Path)
				if err != nil {
					i.ackChan <- ackMsg{err: err}
					return
				}
				// Do the rest async from client
				i.ackChan <- ackMsg{}

				fmt.Printf("Analyzing %v\n", i.clientReq.Path)
				builder.Progress = s
				_, err = builder.Build()
				if err != nil {
					fmt.Printf("Failed to analyze %v: %s\n",
						i.clientReq.Path, err)
					return
				}
			}()
		// New tree received, build index
		case m := <-s.treeChan:
			go func() {
				indexChan <- index.Build(m.tree)
				m.ackChan <- ackMsg{}
			}()
		// Received index built in go routine above.
		case m := <-indexChan:
			fmt.Printf("All or parts of %v is searchable\n", m.RootPath)
			indexes[m.RootPath] = m
		// Perform search
		case m := <-s.searchChan:
			search(m, indexes)
			m.ackChan <- ackMsg{}
		}
	}
}

func (s *Service) Start() error {
	r := &RpcService{s}
	if err := rpc.RegisterName("Search", r); err != nil {
		return err
	}
	go s.service()
	return nil
}

func (s *Service) newOrUpdatedTree(t *tree.Tree) {
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
	s.newOrUpdatedTree(t)
}

// RPC function
func (r *RpcService) Search(req *SearchReq, res *SearchRes) error {
	ackChan := make(chan ackMsg)
	r.service.searchChan <- &searchMsg{
		clientReq: req,
		clientRes: res,
		ackChan:   ackChan,
	}
	<-ackChan

	return nil
}

// Called from RPC client. Wrapper.
func Search(c *rpc.Client, req *SearchReq) (*SearchRes, error) {
	res := &SearchRes{}
	if err := c.Call("Search.Search", req, res); err != nil {
		return nil, err
	}
	return res, nil
}

func (s *Service) Index(req *IndexReq, res *IndexRes) error {
	ackChan := make(chan ackMsg)
	s.indexChan <- &indexMsg{
		clientReq: req,
		clientRes: res,
		ackChan:   ackChan,
	}
	a := <-ackChan

	return a.err
}

// RPC function
func (r *RpcService) Index(req *IndexReq, res *IndexRes) error {
	return r.service.Index(req, res)
}

// Called from RPC client. Wrapper.
func Index(c *rpc.Client, req *IndexReq) (*IndexRes, error) {
	res := &IndexRes{}
	if err := c.Call("Search.Index", req, res); err != nil {
		return nil, err
	}
	return res, nil
}
