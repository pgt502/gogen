package main

import (
	"flag"
	"fmt"
	"log"
	"strings"

	"github.com/pgt502/gogen/generate"
	"github.com/pgt502/gogen/mockgen"
)

var (
	pkgName   = flag.String("pkg", ".", "package name of the interface to generate the mock from")
	output    = flag.String("o", ".", "output folder")
	templates = flag.String("t", "./templates", "templates folder")
)

func main() {
	flag.Parse()

	ifaceName := flag.Arg(0)

	//*pkgName = "../../auth"
	//ifaceName = "AuthMiddleware"

	fmt.Printf("pkg is: %s, interface: %s\n", *pkgName, ifaceName)
	if *pkgName == "" {
		log.Println("pkg name needs to be specified")
		return
	}
	if ifaceName == "" {
		log.Println("interface needs to be provided as an argument")
		return
	}
	g, err := mockgen.NewGenerator(ifaceName, *pkgName)
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
