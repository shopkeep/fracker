package main

import (
	// TODO: remove dependence on etcd types, particularly *etcd.Response and *etcd.Node
	"github.com/coreos/go-etcd/etcd"

	"fmt"
	"io"
	"log"
	"strings"
)

// A Fracker drills into etcd directories and outputs the values to a given outlet.
type Fracker interface {
	// Frack() looks up each of the given keys in etcd and walks the tree of nodes returned. Each leaf
	// value is converted into a environment variable declaration and written to the given io.Writer.
	Frack(io.Writer, []string)
}

// New() creates a new Fracker.
func New(client EtcdKeyGetter) Fracker {
	return &fracker{client}
}

// The EtcdKeyGetter interface is to facilitate testing without using a running etcd process.
// XXX: can this be unexported?
type EtcdKeyGetter interface {
	Get(string, bool, bool) (*etcd.Response, error)
}

type fracker struct {
	client EtcdKeyGetter
}

// Frack() reads the configuration values out of etcd and writes them to stdout. Panics if a requested key
// is not found in etcd.
func (self *fracker) Frack(out io.Writer, keys []string) {
	var err error
	var resp *etcd.Response
	envVars := make(map[string]string, 0)
	for _, key := range keys {
		if resp, err = self.client.Get(key, false, true); err != nil {
			log.Panicln(err)
		}
		self.readNode(resp.Node, func(k, v string) {
			envVars[self.envVarName(k)] = v
		})
	}
	for key, val := range envVars {
		fmt.Fprintf(out, "%s=%s\n", key, val)
	}
}

// recurse over the given node and its children, yielding each non-directory
// node's key and value to the given function
func (self *fracker) readNode(node *etcd.Node, fn func(string, string)) {
	if !node.Dir {
		fn(node.Key, node.Value)
	} else {
		for _, n := range node.Nodes {
			self.readNode(n, fn)
		}
	}
}

// turns an etcd key name "/foo/bar" into an environment variable name "FOO_BAR"
// TODO: replace all non-valid characters with underscores
func (self *fracker) envVarName(name string) string {
	return strings.Replace(strings.ToUpper(name[1:]), "/", "_", -1)
}
