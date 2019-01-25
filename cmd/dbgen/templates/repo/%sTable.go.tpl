package {{.TablePackage}}

import (    
    "database/sql"
    {{.Package}} "{{.PackagePath}}"
    {{range $path, $name := .Imports}}
    {{$name}} "{{$path}}"{{end}}
)

type {{.Name}}Table interface{    
    Insert(tx *sql.Tx, {{.NameLower}} {{.Package}}.{{.Name}}) (err error)
    Update(tx *sql.Tx, {{.NameLower}} {{.Package}}.{{.Name}}) (err error)
    GetAll() ({{.NameLower}}s []*{{.Package}}.{{.Name}}, err error)
    Get({{range $i, $f := .PKFields}}{{if $i}},{{end}}{{$f.NameLower}} {{$f.Type}}{{end}}) ({{.NameLower}} {{.Package}}.{{.Name}}, err error)
    Delete(tx *sql.Tx, {{range $i, $f := .PKFields}}{{if $i}},{{end}}{{$f.NameLower}} {{$f.Type}}{{end}}) (err error)
}