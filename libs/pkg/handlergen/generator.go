package handlergen

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"os"
	"path"
	"strings"
	"text/template"

	"github.com/pkg/errors"
	"golang.org/x/tools/imports"
)

type Config struct {
	Package          string
	ModelPackage     string
	InterfacePackage string
	HelperPackage    string
	OutputPath       string
	Imports          []string
}

type Generator struct {
}

func execute(t *template.Template, buf io.Writer, value interface{}) error {
	// t := template.Must(template.New("header").Parse(ImportHeaderTemplate))
	return t.Execute(buf, value)
}

func errorLog(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

type paramWrapper struct {
	param
	ContainerTag string
}

type structTmpl struct {
	StructName           string
	ContainerTag         string
	InterfaceName        string
	InterfacePackageTail string
	Model                *param
	Services             []paramWrapper
	RootPath             string
	Parsers              []MethodParser
}

func (g *Generator) Generate(p *Parser, conf Config) error {

	notEdit := template.Must(template.New("notEditMark").Parse(NotEditMarkTemplate))
	importHeader := template.Must(template.New("importHeader").Parse(ImportHeaderTemplate))
	structDefine := template.Must(template.New("structDefine").Parse(StructDefineTemplate))
	impStructDefine := template.Must(template.New("impStructDefine").Parse(ImpStructDefineTemplate))
	wrapperFunc := template.Must(template.New("wrapperFunc").Parse(WrapperFuncTemplate))
	impFunc := template.Must(template.New("impFunc").Parse(ImpFuncTemplate))
	conf.Imports = p.visitor.imports

	for _, idefine := range p.visitor.defines {
		buf := bytes.NewBuffer(nil)
		impBuf := bytes.NewBuffer(nil)
		errorLog(execute(notEdit, buf, ""))
		errorLog(execute(importHeader, buf, conf))
		errorLog(execute(importHeader, impBuf, conf))

		wrapperFile := toSnakeStyle(idefine.Name+"Wrapper") + ".gen.go"
		impFile := toSnakeStyle(idefine.Name+"Func") + ".example"
		splits := strings.Split(conf.InterfacePackage, "/")
		var parsers []MethodParser
		for _, m := range idefine.Methods {
			mp := MethodParser{
				MethodName:    m.Name,
				InterfaceName: idefine.Name,
				Params:        m.Params,
				Results:       m.Results,
				Doc:           m.Doc,
				FuncDefine:    m.Define,
			}
			err := mp.Parse()
			if err != nil {
				return err
			}
			parsers = append(parsers, mp)
		}
		sTmpl := structTmpl{
			ContainerTag:         "`container:\"type\"`",
			InterfaceName:        idefine.Name,
			InterfacePackageTail: splits[len(splits)-1],
			Model:                &idefine.Model,
			RootPath:             idefine.RootPath,
			Parsers:              parsers,
		}
		for _, s := range idefine.GetServices() {
			sTmpl.Services = append(sTmpl.Services, paramWrapper{
				param:        s,
				ContainerTag: "`container:\"type\"`",
			})
		}
		errorLog(execute(structDefine, buf, sTmpl))
		errorLog(execute(impStructDefine, buf, sTmpl))

		for _, mp := range parsers {
			if mp.HasRequestMapping() {
				errorLog(wrapperFunc.Execute(buf, &mp))
				errorLog(impFunc.Execute(impBuf, &mp))
			}
		}

		fullPath := path.Join(conf.OutputPath, wrapperFile)
		if err := save(fullPath, buf.Bytes()); err != nil {
			return err
		}
		fullPath = path.Join(conf.OutputPath, impFile)
		if err := save(fullPath, impBuf.Bytes()); err != nil {
			return err
		}
	}

	return nil
}

func save(file string, data []byte) error {
	f, err := os.OpenFile(file,
		os.O_CREATE|os.O_RDWR|os.O_TRUNC, os.ModePerm)
	if err != nil {
		return errors.Wrapf(err, "open file: %s", file)
	}
	res, err := imports.Process("", data, &imports.Options{
		TabWidth:  8,
		TabIndent: true,
		Comments:  true,
		Fragment:  true,
	})
	if err != nil {
		fmt.Println(string(data))
		return errors.Wrap(err, "import fmt")
	}
	f.Write(res)
	f.Close()
	return nil
}
