package fracker_test

import (
	"github.com/coreos/go-etcd/etcd"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	f "github.com/shopkeep/fracker"
)

var _ = Describe("Node", func() {
	// a counter to check the number of times a function yields
	var yields int
	var valMap map[string]string
	var node f.Node

	BeforeEach(func() {
		yields = 0
		valMap = make(map[string]string, 0)
	})

	JustBeforeEach(func() {
		node.Each(func(k, v string) {
			valMap[k] = v
			yields++
		})
	})

	Context(`wrapping a file node`, func() {
		BeforeEach(func() {
			node = f.NewNode(&etcd.Node{
				Dir:   false,
				Key:   "/foo",
				Value: "1234",
			})
		})

		It(`yields once`, func() {
			Expect(yields).To(Equal(1))
		})

		It(`yields the file's name and value`, func() {
			Expect(valMap).To(Equal(map[string]string{
				"/foo": "1234",
			}))
		})
	})

	Context(`wrapping a directory node`, func() {
		BeforeEach(func() {
			node = f.NewNode(&etcd.Node{
				Dir: true,
				Key: "/foo",
				Nodes: []*etcd.Node{
					&etcd.Node{
						Dir:   false,
						Key:   "/foo/bar",
						Value: "1234",
					},
					&etcd.Node{
						Dir: true,
						Key: "/foo/baz",
						Nodes: []*etcd.Node{
							&etcd.Node{
								Dir:   false,
								Key:   "/foo/baz/qux",
								Value: "crunch",
							},
							&etcd.Node{
								Dir:   false,
								Key:   "/foo/baz/goo",
								Value: "munch",
							},
						},
					},
				},
			})
		})

		It(`yields once for each file`, func() {
			Expect(yields).To(Equal(3))
		})

		It(`yields each pair of file names and values`, func() {
			Expect(valMap).To(Equal(map[string]string{
				"/foo/bar":     "1234",
				"/foo/baz/qux": "crunch",
				"/foo/baz/goo": "munch",
			}))
		})
	})
})
