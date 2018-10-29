package generate

import "strings"

type Param struct {
	Name string
	Type string
}

func (p Param) NameLower() string {
	return strings.ToLower(p.Name)
}

func (p Param) NameCapital() string {
	return strings.Title(p.Name)
}
