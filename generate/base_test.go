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
				Expect(fields[0].Type()).To(Equal("string"), "Name type")
				Expect(fields[1].Type()).To(Equal("int64"), "Id type")
				Expect(fields[2].Type()).To(Equal("time.Time"), "CreatedAt type")
			})
		})
		Context("from interface", func() {
			BeforeEach(func() {
				generator, err = NewGenerator("TestInterface", "github.com/pgt502/gogen/testdata")
				Expect(err).NotTo(HaveOccurred())
			})
			It("should return correct methods with params and return types", func() {
				methods := generator.Methods()
				Expect(len(methods)).To(Equal(5))
				Expect(methods[0].Name()).To(Equal("Method1"))
				Expect(methods[1].Name()).To(Equal("Method2"))
				Expect(methods[2].Name()).To(Equal("Method3"))
				Expect(methods[3].Name()).To(Equal("Method4"))
				Expect(methods[0].ParamTypes()).To(Equal([]*Param{}), "1st method params")
				Expect(methods[1].ParamTypes()).To(Equal([]*Param{&Param{Type: "string", Name: "Param0"}, &Param{Type: "interface{}", Name: "Param1"}}), "2nd method params")
				Expect(methods[2].ParamTypes()).To(Equal([]*Param{&Param{Type: "[]*testdata.TestType", Name: "Param0"}}), "3rd method params")
				Expect(methods[3].ParamTypes()).To(Equal([]*Param{&Param{Type: "[]string", Name: "inputs"}}), "4th method params")
				Expect(methods[3].ParamTypesVariadic()).To(Equal([]*Param{&Param{Type: "...string", Name: "inputs"}}), "4th method params variadic")
				Expect(methods[4].ParamTypes()).To(Equal([]*Param{&Param{Type: "*testdata.TestType", Name: "test"}}), "5th method params")
				Expect(methods[0].ReturnTypes()).To(Equal([]*Param{&Param{Type: "error", Name: "Ret0"}}), "1st method return types")
				Expect(methods[1].ReturnTypes()).To(Equal([]*Param{&Param{Type: "interface{}", Name: "Ret0"}, &Param{Type: "error", Name: "Ret1"}}), "2nd method return types")
				Expect(methods[2].ReturnTypes()).To(Equal([]*Param{&Param{Type: "[]testdata.TestType", Name: "Ret0"}}), "3rd method return types")
				Expect(methods[3].ReturnTypes()).To(Equal([]*Param{}), "4th method return types")
				Expect(methods[4].ReturnTypes()).To(Equal([]*Param{&Param{Type: "*testdata.TestType", Name: "ret"}, &Param{Type: "error", Name: "err"}}), "5th method return types")
			})
		})
	})
})
