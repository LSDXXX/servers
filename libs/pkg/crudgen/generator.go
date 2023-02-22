package crudgen

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

type structTmpl struct {
	TableName            string
	StructName           string
	ContainerTag         string
	InterfaceName        string
	InterfacePackageTail string
	Model                *param
}

func (g *Generator) Generate(p *Parser, conf Config) error {

	notEdit := template.Must(template.New("notEditMark").Parse(NotEditMarkTemplate))
	importHeader := template.Must(template.New("importHeader").Parse(ImportHeaderTemplate))
	structDefine := template.Must(template.New("structDefine").Parse(StructDefineTemplate))

	for _, idefine := range p.visitor.defines {
		buf := bytes.NewBuffer(nil)
		errorLog(execute(notEdit, buf, ""))
		errorLog(execute(importHeader, buf, conf))
		structName := idefine.Name + "Imp"
		fileName := toSnakeStyle(structName) + ".go"
		splits := strings.Split(conf.InterfacePackage, "/")
		errorLog(execute(structDefine, buf, structTmpl{
			TableName:            idefine.TableName,
			StructName:           structName,
			ContainerTag:         "`container:\"type\"`",
			InterfaceName:        idefine.Name,
			InterfacePackageTail: splits[len(splits)-1],
			Model:                &idefine.Model,
		}))
		for _, m := range idefine.Methods {
			mp := MethodParser{
				MethodName: m.Name,
				StructName: structName,
				Params:     m.Params,
				Results:    m.Results,
				Doc:        m.Doc,
				Table:      idefine.TableName,
				FuncDefine: m.Define,
			}
			err := mp.Parse()
			if err != nil {
				return err
			}
			errorLog(mp.MethodTemplate.Execute(buf, &mp))
		}
		fullPath := path.Join(conf.OutputPath, fileName)
		f, err := os.OpenFile(fullPath,
			os.O_CREATE|os.O_RDWR|os.O_TRUNC, os.ModePerm)
		if err != nil {
			return errors.Wrapf(err, "open file: %s", fullPath)
		}
		data := buf.Bytes()
		fmt.Println(string(data))
		res, err := imports.Process("", data, &imports.Options{
			TabWidth:  8,
			TabIndent: true,
			Comments:  true,
			Fragment:  true,
		})
		if err != nil {
			return errors.Wrap(err, "import fmt")
		}
		f.Write(res)
		f.Close()
	}

	return nil
}
