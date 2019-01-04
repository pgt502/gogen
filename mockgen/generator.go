package mockgen

import (
	"bytes"
	"html/template"
	"log"

	"github.com/pgt502/gogen/generate"
)

type Generator struct {
	generate.Generator
}

func NewGenerator(ifaceName, pkgName string) (*Generator, error) {
	base, err := generate.NewGenerator(ifaceName, pkgName)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	g := &Generator{
		Generator: base,
	}
	return g, nil
}

// Name returns the mock type's name by default it is {interfaceName}Mock. Overrides the base name
func (g *Generator) Name() string {
	return "Mock" + g.Generator.Name()
}

func (g *Generator) Generate(tpl string) string {
	var buf bytes.Buffer
	tmp := template.New("test")
	tmp, err := tmp.Parse(tpl)
	if err != nil {
		log.Printf("error parsing template: %s\n", err)
		return ""
	}
	tmp.Execute(&buf, g)

	return buf.String()
}
