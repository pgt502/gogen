package {{.PGTablePackage}}

import (
    "fmt"
    core "git.uk.guardtime.com/guardtime/volta/voltalib/core"
    {{.Package}} "{{.PackagePath}}"
    {{range $path, $name := .Imports}}
    {{$name}} "{{$path}}"{{end}}
)

type pg{{.Name}}Table struct{
    tableName string
    db core.Queryable
    columns []string
    values string
}

func NewPg{{.Name}}Table(q core.Queryable) (t tables.{{.Name}}Table){
    return &pg{{.Name}}Table{
        tableName : "{{.Table}}",
        db: q,
        columns : []string{
            {{range $i, $col := .Columns}} "{{$col}}",
            {{end}}
        },
        values: "{{range $i, $col := .Columns}}{{if $i}},{{end}}${{inc $i}}{{end}}",
    }
}

func (t *pg{{.Name}}Table) Insert(el {{.Package}}.{{.Name}}) (err error) {
    sqlStatement := fmt.Sprintf("INSERT INTO %s (%s) VALUES (%s)", t.tableName, strings.Join(t.columns, ","), t.values)

    _, err = t.db.Exec(sqlStatement,
        {{range $i, $f := .Fields}} el.{{$f.Name}},
        {{end}}
    )

    if err != nil {
        
    }

    return
}

func (t *pg{{.Name}}Table) Update(el {{.Package}}.{{.Name}}) (err error) {
    return
}


//GetAll() ([]*{{.Package}}.{{.Name}}, error)