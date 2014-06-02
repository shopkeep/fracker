package main_test

import (
	"github.com/codegangsta/cli"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	cmd "github.com/shopkeep/fracker/cmd"

	"io"
	"io/ioutil"
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
				args = strings.Split("fracker key1 key2 key3", " ")
			})

			It("sets the output to stdout", func() {
				Expect(outf).To(Equal(os.Stdout))
			})
		})

		Context("when an --output option is given", func() {
			BeforeEach(func() {
				args = strings.Split("fracker --output=foo.env key1 key2 key3", " ")
			})

			AfterEach(func() {
				os.Remove("foo.env")
			})

			It("sets the output to the given file", func() {
				info, _ := outf.Stat()
				Expect(info.Name()).To(Equal("foo.env"))
			})

			Context("when the file already exists", func() {
				BeforeEach(func() {
					if err := ioutil.WriteFile("foo.env", []byte("woohoo"), 0666); err != nil {
						panic(err)
					}
				})

				It("truncates the given file", func() {
					info, _ := outf.Stat()
					Expect(info.Size()).To(Equal(int64(0)))
				})
			})

			Context("when the file doesn't already exist", func() {
				BeforeEach(func() {
					_ = os.Remove("foo.env")
				})

				It("creates the given file", func() {
					_, err := outf.Stat()
					Expect(err).To(BeNil())
				})
			})
		})
	})
})
