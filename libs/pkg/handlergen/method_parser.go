package handlergen

import (
	"errors"
	"fmt"
	"strings"
)

func findParamByName(params []param, name string) (param, bool) {
	for _, p := range params {
		if p.Name == name {
			return p, true
		}
	}
	return param{}, false
}

type requestMappingDefine struct {
	Path   string
	Method string
}

type paramDefine struct {
	Name     string
	Param    param
	Required bool
}

type MethodParser struct {
	MethodName     string
	StructName     string
	InterfaceName  string
	FuncDefine     string
	Params         []param
	Results        []param
	Doc            string
	RequestMapping *requestMappingDefine
	RequestParams  []paramDefine
	PathVariables  []paramDefine
	URIBinding     param
	BodyBinding    param
	ParamBinding   param
}

func (m *MethodParser) HasResponseData() bool {
	return len(m.Results) == 2
}

func (m *MethodParser) HasRequestMapping() bool {
	return m.RequestMapping != nil
}

func (m *MethodParser) HasURIBinding() bool {
	return !m.URIBinding.IsNull()
}

func (m *MethodParser) HasBodyBinding() bool {
	return !m.BodyBinding.IsNull()
}

func (m *MethodParser) HasParamBinding() bool {
	return !m.ParamBinding.IsNull()
}

func (m *MethodParser) HasResultError() bool {
	for _, p := range m.Results {
		if p.Type == "error" {
			return true
		}
	}
	return false
}

func (m *MethodParser) ResultErrorName() string {
	for _, p := range m.Results {
		if p.Type == "error" {
			if len(p.Name) > 0 {
				return p.Name
			}
			return "err"
		}
	}
	return "err"
}

func (m *MethodParser) GetParamInTmpl() string {
	return paramToString(m.Params)
}

func (m *MethodParser) GetResultsInTmpl() string {
	return paramToString(m.Results)
}

func (m *MethodParser) GetParamInFunc() string {
	var names []string
	for _, p := range m.Params {
		names = append(names, p.Name)
	}
	return strings.Join(names, ",")
}

func paramToString(params []param) string {
	var res []string
	for _, param := range params {
		name := param.Name
		tmplString := fmt.Sprintf("%s ", name)
		if param.IsArray {
			tmplString += "[]"
		}
		if param.IsPointer {
			tmplString += "*"
		}
		if param.Package != "" {
			tmplString += fmt.Sprintf("%s.", param.Package)
		}
		tmplString += param.Type
		res = append(res, tmplString)
	}
	return strings.Join(res, ",")
}

func (m *MethodParser) getDocString() string {
	docString := strings.TrimSpace(m.Doc)

	if index := strings.Index(docString, "\n\n"); index != -1 {
		if strings.Contains(docString[index+2:], m.MethodName) {
			docString = docString[:index]
		} else {
			docString = docString[index+2:]
		}
	}

	docString = strings.TrimPrefix(docString, m.MethodName)
	return docString
}

type kv struct {
	key string
	val string
}

func parseValue(data string) []kv {
	sp := strings.Split(data, ",")
	var out []kv
	for _, s := range sp {
		s = strings.Trim(s, " ")
		res := strings.Split(s, "=")
		if len(res) == 2 {
			out = append(out, kv{
				key: strings.Trim(res[0], " "),
				val: strings.Trim(res[1], " "),
			})
		}
	}
	return out
}

func (m *MethodParser) parseDoc() error {
	docString := strings.TrimSpace(m.getDocString())
	lines := strings.Split(strings.ReplaceAll(docString, "\n\r", "\n"), "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		key, value, ok := parseAnnotation(line)
		if !ok {
			continue
		}
		switch key {
		case "RequestMapping":
			path, method, _ := parseRequestMapping(line)
			m.RequestMapping = &requestMappingDefine{
				Path:   path,
				Method: method,
			}
		case "RequestParam":
			kvs := parseValue(value)
			for _, kv := range kvs {
				v := kv.val
				if strings.HasPrefix(kv.val, "@") {
					v = v[1:]
				}
				p, ok := findParamByName(m.Params, v)
				if !ok {
					return errors.New(fmt.Sprintf("param %s not found in method", v))
				}
				m.RequestParams = append(m.RequestParams, paramDefine{
					Name:  kv.key,
					Param: p,
				})
			}
		case "RequestParamRequired":
			kvs := parseValue(value)
			for _, kv := range kvs {
				v := kv.val
				if strings.HasPrefix(kv.val, "@") {
					v = v[1:]
				}
				p, ok := findParamByName(m.Params, v)
				if !ok {
					return errors.New(fmt.Sprintf("param %s not found in method", v))
				}
				m.RequestParams = append(m.RequestParams, paramDefine{
					Name:     kv.key,
					Param:    p,
					Required: true,
				})
			}

		case "PathVariable":
			kvs := parseValue(value)
			for _, kv := range kvs {
				v := kv.val
				if strings.HasPrefix(kv.val, "@") {
					v = v[1:]
				}
				p, ok := findParamByName(m.Params, v)
				if !ok {
					return errors.New(fmt.Sprintf("param %s not found in method", v))
				}
				m.PathVariables = append(m.RequestParams, paramDefine{
					Name:  kv.key,
					Param: p,
				})
			}
		case "BindURI":
			p, ok := findParamByName(m.Params, value)
			if !ok {
				return errors.New(fmt.Sprintf("param %s not found in method", value))
			}
			m.URIBinding = p
		case "BindQuery":
			p, ok := findParamByName(m.Params, value)
			if !ok {
				return errors.New(fmt.Sprintf("param %s not found in method", value))
			}
			m.ParamBinding = p
		case "BindBody":
			p, ok := findParamByName(m.Params, value)
			if !ok {
				return errors.New(fmt.Sprintf("param %s not found in method", value))
			}
			m.BodyBinding = p
		default:
			return errors.New("invalid annotation")
		}
	}
	return nil
}

func (m *MethodParser) Parse() error {
	err := m.parseDoc()
	if err != nil {
		return err
	}
	if m.HasRequestMapping() {
		if !(len(m.Results) == 2 && m.Results[1].IsError()) && !(len(m.Results) == 1 && m.Results[0].IsError()) {
			return errors.New("method result type invalid")
		}
	}
	return nil
}
