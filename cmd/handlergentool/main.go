package main

import (
	"flag"
	"log"

	"github.com/LSDXXX/libs/pkg/handlergen"
)

var (
	fileName    string
	outputPath  string
	packageName string
)

func init() {
	flag.StringVar(&fileName, "f", "", "input file")
	flag.StringVar(&outputPath, "op", "", "output path")
	flag.StringVar(&packageName, "pkg", "", "package name")
}

func main() {
	flag.Parse()
	if len(fileName) == 0 || len(outputPath) == 0 {
		log.Fatalf("input file or outputPath is empty")
	}
	p := handlergen.Parser{}
	p.ParseFile(fileName)
	g := handlergen.Generator{}
	err := g.Generate(&p, handlergen.Config{
		Package:       packageName,
		HelperPackage: "github.com/LSDXXX/libs/pkg/handlergen/helper",
		OutputPath:    outputPath,
	})
	if err != nil {
		panic(err)
	}
}
