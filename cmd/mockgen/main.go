package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"path"
	"path/filepath"

	"github.com/pgt502/gogen/mockgen"
)

//import "github.com/josharian/impl"

var (
	pkgName = flag.String("pkg", ".", "package name of the interface to generate the mock from")
	output  = flag.String("o", ".", "output folder")
)

func main() {
	flag.Parse()

	//abs := "../auth"
	//ifaceName := "AuthMiddleware"

	ifaceName := flag.Arg(0)

	//*pkgName = "git.uk.guardtime.com/guardtime/volta/voltapm/messaging/common"
	//ifaceName = "RtbfHandlerFactory"

	fmt.Printf("pkg is: %s, interface: %s\n", *pkgName, ifaceName)
	if *pkgName == "" {
		fmt.Println("pkg name needs to be specified")
		return
	}
	if ifaceName == "" {
		fmt.Println("interface needs to be provided as an argument")
		return
	}
	g, err := mockgen.NewGenerator(ifaceName, *pkgName)
	if err != nil {
		fmt.Println(err)
		return
	}

	content := g.Generate()
	fmt.Printf("generated file: %s\n", content)
	fn := path.Join(*output, fmt.Sprintf("Mock%s.go", ifaceName))
	fn, err = filepath.Abs(fn)
	if err != nil {
		fmt.Println(err)
		return
	}
	writeToFile(fn, content)
}

func writeToFile(name, content string) {
	err := ioutil.WriteFile(name, []byte(content), 0644)
	if err != nil {
		fmt.Println(err)
		return
	}
}
