package dbgen

import (
	"bytes"
	"fmt"
	"html/template"
	"strings"

	"github.com/pgt502/gogen/generate"
)

type Generator struct {
	generate.Generator
	name string
}

func NewGenerator(stuctName, pkgName string) (*Generator, error) {
	base, err := generate.NewGenerator(stuctName, pkgName)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	g := &Generator{
		Generator: base,
	}

	return g, nil
}

func (g *Generator) Table() string {
	return strings.ToLower(g.Name())
}

func (g *Generator) TablePackage() string {
	return "table"
}

func (g *Generator) PGTablePackage() string {
	return "postgress"
}

func (g *Generator) Columns() []string {

	fields := g.Fields()
	ret := make([]string, len(fields))
	for i, el := range fields {
		ret[i] = el.Column()
	}
	return ret
}

func (g *Generator) Fields() []DbField {
	if g.IsInterface() {
		return nil
	}
	bfields := g.Generator.Fields()

	dbFields := make([]DbField, len(bfields))
	for i, el := range bfields {
		f := NewDbField(el)
		dbFields[i] = f
	}
	return dbFields
}

func (g *Generator) PKFields() []DbField {
	if g.IsInterface() {
		return nil
	}
	all := g.Fields()

	pkFields := make([]DbField, 0)
	for _, el := range all {
		if el.IsPK() {
			pkFields = append(pkFields, el)
		}
	}
	return pkFields
}

func (g *Generator) NonPKFields() []DbField {
	if g.IsInterface() {
		return nil
	}
	all := g.Fields()

	nonPkFields := make([]DbField, 0)
	for _, el := range all {
		if !el.IsPK() {
			nonPkFields = append(nonPkFields, el)
		}
	}
	return nonPkFields
}

func (g *Generator) Generate(tpl string) string {
	var buf bytes.Buffer
	buf.WriteString("\n\n")
	// templates
	funcMap := template.FuncMap{
		"inc": func(i int) int {
			return i + 1
		},
	}
	tmp := template.New("test").Funcs(funcMap)
	tmp, err := tmp.Parse(tpl)
	if err != nil {
		fmt.Printf("error parsing template: %s\n", err)
		return ""
	}
	tmp.Execute(&buf, g)

	return buf.String()
}

var tableInterfaceTpl = `
	package {{.TablePackage}}

	import (
		"fmt"
		{{range $path, $name := .Imports}}
		{{$name}} "{{$path}}"{{end}}
	)

	type {{.Name}}Table interface{
		Insert({{.Package}}.{{.Name}}) error
		Update({{.Package}}.{{.Name}}) error
		GetAll() ([]*{{.Package}}.{{.Name}}, error)
	}
`
var pgtableStructTpl = `
	package {{.PGTablePackage}}

	import (
		"fmt"
		core "git.uk.guardtime.com/guardtime/volta/voltalib/core"
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


`
