package mockgen

import (
	"bytes"
	"fmt"
	"html/template"
	"log"

	"github.com/pgt502/gogen/generate"
)

type Generator struct {
	generate.Generator
}

func NewGenerator(ifaceName, pkgName, file string) (*Generator, error) {
	var base generate.Generator
	var err error
	if pkgName != "" {
		base, err = generate.NewGenerator(ifaceName, pkgName)
	} else if file != "" {
		base, err = generate.NewGeneratorFromFile(ifaceName, file)
	} else {
		return nil, fmt.Errorf("either package name or file need to be provided")
	}
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
