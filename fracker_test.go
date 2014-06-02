package fracker_test

import (
	"github.com/coreos/go-etcd/etcd"
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

	Context(`fracking a key that doesn't exist`, func() {
		BeforeEach(func() {
			client.StubGet = func(key string) (f.Node, error) {
				return nil, errors.New("no key")
			}
		})

		JustBeforeEach(func() {
			err = fracker.Frack(out, []string{"/foo"})
		})

		It(`returns an error`, func() {
			Expect(err).ToNot(BeNil())
		})
	})

	Context(`fracking a key that exists`, func() {
		Context(`and is a file`, func() {
			BeforeEach(func() {
				client.StubGet = func(key string) (f.Node, error) {
					n := &etcd.Node{
						Dir:   false,
						Key:   "/foo/bar",
						Value: "crunch",
					}
					return f.NewNode(n), nil
				}
			})

			JustBeforeEach(func() {
				err = fracker.Frack(out, []string{"/foo/bar"})
			})

			It(`does not return an error`, func() {
				Expect(err).To(BeNil())
			})

			It(`writes the last segment of the node's value in KEY=VALUE format`, func() {
				Expect(out.String()).To(Equal("BAR=crunch\n"))
			})
		})

		Context(`and is a directory`, func() {
			BeforeEach(func() {
				client.StubGet = func(key string) (f.Node, error) {
					n := &etcd.Node{
						Dir: true,
						Key: "/foo",
						Nodes: []*etcd.Node{
							&etcd.Node{
								Dir: true,
								Key: "/foo/bar",
								Nodes: []*etcd.Node{
									&etcd.Node{
										Dir:   false,
										Key:   "/foo/bar/baz",
										Value: "crunch",
									},
									&etcd.Node{
										Dir:   false,
										Key:   "/foo/bar/qux",
										Value: "munch",
									},
								},
							},
						},
					}
					return f.NewNode(n), nil
				}
			})

			JustBeforeEach(func() {
				err = fracker.Frack(out, []string{"/foo"})
			})

			It(`doesn't return an error`, func() {
				Expect(err).To(BeNil())
			})

			It(`writes each value out in KEY=VALUE format (removing the prefix)`, func() {
				Expect(out.String()).To(Equal("BAR_BAZ=crunch\nBAR_QUX=munch\n"))
			})
		})
	})
})
