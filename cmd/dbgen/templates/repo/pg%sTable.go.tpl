package {{.PGTablePackage}}

import (
    "database/sql"
	"fmt"
	"strings"
    
    core "github.com/pgt502/gogen/core"
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
        //logging.LogErrore(err)
        return
    }

    return
}

func (t *pg{{.Name}}Table) Update(el {{.Package}}.{{.Name}}) (err error) {
    primaryKey := "{{range $i, $f := .PKFields}}{{if $i}} AND {{end}}{{$f.Column}}=${{$f.ColumnIndex}}{{end}}"
    valueSet := "{{range $i, $f := .NonPKFields}}{{if $i}},{{end}}{{$f.Column}}=${{$f.ColumnIndex}}{{end}}"

    sqlStatement := fmt.Sprintf("UPDATE %s SET %s WHERE %s", t.tableName, valueSet, primaryKey)
    _, err = t.db.Exec(sqlStatement,
        {{range $i, $f := .Fields}} el.{{$f.Name}},
        {{end}}
    )

    if err != nil {
        //logging.LogErrore(err)
        return
    }

    return
}


func (t *pg{{.Name}}Table) GetAll() (ret []*{{.Package}}.{{.Name}}, err error){
    sqlStatement := fmt.Sprintf(`SELECT %s
        FROM %s
    `, strings.Join(t.columns, ","), t.tableName)

    rows, err := t.db.Query(sqlStatement)
    if err != nil && err != sql.ErrNoRows {
        //logging.LogErrore(err)
        return
    }

    ret, err = ReadRows(rows)
    if err != nil {
        //logging.LogErrore(err)
        return
    }

    return
}

func (t *pg{{.Name}}Table) Get({{range $i, $f := .PKFields}}{{if $i}},{{end}}{{$f.NameLower}} {{$f.Type}}{{end}}) (ret {{.Package}}.{{.Name}}, err error) {
    where := "{{range $i, $f := .PKFields}}{{if $i}} AND {{end}}{{$f.Column}}=${{inc $i}}{{end}}"
    sqlStatement := fmt.Sprintf(`SELECT %s
        FROM %s
        WHERE %s`,
        strings.Join(t.columns, ","), 
        t.tableName,
        where,
    )

    row := t.db.QueryRow(sqlStatement,
        {{range $i, $f := .PKFields}}{{$f.NameLower}},
        {{end}}
    )
    ret, err = t.ReadRow(row)
    if err != nil && err != sql.ErrNoRows {
        //logging.LogErrore(err)
        return
    }
    return
}

func (t *pg{{.Name}}Table)Delete({{range $i, $f := .PKFields}}{{if $i}},{{end}}{{$f.NameLower}} {{$f.Type}}{{end}}) (err error){
    where := "{{range $i, $f := .PKFields}}{{if $i}} AND {{end}}{{$f.Column}}=${{inc $i}}{{end}}"
    sqlStatement := fmt.Sprintf(`DELETE 
        FROM %s
        WHERE %s`,        
        t.tableName,
        where,
    )

    _, err = t.db.Exec(sqlStatement,
        {{range $i, $f := .PKFields}}{{$f.NameLower}},
        {{end}}
    )

    if err != nil {
        //logging.LogErrore(err)
        return
    }

    return
}

func (t *pg{{.Name}}Table)ReadRows(rows core.ScannableExt) (items []*{{.Package}}.{{.Name}}, err error){
    for rows.Next() {
        var item {{.Package}}.{{.Name}}
        item, err = t.ReadRow(rows)
        if err != nil {
            //logging.LogErrore(err)
            return
        }
        items = append(items, &item)
    }
    return
}

func (t *pg{{.Name}}Table)ReadRow(row core.Scannable) (item {{.Package}}.{{.Name}}, err error){
    err = row.Scan(
        {{range $i, $f := .Fields}} &item.{{$f.Name}},
        {{end}}
    )
    return
}