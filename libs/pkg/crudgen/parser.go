package crudgen

import (
	"fmt"
	"go/ast"
	"go/doc"
	"go/parser"
	"go/token"
	"io/ioutil"
	"log"
	"regexp"
	"strings"
	"unicode"
)

type Parser struct {
	visitor *interfaceVisitor
}

func toSnakeStyle(in string) string {
	beforeUpper := true
	var builder strings.Builder
	for _, r := range in {
		if unicode.IsUpper(r) {
			if !beforeUpper {
				builder.WriteByte('_')
				beforeUpper = true
			}
		} else {
			beforeUpper = false
		}
		builder.WriteRune(unicode.ToLower(r))
	}
	return builder.String()
}

type interfaceVisitor struct {
	data    []byte
	defines []interfaceDefine
	docs    map[string]string
}

type param struct { // (user model.User)
	Package   string // package's name: model
	Name      string // param's name: user
	Type      string // param's type: User
	IsArray   bool   // is array or not
	IsPointer bool   // is pointer or not
}

func (p *param) FullType() string {
	if len(p.Package) != 0 {
		return p.Package + "." + p.Type
	}
	return p.Type
}

func (p *param) IsMap() bool {
	return strings.HasPrefix(p.Type, "map[")
}

func (p *param) IsInterface() bool {
	return p.Type == "interface{}"
}

func (p *param) IsNull() bool {
	return p.Package == "" && p.Type == "" && p.Name == ""
}

func (p *param) IsError() bool {
	return p.Type == "error"
}

func (p *param) IsTime() bool {
	return p.Package == "time" && p.Type == "Time"
}

func (p *param) astGetParamType(param ast.Expr) {
	switch v := param.(type) {
	case *ast.Ident:
		p.Type = v.Name
		if v.Obj != nil {
			p.Package = "UNDEFINED" // set a placeholder
		}
	case *ast.SelectorExpr:
		p.astGetEltType(v)
	case *ast.ArrayType:
		p.astGetEltType(v.Elt)
		p.IsArray = true
	case *ast.Ellipsis:
		p.astGetEltType(v.Elt)
		p.IsArray = true
	case *ast.MapType:
		p.astGetMapType(v)
	case *ast.InterfaceType:
		p.Type = "interface{}"
	case *ast.StarExpr:
		p.IsPointer = true
		p.astGetEltType(v.X)
	default:
		log.Fatalf("unknow param type: %+v", v)
	}
}

func (p *param) astGetEltType(expr ast.Expr) {
	switch v := expr.(type) {
	case *ast.Ident:
		p.Type = v.Name
		if v.Obj != nil {
			p.Package = "UNDEFINED"
		}
	case *ast.SelectorExpr:
		p.Type = v.Sel.Name
		p.astGetPackageName(v.X)
	case *ast.MapType:
		p.astGetMapType(v)
	case *ast.StarExpr:
		p.IsPointer = true
		p.astGetEltType(v.X)
	case *ast.InterfaceType:
		p.Type = "interface{}"
	default:
		log.Fatalf("unknow param type: %+v", v)
	}
}

func (p *param) astGetPackageName(expr ast.Expr) {
	switch v := expr.(type) {
	case *ast.Ident:
		p.Package = v.Name
	}
}

func (p *param) astGetMapType(expr *ast.MapType) string {
	p.Type = fmt.Sprintf("map[%s]%s", astGetType(expr.Key), astGetType(expr.Value))
	return ""
}

func astGetType(expr ast.Expr) string {
	switch v := expr.(type) {
	case *ast.Ident:
		return v.Name
	case *ast.InterfaceType:
		return "interface{}"
	}
	return ""

}

type methodDefine struct {
	Name    string
	Define  string
	Params  []param
	Results []param
	Doc     string
}

type interfaceDefine struct {
	Name      string
	Methods   []methodDefine
	Model     param
	TableName string
}

var annotationRegexp = regexp.MustCompile(`^@([a-z A-Z]+)\((.*)\)$`)

func parseAnnotation(s string) (key string, value string, ok bool) {
	matches := annotationRegexp.FindStringSubmatch(s)
	if len(matches) == 0 {
		return "", "", false
	}
	return matches[1], matches[2], true
}

func parseTableName(in string) string {
	lines := strings.Split(strings.ReplaceAll(in, "\n\r", "\n"), "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		key, value, ok := parseAnnotation(line)
		if ok && key == "Table" {
			return value
		}
	}
	return ""
}

func getParamList(fields *ast.FieldList) []param {
	if fields == nil {
		return nil
	}
	var pars []param
	if len(fields.List) < 1 {
		return nil
	}
	for _, field := range fields.List {
		if field.Names == nil {
			par := param{}
			par.astGetParamType(field.Type)
			pars = append(pars, par)
			continue
		}

		for _, name := range field.Names {
			par := param{
				Name: name.Name,
			}
			par.astGetParamType(field.Type)
			pars = append(pars, par)
			continue
		}
	}
	return pars
}

func (i *interfaceVisitor) Visit(n ast.Node) (w ast.Visitor) {
	// if n != nil {
	// 	fmt.Println(reflect.TypeOf(n).String(), n)
	// }
	switch n := n.(type) {
	case *ast.TypeSpec:
		define := interfaceDefine{}
		if data, ok := n.Type.(*ast.InterfaceType); ok {
			interfaceDoc := i.docs[n.Name.Name]
			define.TableName = parseTableName(interfaceDoc)
			if len(define.TableName) == 0 {
				// define.TableName = toSnakeStyle(n.Name.Name)
				return
			}
			define.Name = n.Name.Name
			methods := data.Methods.List
			for _, method := range methods {
				if len(method.Names) == 0 {
					if v, ok := method.Type.(*ast.IndexListExpr); ok {
						p := param{}
						p.astGetParamType(v.X)
						if p.Package == "helper" && p.Type == "DAO" {
							p.astGetParamType(v.Indices[1])
							define.Model = p
						}
					}
					continue
				}
				for _, name := range method.Names {
					var mdefine methodDefine
					mdefine.Define = string(i.data[method.Pos()-1 : method.End()-1])
					mdefine.Name = name.Name
					fmt.Println(name.Name)
					mdefine.Doc = method.Doc.Text()
					mdefine.Params = getParamList(method.Type.(*ast.FuncType).Params)
					mdefine.Results = getParamList(method.Type.(*ast.FuncType).Results)
					define.Methods = append(define.Methods, mdefine)
				}
			}
			i.defines = append(i.defines, define)
		}
	}
	return i
}

func (p *Parser) ParseFile(name string) {
	data, err := ioutil.ReadFile(name)
	if err != nil {
		panic(err)
	}
	fset := token.NewFileSet()

	f, err := parser.ParseFile(fset, name, nil, parser.ParseComments)
	if err != nil {
		panic(err)
	}

	doc, err := doc.NewFromFiles(fset, []*ast.File{f}, "main")
	if err != nil {
		panic(err)
	}
	m := make(map[string]string)
	for _, t := range doc.Types {
		m[t.Name] = t.Doc
	}

	v := &interfaceVisitor{data: data, docs: m}
	ast.Walk(v, f)
	p.visitor = v
}
