package main

import (
	"github.com/codegangsta/cli"
	"github.com/shopkeep/fracker"

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

	client := fracker.NewClient(hosts)
	frk := fracker.New(client)

	app := cli.NewApp()
	app.Name = "fracker"
	app.Usage = "convert etcd hierarchies to environment variables"
	app.Flags = []cli.Flag{
		cli.StringFlag{Name: "output", Value: "", Usage: "output file (stdout by default)"},
		cli.BoolFlag{Name: "append", Usage: "append to output file instead of overwriting"},
	}

	app.Action = func(ctx *cli.Context) {
		var out io.Writer = os.Stdout
		var err error

		if ctx.String("output") != "" {
			if ctx.Bool("append") {
				out, err = os.OpenFile(ctx.String("output"), os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
			} else {
				out, err = os.Create(ctx.String("output"))
			}
			if err != nil {
				log.Fatalln(err)
			}
		}

		if err := frk.Frack(out, ctx.Args()); err != nil {
			log.Fatalln(err)
		}
	}

	app.Run(os.Args)
}
