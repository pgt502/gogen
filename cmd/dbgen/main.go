package main

import (
	"flag"
	"fmt"
	"log"
	"path"
	"strings"

	"github.com/pgt502/gogen/generate"

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

	ifaceName := flag.Arg(0)

	//*pkgName = "../../testdata"
	//ifaceName = "Order"

	log.Printf("pkg is: %s, interface: %s\n", *pkgName, ifaceName)
	if *pkgName == "" {
		fmt.Println("pkg name needs to be specified")
		return
	}
	if ifaceName == "" {
		log.Println("struct needs to be provided as an argument")
		return
	}
	g, err := dbgen.NewGenerator(ifaceName, *pkgName)
	if err != nil {
		log.Println(err)
		return
	}
	if g.IsInterface() {
		log.Println("struct needs to be provided as an argument - got interface")
		return
	}

	// generate interface:
	generate.GenerateFile("./templates/tableInterface.tpl", path.Join(*output, fmt.Sprintf("%sTable.go", strings.ToLower(ifaceName))), g)
	// generate pgTable:
	generate.GenerateFile("./templates/pgTableStruct.tpl", path.Join(*output, fmt.Sprintf("pg%sTable.go", ifaceName)), g)
	// generate postgresql script
	generate.GenerateFile("./templates/tableScript.up.tpl", path.Join(*output, fmt.Sprintf("nn_table_%s.up.sql", strings.ToLower(ifaceName))), g)
	generate.GenerateFile("./templates/tableScript.down.tpl", path.Join(*output, fmt.Sprintf("nn_table_%s.down.sql", strings.ToLower(ifaceName))), g)
	// generate repo:
	generate.GenerateFile("./templates/tableRepo.tpl", path.Join(*output, fmt.Sprintf("%sRepo.go", strings.ToLower(ifaceName))), g)
}
