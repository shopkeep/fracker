package main_test

import (
	"github.com/coreos/go-etcd/etcd"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	f "github.com/shopkeep/fracker"

	"bytes"
	"errors"
	"testing"
)

type StubEtcd struct {
	StubGet func(string, bool, bool) (*etcd.Response, error)
}

func (self *StubEtcd) Get(key string, sort, rec bool) (*etcd.Response, error) {
	return self.StubGet(key, sort, rec)
}

func TestFracker(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Fracker")
}

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
			client.StubGet = func(key string, sort, rec bool) (*etcd.Response, error) {
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
				client.StubGet = func(key string, sort, rec bool) (*etcd.Response, error) {
					resp := &etcd.Response{
						Node: &etcd.Node{
							Dir:   false,
							Key:   "/foo",
							Value: "crunch",
						},
					}
					return resp, nil
				}
			})

			It(`writes the value out in KEY=VALUE format`, func() {
				Expect(out.String()).To(Equal("FOO=crunch\n"))
			})
		})

		Context(`and is a directory`, func() {
			BeforeEach(func() {
				client.StubGet = func(key string, sort, rec bool) (*etcd.Response, error) {
					resp := &etcd.Response{
						Node: &etcd.Node{
							Dir: true,
							Key: "/foo",
							Nodes: []*etcd.Node{
								&etcd.Node{
									Dir:   false,
									Key:   "/foo/bar",
									Value: "crunch",
								},
								&etcd.Node{
									Dir:   false,
									Key:   "/foo/baz",
									Value: "munch",
								},
							},
						},
					}
					return resp, nil
				}
			})

			It(`writes each value out in KEY=VALUE format`, func() {
				Expect(out.String()).To(Equal("FOO_BAR=crunch\nFOO_BAZ=munch\n"))
			})
		})
	})
})
