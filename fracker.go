package fracker

import (
	"fmt"
	"io"
	"strings"
)

// A Fracker drills into etcd directories and outputs the values to a given outlet.
type Fracker interface {
	// Frack() looks up each of the given keys in etcd and walks the tree of nodes returned. Each leaf
	// value is converted into a environment variable declaration and written to the given io.Writer.
	Frack(io.Writer, []string) error
}

// New() creates a new Fracker.
func New(client EtcdClient) Fracker {
	return &fracker{client}
}

type fracker struct {
	client EtcdClient
}

// Frack() reads the configuration values out of etcd and writes them to stdout. Panics if a requested key
// is not found in etcd.
func (self *fracker) Frack(out io.Writer, keys []string) error {
	envVars := make(map[string]string, 0)
	for _, key := range keys {
		if node, err := self.client.Get(key); err != nil {
			return err
		} else {
			node.Each(func(k, v string) {
				envVars[self.envVarName(k)] = v
			})
		}
	}
	for key, val := range envVars {
		fmt.Fprintf(out, "%s=%s\n", key, val)
	}
	return nil
}

// turns an etcd key name "/foo/bar" into an environment variable name "FOO_BAR"
// TODO: replace all non-valid characters with underscores
func (self *fracker) envVarName(name string) string {
	return strings.Replace(strings.ToUpper(name[1:]), "/", "_", -1)
}
