package generate

import "go/types"

type Field struct {
	gen   Generator
	field *types.Var
	tag   string
}

func (f Field) Type() string {
	return f.field.Type().String()
}

func (f Field) Name() string {
	return f.field.Name()
}

func (f Field) Tag() string {
	return f.tag
}
