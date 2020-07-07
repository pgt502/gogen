package main

import (
	"flag"
	"log"
	"strings"

	"github.com/pgt502/gogen/generate"
	"github.com/pgt502/gogen/mockgen"
)

var (
	pkgName   = flag.String("pkg", "", "package name of the interface to generate the mock from")
	file      = flag.String("f", "", "path to the file where the type to generate the code from is define")
	output    = flag.String("o", ".", "output folder")
	templates = flag.String("t", "./templates/repo", "templates folder")
)

func main() {
	flag.Parse()

	ifaceName := flag.Arg(0)

	//*pkgName = "../../auth"
	//ifaceName = "AuthMiddleware"

	log.Printf("pkg is: %s, file: %s, interface: %s\n", *pkgName, *file, ifaceName)
	if *pkgName == "" && *file == "" {
		log.Println("pkg name or file needs to be specified")
		return
	}
	if ifaceName == "" {
		log.Println("interface needs to be provided as an argument")
		return
	}
	g, err := mockgen.NewGenerator(ifaceName, *pkgName, *file)
	if err != nil {
		log.Println(err)
		return
	}
	if !g.IsInterface() {
		log.Printf("provided type: [%s] is not an interface", ifaceName)
		return
	}
	generate.GenerateFilesFromTemplates(strings.ToLower(ifaceName), *templates, *output, g)
}
