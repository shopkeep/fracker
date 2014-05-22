package main

import (
	"github.com/codegangsta/cli"

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

	fracker := New(NewClient(hosts))

	app := cli.NewApp()
	app.Name = "fracker"
	app.Usage = "convert etcd hierarchies to environment variables"
	app.Action = func(ctx *cli.Context) {
		fracker.Frack(os.Stdout, ctx.Args())
	}

	app.Run(os.Args)
}
