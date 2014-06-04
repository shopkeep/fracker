package fracker_test

import (
	"github.com/onsi/ginkgo"
	"github.com/onsi/gomega"
	f "github.com/shopkeep/fracker"

	"log"
	"os"
	"testing"
)

type TestClient struct {
	StubGet func(string) (f.Node, error)
}

func (self *TestClient) Get(key string) (f.Node, error) {
	return self.StubGet(key)
}

func TestFracker(t *testing.T) {
	gomega.RegisterFailHandler(ginkgo.Fail)
	ginkgo.RunSpecs(t, "Fracker")
}

func init() {
	var err error
	var null *os.File
	if null, err = os.OpenFile(os.DevNull, os.O_WRONLY, 0666); err != nil {
		panic(err)
	}
	log.SetOutput(null)
}
