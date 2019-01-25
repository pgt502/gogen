package repos

import (
    {{.Package}} "{{.PackagePath}}"
    {{range $path, $name := .Imports}}
    {{$name}} "{{$path}}"{{end}}
)

type {{.Name}}Repo interface{
    Create({{.NameLower}} {{.Package}}.{{.Name}}) (err error)
    Update({{.NameLower}} {{.Package}}.{{.Name}}) (err error)
    GetAll() ({{.NameLower}}s []*{{.Package}}.{{.Name}}, err error)
    Get({{range $i, $f := .PKFields}}{{if $i}},{{end}}{{$f.NameLower}} {{$f.Type}}{{end}}) ({{.NameLower}} {{.Package}}.{{.Name}}, err error)
    Delete({{range $i, $f := .PKFields}}{{if $i}},{{end}}{{$f.NameLower}} {{$f.Type}}{{end}}) (err error)
}

type {{.NameLower}}Repo struct{
    db tables.{{.Name}}Table
}

func New{{.Name}}Repo(tb tables.{{.Name}}Table) ({{.Name}}Repo){
    return &{{.NameLower}}Repo{
        db : tb,
    }
}

func (r *{{.NameLower}}Repo) Create(el {{.Package}}.{{.Name}}) (err error) {
    err = r.db.Insert(nil, el)
    return
}

func (r *{{.NameLower}}Repo) Update(el {{.Package}}.{{.Name}}) (err error) {
    err = r.db.Update(nil, el)
    return
}

func (r *{{.NameLower}}Repo) GetAll() (ret []*{{.Package}}.{{.Name}}, err error){
    ret, err = r.db.GetAll()
    return
}

func (r *{{.NameLower}}Repo) Get({{range $i, $f := .PKFields}}{{if $i}},{{end}}{{$f.NameLower}} {{$f.Type}}{{end}}) (ret {{.Package}}.{{.Name}}, err error) {
    ret, err = r.db.Get({{range $i, $f := .PKFields}}{{if $i}},{{end}}{{$f.NameLower}}{{end}})
    return
}

func (r *{{.NameLower}}Repo) Delete({{range $i, $f := .PKFields}}{{if $i}},{{end}}{{$f.NameLower}} {{$f.Type}}{{end}}) (err error) {
    err = r.db.Delete(nil, {{range $i, $f := .PKFields}}{{if $i}},{{end}}{{$f.NameLower}}{{end}})
    return
}