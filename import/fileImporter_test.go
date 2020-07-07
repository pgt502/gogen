package generate_test

import (
	"github.com/ernesto-jimenez/gogen/importer"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	. "github.com/pgt502/gogen/import"
)

var _ = Describe("FileImporter", func() {
	Context("importing a file", func() {
		It("should import the file", func() {
			imp := NewImporter()
			pkg, err := imp.ImportFile("../testdata/order.go")
			Expect(err).To(BeNil(), "error importing file")
			Expect(pkg).NotTo(BeNil(), "package nil")
			Expect(pkg.Name()).To(Equal("testdata"), "package name")

			obj := pkg.Scope().Lookup("Order")
			Expect(obj).NotTo(BeNil(), "object nil")
			Expect(obj.Name()).To(Equal("Order"), "object name")
		})
		It("using importer", func() {
			imp := importer.DefaultWithTestFiles()
			pkgName := "github.com/pgt502/gogen/testdata"
			pkg, err := imp.Import(pkgName)
			Expect(err).To(BeNil(), "error importing file")
			Expect(pkg).NotTo(BeNil(), "package nil")
			Expect(pkg.Name()).To(Equal("testdata"), "package name")

			obj := pkg.Scope().Lookup("Order")
			Expect(obj).NotTo(BeNil(), "object nil")
			Expect(obj.Name()).To(Equal("Order"), "object name")
		})
	})
})
