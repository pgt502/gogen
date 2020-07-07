package generate_test

import (
	"fmt"
	"log"

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

			log.Printf("package name is: %s, path: [%s]\n", pkg.Name(), pkg.Path())
			obj := pkg.Scope().Lookup("Order")
			fmt.Printf("object: %+v\n", obj)
			fmt.Print(pkg.Imports())
		})
		FIt("should import a file outside the project", func() {
			imp := NewImporter()
			pkg, err := imp.ImportFile("/Users/pawel/dev/go/hello/pkg/types/person.go")
			Expect(err).To(BeNil(), "error importing file")
			Expect(pkg).NotTo(BeNil(), "package nil")

			log.Printf("package name is: %s, path: [%s]\n", pkg.Name(), pkg.Path())
			obj := pkg.Scope().Lookup("Person")
			fmt.Printf("object: %+v\n", obj)
			fmt.Print(pkg.Imports())
		})
		It("using importer", func() {
			imp := importer.DefaultWithTestFiles()
			pkgName := "github.com/pgt502/gogen/testdata"
			pkg, err := imp.Import(pkgName)
			Expect(err).To(BeNil(), "error importing file")
			Expect(pkg).NotTo(BeNil(), "package nil")
			log.Printf("package name is: %s, path: [%s]\n", pkg.Name(), pkg.Path())
			obj := pkg.Scope().Lookup("Order")
			fmt.Printf("object: %+v\n", obj)
		})
	})
})
