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
}