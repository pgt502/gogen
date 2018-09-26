package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"path"
	"path/filepath"
	"strings"

	"github.com/pgt502/gogen/dbgen"
)

type arrayFlags []string

func (i *arrayFlags) String() string {
	return "my string representation"
}

func (i *arrayFlags) Set(value string) error {
	*i = append(*i, value)
	return nil
}

var myFlags arrayFlags

var (
	pkgName = flag.String("pkg", ".", "package name of the interface to generate the mock from")
	output  = flag.String("o", ".", "output folder")
)

func main() {
	flag.Parse()

	flag.Var(&myFlags, "list", "list of params")

	//abs := "../auth"
	//ifaceName := "AuthMiddleware"

	ifaceName := flag.Arg(0)

	*pkgName = "../../testdata"
	ifaceName = "Order"

	fmt.Printf("pkg is: %s, interface: %s\n", *pkgName, ifaceName)
	if *pkgName == "" {
		fmt.Println("pkg name needs to be specified")
		return
	}
	if ifaceName == "" {
		fmt.Println("interface needs to be provided as an argument")
		return
	}
	g, err := dbgen.NewGenerator(ifaceName, *pkgName)
	if err != nil {
		fmt.Println(err)
		return
	}

	// generate interface:
	generateFile("tableInterface.tpl", path.Join(*output, fmt.Sprintf("%sTable.go", strings.ToLower(ifaceName))), g)
	// generate pgTable:
	generateFile("pgTableStruct.tpl", path.Join(*output, fmt.Sprintf("pg%sTable.go", ifaceName)), g)
	// generate postgresql script
	generateFile("tableScript.up.tpl", path.Join(*output, fmt.Sprintf("nn_table_%s.up.sql", strings.ToLower(ifaceName))), g)
	generateFile("tableScript.down.tpl", path.Join(*output, fmt.Sprintf("nn_table_%s.down.sql", strings.ToLower(ifaceName))), g)
	// generate repo:
	generateFile("tableRepo.tpl", path.Join(*output, fmt.Sprintf("%sRepo.go", strings.ToLower(ifaceName))), g)
}

func generateFile(tplName, fn string, g *dbgen.Generator) {
	tpl, err := readTemplate(tplName)
	if err != nil {
		fmt.Println(err)
		return
	}
	content := g.Generate(tpl)
	fmt.Printf("generated file: %s\n", content)

	fn, err = filepath.Abs(fn)
	if err != nil {
		fmt.Println(err)
		return
	}
	writeToFile(fn, content)
}

func readTemplate(name string) (content string, err error) {
	fn := path.Join("./templates", name)
	fn, err = filepath.Abs(fn)
	if err != nil {
		fmt.Println(err)
		return
	}
	bytes, err := ioutil.ReadFile(fn)
	if err != nil {
		fmt.Println(err)
		return
	}
	content = string(bytes)
	return
}

func writeToFile(name, content string) {
	err := ioutil.WriteFile(name, []byte(content), 0644)
	if err != nil {
		fmt.Println(err)
		return
	}
}
