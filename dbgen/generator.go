package dbgen

import (
	"bytes"
	"fmt"
	"html/template"
	"log"
	"strings"

	"github.com/pgt502/gogen/generate"
)

// Generator is the main db generator type
type Generator struct {
	generate.Generator
	name string
	opts Options
}

// NewGenerator creates a new generator
func NewGenerator(stuctName string, options ...Option) (*Generator, error) {
	opts := &Options{}
	for _, op := range options {
		op(opts)
	}
	var base generate.Generator
	var err error
	if opts.PackageName != nil && *opts.PackageName != "" {
		base, err = generate.NewGenerator(stuctName, *opts.PackageName)
	} else if opts.File != nil && *opts.File != "" {
		base, err = generate.NewGeneratorFromFile(stuctName, *opts.File)
	} else {
		return nil, fmt.Errorf("either package name or file need to be provided")
	}

	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	g := &Generator{
		Generator: base,
		opts:      *opts,
	}
	return g, nil
}

func (g *Generator) Table() string {
	name := strings.ToLower(g.Name())
	if g.opts.Pluralise {
		name = fmt.Sprintf("%ss", name)
	}
	return name
}

func (g *Generator) Schema() string {
	return "public"
}

func (g *Generator) TablePackage() string {
	return "tables"
}

func (g *Generator) PGTablePackage() string {
	return "postgres"
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

	dbFields := make([]DbField, 0)
	for _, el := range bfields {
		f := NewDbField(el)
		if f.Ignore() {
			continue
		}
		dbFields = append(dbFields, f)
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
	funcMap := template.FuncMap{
		"inc": func(i int) int {
			return i + 1
		},
	}
	tmp := template.New("tpl").Funcs(funcMap)
	tmp, err := tmp.Parse(tpl)
	if err != nil {
		log.Printf("error parsing template: %s\n", err)
		return ""
	}
	tmp.Execute(&buf, g)

	return buf.String()
}
