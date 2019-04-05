package index

import (
	"net/rpc"
)

// Exposed over RPC
type IndexRequest struct {
	Path string
}

// Exposed over RPC
type IndexResponse struct {
}

type RpcInstance struct {
	indexPathChan chan *indexPathMsg
}

func registerRpcInstance(ch chan *indexPathMsg) error {
	i := &RpcInstance{
		indexPathChan: ch,
	}
	err := rpc.RegisterName("Index", i)
	if err != nil {
		return err
	}
	return nil
}

func (i *RpcInstance) Add(req *IndexRequest, res *IndexResponse) error {
	i.indexPathChan <- &indexPathMsg{path: req.Path}
	return nil
}
