package main_test

import (
	"github.com/codegangsta/cli"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	// "github.com/shopkeep/fracker"
	cmd "github.com/shopkeep/fracker/cmd"

	"io"
	"os"
	"strings"
)

var _ = Describe("fracker", func() {
	var app *cli.App

	BeforeEach(func() {
		app = cmd.App()
	})

	Describe("setting the output file", func() {
		var out io.Writer
		var outf *os.File
		var err error
		var ok bool

		// args to pass to app.Run(). Easier than trying to construct a
		// properly parsed *cli.Context by hand
		var args []string

		BeforeEach(func() {
			app.Action = func(ctx *cli.Context) {
				out, err = cmd.GetOutputFile(ctx)
				outf, ok = out.(*os.File)
			}
		})

		JustBeforeEach(func() {
			app.Run(args)
		})

		Context("when no --output option is given", func() {
			BeforeEach(func() {
				args = strings.Split("key1 key2 key3", " ")
			})

			It("sets the output to stdout", func() {
				Expect(outf).To(Equal(os.Stdout))
			})
		})

		Context("when an --output option is given", func() {
			Context("when no --append option is given", func() {
				BeforeEach(func() {
					args = strings.Split("--output foo.env key1 key2 key3", " ")
				})

				It("sets the output to the given file", func() {
					info, _ := outf.Stat()
					Expect(info.Name()).To(Equal("foo.env"))
				})

				It("truncates the given file", func() {
					info, _ := outf.Stat()
					Expect(info.Size()).To(Equal(0))
				})
			})

			Context("when the --append option is given", func() {
				BeforeEach(func() {
					args = strings.Split("--output foo.env --append key1 key2 key3", " ")
					app.Action = func(ctx *cli.Context) {
						out, err = cmd.GetOutputFile(ctx)
						outf, ok = out.(*os.File)
					}
				})

				It("sets the output to the given file", func() {
					info, _ := outf.Stat()
					Expect(info.Name()).To(Equal("foo.env"))
				})

				It("sets append-only mode on the file", func() {
					info, _ := outf.Stat()
					Expect(info.Mode() & os.ModeAppend).ToNot(Equal(0))
				})
			})
		})
	})
})
