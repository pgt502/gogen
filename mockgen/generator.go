package mockgen

import (
	"bytes"
	"fmt"
	"html/template"

	"github.com/pgt502/gogen/generate"
)

type Generator struct {
	generate.Generator
}

func NewGenerator(ifaceName, pkgName string) (*Generator, error) {
	base, err := generate.NewGenerator(ifaceName, pkgName)
	if err != nil {
		fmt.Println(err)
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

func (g *Generator) Generate() string {
	var buf bytes.Buffer
	buf.WriteString("\n\n")
	// templates
	tmp := template.New("test")
	tmp, err := tmp.Parse(mockTmpl)
	if err != nil {
		fmt.Printf("error parsing template: %s\n", err)
		return ""
	}
	tmp.Execute(&buf, g)

	return buf.String()
}

var mockTmpl = `package {{.Package}}

import (
	"fmt"
	{{range $path, $name := .Imports}}
	{{$name}} "{{$path}}"{{end}}
)

{{$gen := .}}

type {{.Name}} struct{
	{{range $i, $el := .Methods}}
		{{$el.Name}}Call struct{ 
			Receives struct{
				{{range $j, $p := $el.ParamTypesCleanArr}}Param{{$j}} {{$p}}
				{{end}}
			}
			Returns struct{
				{{range $j, $r := $el.ReturnTypesClean}}Ret{{$j}} {{$r}}
				{{end}}
			}
			GetsCalled struct{
				Times int
			}
		}
		{{end}}
}

{{range $i, $el := .Methods}}
	func(m *{{$gen.Name}}) {{$el.Name}}({{range $j, $p := $el.ParamTypesClean}}{{if $j}},{{end}}p{{$j}} {{$p}}{{end}})({{range $j, $r := $el.ReturnTypesClean}}{{if $j}},{{end}}r{{$j}} {{$r}}{{end}}){
		m.{{$el.Name}}Call.GetsCalled.Times++		
		{{range $j, $p := $el.ParamTypes}}m.{{$el.Name}}Call.Receives.Param{{$j}}=p{{$j}}
		{{end}}
		{{range $j, $r := $el.ReturnTypes}}r{{$j}}=m.{{$el.Name}}Call.Returns.Ret{{$j}}
		{{end}}
		return
	}
{{end}}

`
