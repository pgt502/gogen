package {{.Package}}

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
				{{range $j, $p := $el.ParamTypes}}Param{{$j}} {{$p}}
				{{end}}
			}
			Returns struct{
				{{range $j, $r := $el.ReturnTypes}}Ret{{$j}} {{$r}}
				{{end}}
			}
			GetsCalled struct{
				Times int
			}
		}
		{{end}}
}

{{range $i, $el := .Methods}}
	func(m *{{$gen.Name}}) {{$el.Name}}({{range $j, $p := $el.ParamTypesVariadic}}{{if $j}},{{end}}p{{$j}} {{$p}}{{end}})({{range $j, $r := $el.ReturnTypes}}{{if $j}},{{end}}r{{$j}} {{$r}}{{end}}){
		m.{{$el.Name}}Call.GetsCalled.Times++		
		{{range $j, $p := $el.ParamTypes}}m.{{$el.Name}}Call.Receives.Param{{$j}}=p{{$j}}
		{{end}}
		{{range $j, $r := $el.ReturnTypes}}r{{$j}}=m.{{$el.Name}}Call.Returns.Ret{{$j}}
		{{end}}
		return
	}
{{end}}