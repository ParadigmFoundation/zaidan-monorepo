package tests_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/ParadigmFoundation/zaidan-monorepo/tests"
)

var client *tests.Client

func TestTesting(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Zaidan")
}

var _ = BeforeSuite(func() {
	var err error
	client, err = tests.NewClient("ws://localhost:8000/ws")
	Expect(err).NotTo(HaveOccurred())
})

var _ = AfterSuite(func() {
	if client != nil {
		client.Close()
	}
})
