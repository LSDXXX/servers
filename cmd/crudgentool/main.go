package main

import (
	"flag"
	"log"

	"github.com/LSDXXX/libs/pkg/crudgen"
)

var (
	fileName   string
	outputPath string
)

func init() {
	flag.StringVar(&fileName, "f", "", "input file")
	flag.StringVar(&outputPath, "op", "", "output path")
}

func main() {
	flag.Parse()
	if len(fileName) == 0 || len(outputPath) == 0 {
		log.Fatalf("input file or outputPath is empty")
	}
	p := crudgen.Parser{}
	p.ParseFile(fileName)
	g := crudgen.Generator{}
	err := g.Generate(&p, crudgen.Config{
		Package:          "infra",
		ModelPackage:     "github.com/LSDXXX/libs/model",
		InterfacePackage: "github.com/LSDXXX/libs/repo",
		HelperPackage:    "github.com/LSDXXX/libs/pkg/crudgen/helper",
		OutputPath:       outputPath,
	})
	if err != nil {
		panic(err)
	}
}
