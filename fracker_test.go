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
	var client *TestClient
	var err error

	BeforeEach(func() {
		out = bytes.NewBuffer([]byte{})
		client = &TestClient{}
		fracker = f.New(client)
	})

	Describe(`Frack`, func() {
		Context(`when a key doesn't exist`, func() {
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

		Context(`when a key does exist`, func() {
			Context(`and is a file`, func() {
				BeforeEach(func() {
					client.StubGet = func(key string) (f.Node, error) {
						n := &etcd.Node{
							Dir:   false,
							Key:   "/foo/baaaaaz",
							Value: "crunch",
						}
						return f.NewNode(n), nil
					}
				})

				JustBeforeEach(func() {
					err = fracker.Frack(out, []string{"/foo/baaaaaz"})
				})

				It(`does not return an error`, func() {
					Expect(err).To(BeNil())
				})

				It(`writes the last segment of the node's value in KEY=VALUE format`, func() {
					Expect(out.String()).To(ContainSubstring("BAAAAAZ=crunch\n"))
				})
			})

			Context(`and the name contains hyphens`, func() {
				BeforeEach(func() {
					client.StubGet = func(key string) (f.Node, error) {
						n := &etcd.Node{
							Dir:   false,
							Key:   "/foo/fun-times",
							Value: "crunch",
						}
						return f.NewNode(n), nil
					}
				})

				JustBeforeEach(func() {
					err = fracker.Frack(out, []string{"/foo/fun-times"})
				})

				It(`does not return an error`, func() {
					Expect(err).To(BeNil())
				})

				It(`writes the key's name replacing the hyphens with underscores`, func() {
					Expect(out.String()).To(ContainSubstring("FUN_TIMES=crunch\n"))
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
										&etcd.Node{
											Dir:   false,
											Key:   "/foo/bar/zoot-suit",
											Value: "funch",
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
					Expect(out.String()).To(ContainSubstring("BAR_BAZ=crunch\n"))
					Expect(out.String()).To(ContainSubstring("BAR_QUX=munch\n"))
					Expect(out.String()).To(ContainSubstring("BAR_ZOOT_SUIT=funch\n"))
				})
			})
		})
	})
})
