package generate_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestImport(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Import Suite")
}
