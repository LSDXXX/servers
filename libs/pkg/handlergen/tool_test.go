package handlergen

import "testing"

func TestG(t *testing.T) {

	p := Parser{}
	p.ParseFile("test.go")
	g := Generator{}
	err := g.Generate(&p, Config{
		Package:          "infra",
		ModelPackage:     "github.com/LSDXXX//libs/model",
		InterfacePackage: "github.com/LSDXXX//libs/repo",
		HelperPackage:    "github.com/LSDXXX//libs/pkg/handlergen/helper",
		OutputPath:       "./",
	})
	if err != nil {
		panic(err)
	}
}
