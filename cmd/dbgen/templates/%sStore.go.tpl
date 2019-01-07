package stores

import (    
    {{.Package}} "{{.PackagePath}}"
    {{range $path, $name := .Imports}}
    {{$name}} "{{$path}}"{{end}}
)

type {{.Name}}Store interface{
    Add({{.Package}}.{{.Name}}) error
    Update({{.Package}}.{{.Name}}) error
    GetAll() ([]*{{.Package}}.{{.Name}}, error)
    Get({{range $i, $f := .PKFields}}{{if $i}},{{end}}{{$f.NameLower}} {{$f.Type}}{{end}}) ({{.Package}}.{{.Name}}, error)
    Delete({{range $i, $f := .PKFields}}{{if $i}},{{end}}{{$f.NameLower}} {{$f.Type}}{{end}}) error
}