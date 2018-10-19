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
func (m Method) ParamTypesVariadic() []string {
	sig := m.Signature()
	types := m.listTypeNames(sig.Params())
	n := len(types)
	if n > 0 && sig.Variadic() {
		types[n-1] = strings.Replace(types[n-1], "[]", "...", 1)
	}
	return types
}

// ParamTypes returns the list of types for the params
func (m Method) ParamTypes() []string {
	sig := m.Signature()
	types := m.listTypeNames(sig.Params())
	return types
}

// ReturnTypes returns the list of types for the params
func (m Method) ReturnTypes() []string {
	sig := m.Signature()
	return m.listTypeNames(sig.Results())
}

// ReturnTypesVariadic returns the list of types for the params with ... for variadic params
func (m Method) ReturnTypesVariadic() []string {
	sig := m.Signature()
	types := m.listTypeNames(sig.Results())
	n := len(types)
	if n > 0 && sig.Variadic() {
		types[n-1] = strings.Replace(types[n-1], "[]", "...", 1)
	}
	return types
}

func (m Method) listTypeNames(t *types.Tuple) []string {
	num := t.Len()
	list := make([]string, num)
	for i := 0; i < num; i++ {
		p := t.At(i).Type().String()
		if strings.Contains(p, "/") {
			if strings.HasPrefix(p, "[]*") {
				list[i] = fmt.Sprintf("[]*%s", path.Base(p))
			} else if strings.HasPrefix(p, "[]") {
				list[i] = fmt.Sprintf("[]%s", path.Base(p))
			} else {
				list[i] = path.Base(p)
			}
		} else {
			list[i] = p
		}
	}
	return list
}

func (m Method) Signature() *types.Signature {
	return m.fn.Type().(*types.Signature)
}
