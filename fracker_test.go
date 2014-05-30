package fracker_test

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
	var err error

	BeforeEach(func() {
		out = bytes.NewBuffer([]byte{})
		client = &StubEtcd{}
		fracker = f.New(client)
	})

	JustBeforeEach(func() {
		err = fracker.Frack(out, []string{"/foo"})
	})

	Context(`fracking a key that doesn't exist`, func() {
		BeforeEach(func() {
			client.StubGet = func(key string) (f.Node, error) {
				return nil, errors.New("no key")
			}
		})

		It(`returns an error`, func() {
			Expect(err).ToNot(BeNil())
		})
	})

	Context(`fracking a key that exists`, func() {
		Context(`and is a file`, func() {
			BeforeEach(func() {
				client.StubGet = func(key string) (f.Node, error) {
					n := NewStubFileNode(map[string]string{
						"/foo": "crunch",
					})
					return n, nil
				}
			})

			It(`returns an error`, func() {
				Expect(err).To(BeNil())
			})
		})

		Context(`and is a directory`, func() {
			BeforeEach(func() {
				client.StubGet = func(key string) (f.Node, error) {
					n := NewStubDirNode(map[string]string{
						"/foo/bar": "crunch",
						"/foo/baz": "munch",
					})
					return n, nil
				}
			})

			It(`doesn't return an error`, func() {
				Expect(err).To(BeNil())
			})

			It(`writes each value out in KEY=VALUE format (removing the prefix)`, func() {
				Expect(out.String()).To(Equal("BAR=crunch\nBAZ=munch\n"))
			})
		})
	})
})
