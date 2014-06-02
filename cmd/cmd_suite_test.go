package main_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/shopkeep/fracker"

	"testing"
)

type TestEtcd struct{}

func (self *TestEtcd) Get(key string) (fracker.Node, error) {
	return &TestNode{key: key}, nil
}

type TestNode struct {
	key string
}

func (self *TestNode) IsFile() bool {
	return true
}

func (self *TestNode) Each(fn func(string, string)) {
	fn(self.key, "foo")
}

func TestCmd(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Cmd Suite")
}
