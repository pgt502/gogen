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
		var checkExpected = func(fields []Field, names, tags, types []string) {
			for i, field := range fields {
				Expect(field.Name()).To(Equal(names[i]), "fieldName: %s", names[i])
				Expect(field.Tag()).To(Equal(tags[i]), "tagName of %s", names[i])
				Expect(field.Type()).To(Equal(types[i]), "typeName of %s", names[i])
			}
		}
		Context("using package name", func() {
			Context("from struct", func() {
				BeforeEach(func() {
					generator, err = NewGenerator("TestType", "github.com/pgt502/gogen/testdata")
					Expect(err).NotTo(HaveOccurred())
				})
				It("should return correct fields", func() {
					fields := generator.Fields()
					Expect(len(fields)).To(Equal(3))
					names := []string{"Name", "Id", "CreatedAt"}
					tags := []string{"", `json:"id"`, `json:"created_at"`}
					types := []string{"string", "int64", "time.Time"}
					checkExpected(fields, names, tags, types)
				})
			})
		})
		Context("using file", func() {
			Context("from struct", func() {
				BeforeEach(func() {
					generator, err = NewGeneratorFromFile("TestType", "../testdata/order.go")
					Expect(err).NotTo(HaveOccurred())
				})
				It("should return correct fields", func() {
					fields := generator.Fields()
					Expect(len(fields)).To(Equal(3))
					names := []string{"Name", "Id", "CreatedAt"}
					tags := []string{"", `json:"id"`, `json:"created_at"`}
					types := []string{"string", "int64", "time.Time"}
					checkExpected(fields, names, tags, types)
				})
			})
		})
		Context("using package name", func() {
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
		Context("using file", func() {
			Context("from interface", func() {
				BeforeEach(func() {
					generator, err = NewGeneratorFromFile("TestInterface", "../testdata/order.go")
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
})
