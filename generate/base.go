package generate

import (
	"fmt"
	"go/types"
	"strings"

	"github.com/ernesto-jimenez/gogen/importer"
	"github.com/ernesto-jimenez/gogen/imports"
)

type Generator interface {
	Imports() map[string]string
	Name() string
	NameLower() string
	Package() string
	PackagePath() string
	Methods() []Method
	Fields() []Field
	AllFields() []Field
	IsInterface() bool
}

type generator struct {
	pkg         *types.Package
	typeName    string
	iface       *types.Interface
	strct       *types.Struct
	isInterface bool
	name        string
}

func NewGenerator(typeName, pkgName string) (Generator, error) {
	imp := importer.DefaultWithTestFiles()
	pkg, err := imp.Import(pkgName)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	fmt.Printf("package name is: %s, path: [%s]\n", pkg.Name(), pkg.Path())
	scope := pkg.Scope()
	if scope == nil {
		err = fmt.Errorf("package scope is nil")
		fmt.Println(err)
		return nil, err
	}
	obj := scope.Lookup(typeName)
	if obj == nil {
		err = fmt.Errorf("type '%s' does not exist", typeName)
		return nil, err
	}
	g := &generator{
		typeName: typeName,
		pkg:      pkg,
	}
	if types.IsInterface(obj.Type()) {
		g.iface = obj.Type().Underlying().(*types.Interface).Complete()
		g.isInterface = true
	} else {
		g.strct = obj.Type().Underlying().(*types.Struct)
	}

	return g, nil
}

// Imports returns all the packages that have to be imported for the
func (g *generator) Imports() map[string]string {
	imports := imports.New(g.Package())
	if g.isInterface {
		for _, m := range g.Methods() {
			s := m.Signature()
			imports.AddImportsFrom(s.Params())
			imports.AddImportsFrom(s.Results())
		}
	} else {
		for _, m := range g.Fields() {
			imports.AddImportsFrom(m.UnderlyingType())
		}
	}

	return imports.Imports()
}

func (g *generator) AllFields() []Field {
	if g.isInterface {
		return nil
	}
	numFields := g.strct.NumFields()
	fields := make([]Field, numFields)
	for i := 0; i < numFields; i++ {
		fields[i] = NewField(g, g.strct.Field(i), g.strct.Tag(i), i)
	}
	return fields
}

func (g *generator) Fields() []Field {
	if g.isInterface {
		return nil
	}
	numFields := g.strct.NumFields()
	fields := make([]Field, 0)
	for i := 0; i < numFields; i++ {
		f := g.strct.Field(i)
		if f.Exported() {
			fields = append(fields, NewField(g, f, g.strct.Tag(i), i))
		}
	}
	return fields
}

// Name returns the type's name
func (g *generator) Name() string {
	return g.typeName
}

func (g *generator) NameLower() string {
	return strings.ToLower(g.typeName)
}

func (g *generator) Package() string {
	return g.pkg.Name()
}

func (g *generator) PackagePath() string {
	return g.pkg.Path()
}

func (g *generator) IsInterface() bool {
	return g.isInterface
}

func (g *generator) Methods() []Method {
	if !g.isInterface {
		return nil
	}
	numMethods := g.iface.NumMethods()
	methods := make([]Method, numMethods)
	for i := 0; i < numMethods; i++ {
		//methods[i] = generate.Method{g, g.iface.Method(i)}
		methods[i] = *NewMethod(g.iface.Method(i))
	}
	return methods
}
