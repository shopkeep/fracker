package fracker

import (
	"github.com/coreos/go-etcd/etcd"
)

type Client interface {
	Get(key string) (Node, error)
}

func NewClient(hosts []string) Client {
	return &etcdClient{etcd.NewClient(hosts)}
}

type etcdClient struct {
	*etcd.Client
}

func (self *etcdClient) Get(key string) (Node, error) {
	var err error
	var resp *etcd.Response
	if resp, err = self.Client.Get(key, false, true); err != nil {
		return nil, err
	}
	return NewNode(resp.Node), nil
}
