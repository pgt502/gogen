package repos

import (
    {{.Package}} "{{.PackagePath}}"
    {{range $path, $name := .Imports}}
    {{$name}} "{{$path}}"{{end}}
)

type {{.Name}}Repo interface{
    Create({{.Package}}.{{.Name}}) error
    Update({{.Package}}.{{.Name}}) error
    GetAll() ([]*{{.Package}}.{{.Name}}, error)
    Get({{range $i, $f := .PKFields}}{{if $i}},{{end}}{{$f.NameLower}} {{$f.Type}}{{end}}) ({{.Package}}.{{.Name}}, error)
}

type {{.NameLower}}Repo struct{
    db tables.{{.Name}}Db
}

func New{{.Name}}Repo(tb tables.{{.Name}}Table) ({{.Name}}Repo){
    return &{{.NameLower}}Repo{
        db : tb,
    }
}

func (r *{{.NameLower}}Repo) Create(el {{.Package}}.{{.Name}}) (err error) {
    err = r.db.Insert(el)
    return
}

func (r *{{.NameLower}}Repo) Update(el {{.Package}}.{{.Name}}) (err error) {
    err = r.db.Update(el)
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