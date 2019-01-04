package generate

import (
	"fmt"
	"go/format"
	"io/ioutil"
	"log"
	"path"
	"path/filepath"
	"strings"
)

func GenerateFileFromTemplate(tplPath, dstFile string, g BasicGenerator) {
	tpl, err := readTemplate(tplPath)
	if err != nil {
		log.Println(err)
		return
	}
	content := g.Generate(string(tpl))
	if strings.HasSuffix(dstFile, ".go") {
		formatted, err := format.Source([]byte(content))
		if err != nil {
			log.Printf("error formatting the source code, err: %s", err)
			return
		}
		content = string(formatted)
	}
	log.Printf("generated file: %s\n", content)

	dstFile, err = filepath.Abs(dstFile)
	if err != nil {
		log.Println(err)
		return
	}
	writeToFile(dstFile, content)
}

func GenerateFilesFromTemplates(ifaceName, tplPath, dstDir string, g BasicGenerator) {
	fis, err := ioutil.ReadDir(tplPath)
	if err != nil {
		log.Printf("error reading files from dir: %s, err: %s", tplPath, err)
		return
	}
	var templates []string
	for _, fi := range fis {
		if fi.IsDir() {
			continue
		}
		templates = append(templates, fi.Name())
	}

	for _, el := range templates {
		var fn string

		ind := strings.Index(el, "%s")
		switch ind {
		case 0:
			fn = fmt.Sprintf(el, ifaceName)
		case -1:
			fn = el
		default:
			fn = fmt.Sprintf(el, strings.Title(ifaceName))
		}
		dstFile := path.Join(dstDir, fn)
		dstFile = strings.TrimSuffix(dstFile, ".tpl")
		GenerateFileFromTemplate(path.Join(tplPath, el), dstFile, g)
	}
}

func readTemplate(fn string) (content []byte, err error) {
	fn, err = filepath.Abs(fn)
	if err != nil {
		log.Println(err)
		return
	}
	bytes, err := ioutil.ReadFile(fn)
	if err != nil {
		log.Println(err)
		return
	}
	content = bytes
	return
}

func writeToFile(name, content string) {
	err := ioutil.WriteFile(name, []byte(content), 0644)
	if err != nil {
		log.Println(err)
		return
	}
}
