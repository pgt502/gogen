package generate

import (
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

// ParamTypes returns the list of types for the params
func (m Method) ParamTypes() []string {
	sig := m.Signature()
	types := m.listTypes(sig.Params())
	n := len(types)
	if n > 0 && sig.Variadic() {
		types[n-1] = strings.Replace(types[n-1], "[]", "...", 1)
	}
	return types
}

// ParamTypes returns the list of types for the params
func (m Method) ParamTypesClean() []string {
	sig := m.Signature()
	types := m.listTypeNames(sig.Params())
	n := len(types)
	if n > 0 && sig.Variadic() {
		types[n-1] = strings.Replace(types[n-1], "[]", "...", 1)
	}
	return types
}

func (m Method) ParamTypesCleanArr() []string {
	sig := m.Signature()
	types := m.listTypeNames(sig.Params())

	return types
}

// ReturnTypes returns the list of types for the params
func (m Method) ReturnTypes() []string {
	sig := m.Signature()
	return m.listTypes(sig.Results())
}

// ReturnTypes returns the list of types for the params
func (m Method) ReturnTypesClean() []string {
	sig := m.Signature()
	return m.listTypeNames(sig.Results())
}

func (m Method) listTypes(t *types.Tuple) []string {
	num := t.Len()
	list := make([]string, num)
	for i := 0; i < num; i++ {
		//list[i] = types.TypeString(t.At(i).Type(), m.gen.qf)
		list[i] = types.TypeString(t.At(i).Type(), nil)
	}
	return list
}

func (m Method) listTypeNames(t *types.Tuple) []string {
	num := t.Len()
	list := make([]string, num)
	for i := 0; i < num; i++ {
		//list[i] = types.TypeString(t.At(i).Type(), m.gen.qf)
		//list[i] = types.TypeString(t.At(i).Type(), nil)
		p := t.At(i).Type().String()
		list[i] = path.Base(p)
	}
	return list
}

func (m Method) Signature() *types.Signature {
	return m.fn.Type().(*types.Signature)
}
