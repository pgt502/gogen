package generate

import (
	"go/format"
	"io/ioutil"
	"log"
	"path/filepath"
	"strings"
)

func GenerateFile(tplPath, dstFile string, g BasicGenerator) {
	tpl, err := readTemplate(tplPath)
	if err != nil {
		log.Println(err)
		return
	}
	content := g.Generate(tpl)
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

func readTemplate(fn string) (content string, err error) {
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
	content = string(bytes)
	return
}

func writeToFile(name, content string) {
	err := ioutil.WriteFile(name, []byte(content), 0644)
	if err != nil {
		log.Println(err)
		return
	}
}
