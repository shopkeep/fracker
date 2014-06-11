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
	*etcd.Node
}

func (self *node) Each(fn func(string, string)) {
	if self.Dir {
		for _, child := range self.Nodes {
			n := &node{child}
			n.Each(fn)
		}
	} else {
		fn(self.Key, self.Value)
	}
}
