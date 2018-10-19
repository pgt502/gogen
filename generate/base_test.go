package generate_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	. "github.com/pgt502/gogen/generate"
)

var _ = Describe("Base", func() {
	Describe("getting information", func() {
		var (
			generator Generator
			err       error
		)
		Context("from struct", func() {
			BeforeEach(func() {
				generator, err = NewGenerator("TestType", "github.com/pgt502/gogen/testdata")
				Expect(err).NotTo(HaveOccurred())
			})
			It("should return correct fields", func() {
				fields := generator.Fields()
				Expect(len(fields)).To(Equal(3))
				Expect(fields[0].Name()).To(Equal("Name"), "Name")
				Expect(fields[1].Name()).To(Equal("Id"), "Id")
				Expect(fields[2].Name()).To(Equal("CreatedAt"), "CreatedAt")
				Expect(fields[1].Tag()).To(Equal(`json:"id"`), "Id tag")
				Expect(fields[2].Tag()).To(Equal(`json:"created_at"`), "CreatedAt tag")
			})
		})
		Context("from interface", func() {
			BeforeEach(func() {
				generator, err = NewGenerator("TestInterface", "github.com/pgt502/gogen/testdata")
				Expect(err).NotTo(HaveOccurred())
			})
			It("should return correct methods", func() {
				methods := generator.Methods()
				Expect(len(methods)).To(Equal(4))
				Expect(methods[0].Name()).To(Equal("Method1"))
				Expect(methods[1].Name()).To(Equal("Method2"))
				Expect(methods[2].Name()).To(Equal("Method3"))
				Expect(methods[3].Name()).To(Equal("Method4"))
				Expect(methods[0].ParamTypes()).To(Equal([]string{}), "1st method params")
				Expect(methods[1].ParamTypes()).To(Equal([]string{"string", "interface{}"}), "2nd method params")
				Expect(methods[2].ParamTypes()).To(Equal([]string{"[]*testdata.TestType"}), "3rd method params")
				Expect(methods[3].ParamTypes()).To(Equal([]string{"[]string"}), "4th method params")
				Expect(methods[3].ParamTypesVariadic()).To(Equal([]string{"...string"}), "4th method params variadic")
				Expect(methods[0].ReturnTypes()).To(Equal([]string{"error"}), "1st method return types")
				Expect(methods[1].ReturnTypes()).To(Equal([]string{"interface{}", "error"}), "2nd method return types")
				Expect(methods[2].ReturnTypes()).To(Equal([]string{"[]testdata.TestType"}), "3rd method return types")
				Expect(methods[3].ReturnTypes()).To(Equal([]string{}), "4th method return types")
			})
		})
	})
})
