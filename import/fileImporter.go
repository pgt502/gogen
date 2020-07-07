package generate

import (
	"go/ast"
	"go/importer"
	"go/parser"
	"go/token"
	"go/types"
	"io/ioutil"
)

type fileImporter struct {
	base     types.Importer
	pkg      *types.Package
	typeName string
	iface    *types.Interface
	strct    *types.Struct
	imported map[string]*types.Package
}

// NewImporter is a constructor for file importer
func NewImporter() FileImporter {
	return &fileImporter{
		imported: make(map[string]*types.Package),
		base:     importer.Default(),
	}
}

// FileImporter allows to import a package from file
type FileImporter interface {
	types.Importer
	ImportFile(fn string) (pkg *types.Package, err error)
}

func (i *fileImporter) Import(path string) (pkg *types.Package, err error) {
	if pkg, ok := i.imported[path]; ok {
		return pkg, nil
	}
	pkg, err = importOrErr(i.base, path, err)
	i.imported[path] = pkg
	return
}

func (i *fileImporter) ImportFile(fn string) (pkg *types.Package, err error) {
	fset := token.NewFileSet()
	src, err := ioutil.ReadFile(fn)
	if err != nil {
		return nil, err
	}
	f, err := parser.ParseFile(fset, fn, src, 0)
	if err != nil {
		return nil, err
	}
	packageName := f.Name.Name

	conf := types.Config{
		Importer: i,
	}
	pkg, err = conf.Check(packageName, fset, []*ast.File{f}, nil)

	return
}

func importOrErr(base types.Importer, pkg string, err error) (*types.Package, error) {
	p, impErr := base.Import(pkg)
	if impErr != nil {
		return nil, err
	}
	return p, nil
}
