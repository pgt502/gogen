package generate

import (
	"go/types"
	"strings"
)

type Field interface {
	Type() string
	Name() string
	NameLower() string
	Tag() string
	UnderlyingType() types.Type
	Index() int
}

type field struct {
	gen   Generator
	field *types.Var
	tag   string
	index int
}

func NewField(g Generator, f *types.Var, tag string, ind int) Field {
	fld := &field{
		gen:   g,
		field: f,
		tag:   tag,
		index: ind,
	}
	return fld
}

func (f field) UnderlyingType() types.Type {
	return f.field.Type()
}

func (f field) Type() string {
	return f.field.Type().String()
}

func (f field) Name() string {
	return f.field.Name()
}

func (f field) NameLower() string {
	return strings.ToLower(f.field.Name())
}

func (f field) Tag() string {
	return f.tag
}

func (f field) Index() int {
	return f.index
}
