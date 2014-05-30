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
	env := make(map[string]string, 0)
	for _, key := range keys {
		// normalize the key so that it's in a known format
		key = self.normalizeKeyName(key)
		if node, err := self.client.Get(key); err != nil {
			return err
		} else {
			if node.IsFile() {
				key = key[:len(key)-1]
				key = key[:strings.LastIndex(key, "/")+1]
			}
			node.Each(func(k, v string) {
				n := self.etcdPathToEnvVarName(key, k)
				env[n] = v
			})
		}
	}
	for key, val := range env {
		fmt.Fprintf(out, "%s=%s\n", key, val)
	}
	return nil
}

func (self *fracker) normalizeKeyName(key string) string {
	return "/" + strings.Trim(key, "/") + "/"
}

func (self *fracker) etcdPathToEnvVarName(prefix, key string) string {
	str := strings.TrimPrefix(key, prefix)
	str = strings.ToUpper(str)
	str = strings.Replace(str, "/", "_", -1)
	return str
}
