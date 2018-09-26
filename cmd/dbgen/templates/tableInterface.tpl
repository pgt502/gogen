package {{.TablePackage}}

import (
    "fmt"
    {{.Package}} "{{.PackagePath}}"
    {{range $path, $name := .Imports}}
    {{$name}} "{{$path}}"{{end}}
)

type {{.Name}}Table interface{
    Insert({{.Package}}.{{.Name}}) error
    Update({{.Package}}.{{.Name}}) error
    GetAll() ([]*{{.Package}}.{{.Name}}, error)
    Get({{range $i, $f := .PKFields}}{{if $i}},{{end}}{{$f.NameLower}} {{$f.Type}}{{end}}) ({{.Package}}.{{.Name}}, error)
}