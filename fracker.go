package main

import (
	"github.com/codegangsta/cli"
	"github.com/coreos/go-etcd/etcd"

	"fmt"
	"log"
	"os"
	"strings"
)

func main() {
	var hosts []string
	var env string
	if env = os.Getenv("ETCD_HOSTS"); env == "" {
		hosts = []string{"http://127.0.0.1:4001"}
	} else {
		hosts = strings.Split(env, ",")
	}

	client := etcd.NewClient(hosts)

	app := cli.NewApp()
	app.Name = "fracker"
	app.Usage = "fracker [key...]"
	app.Action = func(ctx *cli.Context) {
		var err error
		var resp *etcd.Response
		envVars := make(map[string]string, 0)
		for _, key := range ctx.Args() {
			//for _, key := range os.Args[1:] {
			if resp, err = client.Get(key, false, true); err != nil {
				log.Panicln(err)
			}
			recurseOverNodes(resp.Node, func(k, v string) {
				envVars[envVarName(k)] = v
			})
		}
		for key, val := range envVars {
			fmt.Fprintf(os.Stdout, "%s=%s\n", key, val)
		}
	}

	app.Run(os.Args)
}

// recurse over the given node and its children, yielding each non-directory
// node's key and value to the given function
func recurseOverNodes(node *etcd.Node, fn func(string, string)) {
	if !node.Dir {
		fn(node.Key, node.Value)
	} else {
		for _, n := range node.Nodes {
			recurseOverNodes(n, fn)
		}
	}
}

// turns an etcd key name "/foo/bar" into an environment variable name "FOO_BAR"
func envVarName(name string) string {
	return strings.Replace(strings.ToUpper(name[1:]), "/", "_", -1)
}
