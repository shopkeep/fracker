package main

import (
	"github.com/codegangsta/cli"
	"github.com/coreos/go-etcd/etcd"

	"fmt"
	"io"
	"log"
	"os"
	"strings"
)

const DefaultEtcdHost string = "http://127.0.0.1:4001"

func main() {
	var hosts []string
	var env string
	if env = os.Getenv("ETCD_HOSTS"); env == "" {
		hosts = []string{DefaultEtcdHost}
	} else {
		hosts = strings.Split(env, ",")
	}

	client := etcd.NewClient(hosts)
	fracker := NewFracker(client)

	app := cli.NewApp()
	app.Name = "fracker"
	app.Usage = "convert etcd hierarchies to environment variables"
	app.Action = func(ctx *cli.Context) {
		fracker.Frack(os.Stdout, ctx.Args())
	}

	app.Run(os.Args)
}

type Fracker interface {
	Frack(io.Writer, []string)
}

func NewFracker(client *etcd.Client) Fracker {
	return &fracker{client}
}

type fracker struct {
	client *etcd.Client
}

func (self *fracker) Frack(out io.Writer, keys []string) {
	var err error
	var resp *etcd.Response
	envVars := make(map[string]string, 0)
	for _, key := range keys {
		if resp, err = self.client.Get(key, false, true); err != nil {
			log.Panicln(err)
		}
		readNode(resp.Node, func(k, v string) {
			envVars[envVarName(k)] = v
		})
	}
	for key, val := range envVars {
		fmt.Fprintf(out, "%s=%s\n", key, val)
	}
}

// recurse over the given node and its children, yielding each non-directory
// node's key and value to the given function
func readNode(node *etcd.Node, fn func(string, string)) {
	if !node.Dir {
		fn(node.Key, node.Value)
	} else {
		for _, n := range node.Nodes {
			readNode(n, fn)
		}
	}
}

// turns an etcd key name "/foo/bar" into an environment variable name "FOO_BAR"
func envVarName(name string) string {
	return strings.Replace(strings.ToUpper(name[1:]), "/", "_", -1)
}
