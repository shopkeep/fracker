package main

import (
	"github.com/coreos/go-etcd/etcd"

	"fmt"
	"io"
	"log"
	"strings"
)

type Fracker interface {
	Frack(io.Writer, []string)
}

func New(client EtcdKeyGetter) Fracker {
	return &fracker{client}
}

// this is just a wrapper to facilitate testing
type EtcdKeyGetter interface {
	Get(string, bool, bool) (*etcd.Response, error)
}

type fracker struct {
	client EtcdKeyGetter
}

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
func (self *fracker) envVarName(name string) string {
	return strings.Replace(strings.ToUpper(name[1:]), "/", "_", -1)
}
