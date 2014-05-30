package main

import (
	"github.com/codegangsta/cli"
	"github.com/shopkeep/fracker"

	"io"
	"log"
	"os"
	"strings"
)

// If no ETCD_HOSTS variable is defined, default to localhost port 4001
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
	app := App()
	app.Action = func(ctx *cli.Context) {
		if out, err := GetOutputFile(ctx); err != nil {
			log.Fatalln(err)
		} else {
			if err := frk.Frack(out, ctx.Args()); err != nil {
				log.Fatalln(err)
			}
		}
	}
	app.Run(os.Args)
}

// Builds the basic *cli.App with name and flag options
func App() *cli.App {
	app := cli.NewApp()
	app.Name = "fracker"
	app.Usage = "convert etcd hierarchies to environment variables"
	app.Flags = []cli.Flag{
		cli.StringFlag{Name: "output", Value: "", Usage: "output file (stdout by default)"},
		cli.BoolFlag{Name: "append", Usage: "append to output file instead of overwriting"},
	}
	return app
}

// Determines the output file based on the context
func GetOutputFile(ctx *cli.Context) (io.Writer, error) {
	if ctx.String("output") != "" {
		if ctx.Bool("append") {
			return os.OpenFile(ctx.String("output"), os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
		}
		return os.Create(ctx.String("output"))
	}
	return os.Stdout, nil
}
