package generate

import (
	"fmt"
	"go/types"
	"path"
	"strings"
)

type Method struct {
	//gen *Generator
	fn *types.Func
}

func NewMethod(f *types.Func) *Method {
	return &Method{
		fn: f,
	}
}

// Name returns the method name
func (m Method) Name() string {
	return m.fn.Name()
}

// ParamTypesVariadic returns the list of types for the params with ... for variadic params
func (m Method) ParamTypesVariadic() []*Param {
	sig := m.Signature()
	types := m.listTypeNames(sig.Params())
	n := len(types)
	if n > 0 && sig.Variadic() {
		types[n-1].Type = strings.Replace(types[n-1].Type, "[]", "...", 1)
	}
	for i, el := range types {
		if el.Name == "" {
			el.Name = fmt.Sprintf("Param%d", i)
		}
	}
	return types
}

// ParamTypes returns the list of types for the params
func (m Method) ParamTypes() []*Param {
	sig := m.Signature()
	types := m.listTypeNames(sig.Params())
	for i, el := range types {
		if el.Name == "" {
			el.Name = fmt.Sprintf("Param%d", i)
		}
	}
	return types
}

// ReturnTypes returns the list of types for the params
func (m Method) ReturnTypes() []*Param {
	sig := m.Signature()
	types := m.listTypeNames(sig.Results())
	for i, el := range types {
		if el.Name == "" {

			el.Name = fmt.Sprintf("Ret%d", i)
		}

	}
	fmt.Printf("return types: %+v\n", types)
	return types
}

// ReturnTypesVariadic returns the list of types for the params with ... for variadic params
func (m Method) ReturnTypesVariadic() []*Param {
	sig := m.Signature()
	types := m.listTypeNames(sig.Results())
	n := len(types)
	if n > 0 && sig.Variadic() {
		types[n-1].Type = strings.Replace(types[n-1].Type, "[]", "...", 1)
	}
	for i, el := range types {
		if el.Name == "" {
			el.Name = fmt.Sprintf("Ret%d", i)
		}
	}
	return types
}

func (m Method) listTypeNames(t *types.Tuple) []*Param {
	num := t.Len()
	list := make([]*Param, num)
	for i := 0; i < num; i++ {
		p := t.At(i).Type().String()
		name := t.At(i).Name()

		fmt.Printf("type is: %s, name: %s\n", p, name)

		if strings.Contains(p, "/") {
			if strings.HasPrefix(p, "[]*") {
				list[i] = &Param{Type: fmt.Sprintf("[]*%s", path.Base(p)), Name: name}
			} else if strings.HasPrefix(p, "[]") {
				list[i] = &Param{Type: fmt.Sprintf("[]%s", path.Base(p)), Name: name}
			} else if strings.HasPrefix(p, "*") {
				list[i] = &Param{Type: fmt.Sprintf("*%s", path.Base(p)), Name: name}
			} else {
				list[i] = &Param{Type: path.Base(p), Name: name}
			}
		} else {
			list[i] = &Param{Type: p, Name: name}
		}
	}
	return list
}

func (m Method) Signature() *types.Signature {
	return m.fn.Type().(*types.Signature)
}
