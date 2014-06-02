package main_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

type TestNode struct {
	key string
}

func (self *TestNode) Each(fn func(string, string)) {
	fn(self.key, "foo")
}

func TestCmd(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Cmd Suite")
}
