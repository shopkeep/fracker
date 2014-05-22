package main

import (
	"github.com/coreos/go-etcd/etcd"
)

type EtcdClient interface {
	Get(key string) (Node, error)
}

func NewClient(hosts []string) EtcdClient {
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
