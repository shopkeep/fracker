package fracker

import (
	"github.com/coreos/go-etcd/etcd"
)

type Node interface {
	Each(func(string, string))
}

func NewNode(n *etcd.Node) Node {
	if n.Dir {
		return &dirNode{n}
	} else {
		return &fileNode{n}
	}
}

type node struct {
	*etcd.Node
}

type dirNode node

func (self *dirNode) Each(fn func(string, string)) {
	for _, n := range self.Node.Nodes {
		NewNode(n).Each(fn)
	}
}

type fileNode node

func (self *fileNode) Each(fn func(string, string)) {
	fn(self.Node.Key, self.Node.Value)
}
