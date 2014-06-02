package fracker

import (
	"github.com/coreos/go-etcd/etcd"
)

type Node interface {
	Each(func(string, string))
}

func NewNode(n *etcd.Node) Node {
	return &node{n}
}

type node struct {
	n *etcd.Node
}

func (self *node) Each(fn func(string, string)) {
	if self.n.Dir {
		for _, e_node := range self.n.Nodes {
			dir := &node{e_node}
			dir.Each(fn)
		}
	} else {
		fn(self.n.Key, self.n.Value)
	}
}
