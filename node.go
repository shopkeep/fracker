package fracker

import (
	"github.com/coreos/go-etcd/etcd"
	"strings"
)

type KV interface {
	Key() string
	Value() string
}

type Node interface {
	KV
	Each(func(string, string))
	Parent() Node
	Node() *etcd.Node
}

func NewNode(n *etcd.Node) Node {
	if n.Dir {
		return &dirNodeClean{&node{n}}
	} else {
		return &fileNodeClean{&node{n}}
	}
}

type node struct {
	n *etcd.Node
}

func (self *node) Node() *etcd.Node {
	return self.n
}

func (self *node) Key() string {
	return self.n.Key
}

func (self *node) Value() string {
	return self.n.Value
}

type dirNodeClean struct {
	*node
}

func (self *dirNodeClean) Node() *etcd.Node {
	return self.node.Node()
}

func (self *dirNodeClean) Parent() Node {
	return self
}

func (self *dirNodeClean) Each(fn func(string, string)) {
	for _, n := range self.Node().Nodes {
		if n.Dir {
			dir := &dirNodeFromDir{&node{n}, self}
			dir.Each(fn)
		} else {
			file := &fileNodeFromDir{&node{n}, self}
			file.Each(fn)
		}
	}
}

type dirNodeFromDir struct {
	*node
	parent *dirNodeClean
}

func (self *dirNodeFromDir) Parent() Node {
	return self.parent
}

func (self *dirNodeFromDir) Each(fn func(string, string)) {
	for _, n := range self.Node().Nodes {
		if n.Dir {
			dir := &dirNodeFromDir{&node{n}, self.parent}
			dir.Each(fn)
		} else {
			file := &fileNodeFromDir{&node{n}, self.parent}
			file.Each(fn)
		}
	}
}

type fileNodeClean struct {
	*node
}

func (self *fileNodeClean) Node() *etcd.Node {
	return self.node.Node()
}

func (self *fileNodeClean) Each(fn func(string, string)) {
	idx := strings.LastIndex(self.Key(), "/")
	k := self.Key()[idx+1:]
	fn(k, self.Value())
}

func (self *fileNodeClean) Parent() Node {
	return nil
}

type fileNodeFromDir struct {
	*node
	parent Node
}

func (self *fileNodeFromDir) Each(fn func(string, string)) {
	parent := self.Parent()
	fn(strings.TrimPrefix(self.Key(), parent.Key())[1:], self.Value())
}

func (self *fileNodeFromDir) Parent() Node {
	return self.parent
}
