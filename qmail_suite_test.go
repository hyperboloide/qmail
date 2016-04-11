package main_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

func TestQmail(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Qmail Suite")
}
