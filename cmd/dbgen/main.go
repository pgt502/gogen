package main

import (
	"flag"
	"fmt"
	"log"
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
	pkgName   = flag.String("pkg", "", "package name of the type to generate the code from")
	file      = flag.String("f", "", "path to the file where the type to generate the code from is define")
	output    = flag.String("o", ".", "output folder")
	templates = flag.String("t", "./templates", "templates folder")
	pluralise = flag.Bool("p", false, "pluralise the table name")
)

func main() {
	flag.Parse()

	flag.Var(&myFlags, "list", "list of params")

	ifaceName := flag.Arg(0)

	//*pkgName = "../../testdata"
	//ifaceName = "Order"

	log.Printf("pkg is: %s, file: %s, type: %s\n", *pkgName, *file, ifaceName)
	if *pkgName == "" && *file == "" {
		fmt.Println("pkg name or file needs to be specified")
		return
	}
	if ifaceName == "" {
		log.Println("struct needs to be provided as an argument")
		return
	}
	g, err := dbgen.NewGenerator(ifaceName, dbgen.Package(pkgName),
		dbgen.File(file), dbgen.Pluralise(*pluralise))
	if err != nil {
		log.Println(err)
		return
	}
	if g.IsInterface() {
		log.Println("struct needs to be provided as an argument - got interface")
		return
	}

	generate.GenerateFilesFromTemplates(strings.ToLower(ifaceName), *templates, *output, g)
}
