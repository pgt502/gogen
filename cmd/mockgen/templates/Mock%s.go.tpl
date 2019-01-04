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
				{{range $j, $p := $el.ParamTypes}}{{$p.NameCapital}} {{$p.Type}}
				{{end}}
			}
			Returns struct{
				{{range $j, $r := $el.ReturnTypes}}{{$r.NameCapital}} {{$r.Type}}
				{{end}}
			}
			GetsCalled struct{
				Times int
			}
		}
		{{end}}
}

{{range $i, $el := .Methods}}
	func(m *{{$gen.Name}}) {{$el.Name}}({{range $j, $p := $el.ParamTypesVariadic}}{{if $j}},{{end}}{{$p.NameLower}} {{$p.Type}}{{end}})({{range $j, $r := $el.ReturnTypes}}{{if $j}},{{end}}{{$r.NameLower}} {{$r.Type}}{{end}}){
		m.{{$el.Name}}Call.GetsCalled.Times++		
		{{range $j, $p := $el.ParamTypes}}m.{{$el.Name}}Call.Receives.{{$p.NameCapital}}={{$p.NameLower}}
		{{end}}
		{{range $j, $r := $el.ReturnTypes}}{{$r.NameLower}}=m.{{$el.Name}}Call.Returns.{{$r.NameCapital}}
		{{end}}
		return
	}
{{end}}