package fracker

import (
	"fmt"
	"io"
	"path/filepath"
	"strings"
)

// A Fracker drills into etcd directories and outputs the values to a given outlet.
type Fracker interface {
	// Frack() looks up each of the given keys in etcd and walks the tree of nodes returned. Each leaf
	// value is converted into a environment variable declaration and written to the given io.Writer.
	Frack(io.Writer, []string) error
}

// New() creates a new Fracker.
func New(client Client) Fracker {
	return &fracker{client}
}

type fracker struct {
	client Client
}

// Frack() reads the configuration values out of etcd and writes them to stdout. Panics if a requested key
// is not found in etcd.
func (self *fracker) Frack(out io.Writer, keys []string) error {
	var node Node
	var err error
	env := make(map[string]string, 0)

	for _, key := range keys {
		key = filepath.Clean(key)

		node, err = self.client.Get(key)
		if err != nil {
			return err
		}

		node.Each(func(k, v string) {
			name := strings.TrimPrefix(k, key)
			if name == "" {
				name = filepath.Base(k)
			}

			name = strings.TrimPrefix(name, "/")
			name = strings.ToUpper(name)
			name = strings.Replace(name, "/", "_", -1)

			env[name] = v
		})
	}
	for name, val := range env {
		fmt.Fprintf(out, "%s=%s\n", name, val)
	}
	return nil
}
