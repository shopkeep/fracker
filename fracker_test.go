package main_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	f "github.com/shopkeep/fracker"

	"bytes"
	"errors"
)

var _ = Describe("Fracker", func() {
	var out *bytes.Buffer
	var fracker f.Fracker
	var client *StubEtcd

	BeforeEach(func() {
		out = bytes.NewBuffer([]byte{})
		client = &StubEtcd{}
		fracker = f.New(client)
	})

	Context(`fracking a key that doesn't exist`, func() {
		BeforeEach(func() {
			client.StubGet = func(key string) (f.Node, error) {
				return nil, errors.New("no key")
			}
		})

		It(`panics`, func() {
			lambda := func() { fracker.Frack(out, []string{"foo"}) }
			Expect(lambda).To(Panic())
		})
	})

	Context(`fracking a key that exists`, func() {
		JustBeforeEach(func() {
			fracker.Frack(out, []string{"/foo"})
		})

		Context(`and is a file`, func() {
			BeforeEach(func() {
				client.StubGet = func(key string) (f.Node, error) {
					n := StubNode(map[string]string{
						"/foo": "crunch",
					})
					return n, nil
				}
			})

			It(`writes the value out in KEY=VALUE format`, func() {
				Expect(out.String()).To(Equal("FOO=crunch\n"))
			})
		})

		Context(`and is a directory`, func() {
			BeforeEach(func() {
				client.StubGet = func(key string) (f.Node, error) {
					n := StubNode(map[string]string{
						"/foo/bar": "crunch",
						"/foo/baz": "munch",
					})
					return n, nil
				}
			})

			It(`writes each value out in KEY=VALUE format`, func() {
				Expect(out.String()).To(Equal("FOO_BAR=crunch\nFOO_BAZ=munch\n"))
			})
		})
	})
})
