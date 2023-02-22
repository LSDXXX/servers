package crudgen

import (
	"fmt"
	"strings"
	"unicode"
)

type fragment struct {
	Type    StatementType
	Value   string
	IsArray bool
}

func isDigit(str string) bool {
	for _, x := range str {
		if !unicode.IsDigit(x) {
			return false
		}
	}
	return true
}

func (f *fragment) fragmentByParams(params []param) {
	for _, param := range params {
		if param.Name == f.Value {
			f.IsArray = param.IsArray
			switch param.Type {
			case "bool":
				f.Type = BOOL
				return
			case "int":
				f.Type = INT
				return
			case "string":
				f.Type = STRING
				return
			case "Time":
				f.Type = TIME
			default:
				f.Type = OTHER
			}
		} else if param.Name == strings.Split(f.Value, ".")[0] {
			f.Type = UNCERTAIN
			f.IsArray = false
		}
	}
}

func checkFragment(s string, params []param) (f fragment, err error) {
	digital := func(str string) string {
		if isDigit(str) {
			return "<integer>"
		}
		return str
	}

	f = fragment{Type: UNKNOWN, Value: strings.Trim(s, " ")}
	str := strings.ToLower(strings.Trim(s, " "))
	switch digital(str) {
	case "<integer>":
		f.Type = INT
	case "&&", "||":
		f.Type = LOGICAL
	case ">", "<", ">=", "<=", "==", "!=":
		f.Type = EXPRESSION
	case "end":
		f.Type = END
	case "if":
		f.Type = IF
	case "set":
		f.Type = SET
	case "else":
		f.Type = ELSE
	case "where":
		f.Type = WHERE
	case "true", "false":
		f.Type = BOOL
	case "nil":
		f.Type = NIL
	default:
		f.fragmentByParams(params)
		if f.Type == UNKNOWN {
			err = fmt.Errorf("unknow parameter: %s", s)
		}
	}
	return
}

func checkTemplate(tmpl string, params []param) (result statement, err error) {
	fragmentList, err := splitTemplate(tmpl, params)
	if err != nil {
		return
	}
	err = checkTempleFragmentValid(fragmentList)
	if err != nil {
		return
	}
	return fragmentToSLice(fragmentList)
}

func fragmentToSLice(list []fragment) (part statement, err error) {
	var values []string

	if len(list) == 0 {
		return
	}
	for _, t := range list {
		values = append(values, t.Value)
	}
	part.Origin = strings.Join(values, " ")
	switch strings.ToLower(values[0]) {
	case "if":
		if len(values) > 1 {
			part.Type = IF
			part.Value = strings.Join(values[1:], " ")
			return
		}
	case "else":
		if len(values) == 1 {
			part.Type = ELSE
			return
		} else {
			if strings.ToLower(values[1]) == "if" && len(values) > 2 {
				part.Value = strings.Join(values[2:], " ")
				part.Type = ELSEIF
				return
			}
		}
	case "where":
		part.Type = WHERE
		return
	case "set":
		part.Type = SET
		return
	case "end":
		part.Type = END
		return
	}

	err = fmt.Errorf("syntax error:%s", strings.Join(values, " "))
	return
}

func splitTemplate(tmpl string, params []param) (fragList []fragment, err error) {
	var buf SQLBuffer
	var f fragment
	for i := 0; !strOutrange(i, tmpl); i++ {
		switch tmpl[i] {
		case '"':
			_ = buf.WriteByte(tmpl[i])
			for i++; ; i++ {
				if strOutrange(i, tmpl) {
					return nil, fmt.Errorf("incomplete code:%s", tmpl)
				}
				_ = buf.WriteByte(tmpl[i])

				if tmpl[i] == '"' && tmpl[i-1] != '\\' {
					fragList = append(fragList, fragment{Type: STRING, Value: buf.Dump()})
					break
				}
			}
		case ' ':
			if sqlClause := buf.Dump(); sqlClause != "" {
				f, err = checkFragment(sqlClause, params)
				if err != nil {
					return nil, err
				}
				fragList = append(fragList, f)
			}
		case '>', '<', '=', '!':
			if sqlClause := buf.Dump(); sqlClause != "" {
				f, err = checkFragment(sqlClause, params)
				if err != nil {
					return nil, err
				}
				fragList = append(fragList, f)
			}

			_ = buf.WriteByte(tmpl[i])

			if strOutrange(i+1, tmpl) {
				return nil, fmt.Errorf("incomplete code:%s", tmpl)
			}
			if tmpl[i+1] == '=' {
				_ = buf.WriteByte(tmpl[i+1])
				i++
			}

			f, err = checkFragment(buf.Dump(), params)
			if err != nil {
				return nil, err
			}
			fragList = append(fragList, f)
		case '&', '|':
			if strOutrange(i+1, tmpl) {
				return nil, fmt.Errorf("incomplete code:%s", tmpl)
			}

			if tmpl[i+1] == tmpl[i] {
				i++

				if sqlClause := buf.Dump(); sqlClause != "" {
					f, err = checkFragment(sqlClause, params)
					if err != nil {
						return nil, err
					}
					fragList = append(fragList, f)
				}

				// write && or ||
				fragList = append(fragList, fragment{
					Type:  LOGICAL,
					Value: tmpl[i-1 : i+1],
				})
			}
		default:
			_ = buf.WriteByte(tmpl[i])
		}
	}

	if sqlClause := buf.Dump(); sqlClause != "" {
		f, err = checkFragment(sqlClause, params)
		if err != nil {
			return nil, err
		}
		fragList = append(fragList, f)
	}
	return fragList, nil
}

func checkTempleFragmentValid(list []fragment) error {
	for i := 1; i < len(list); i++ {
		switch list[i].Type {
		case IF, ELSE, END, BOOL, LOGICAL, WHERE, SET:
			continue
		case INT, STRING, OTHER, UNCERTAIN, TIME, NIL:
			if i+2 < len(list) {
				if isExpressionValid(list[i : i+3]) {
					i += 2
				} else {
					return fmt.Errorf("condition type not match：%s", fragmentToString(list[i:i+3]))
				}
			}
		default:
			return fmt.Errorf("unknow fragment ： %s ", list[i].Value)
		}
	}
	return nil
}

func isExpressionValid(expr []fragment) bool {
	if len(expr) != 3 {
		return false
	}
	if expr[1].Type != EXPRESSION {
		return false
	}
	//Only arrays can be compared with nil
	if expr[0].Type == NIL || expr[2].Type == NIL {
		return expr[0].IsArray || expr[2].IsArray || expr[0].Type == UNCERTAIN || expr[2].Type == UNCERTAIN
	}

	return expr[0].Type == expr[2].Type ||
		expr[0].Type == UNCERTAIN || expr[2].Type == UNCERTAIN
}

func fragmentToString(list []fragment) string {
	var values []string

	if len(list) == 0 {
		return ""
	}
	for _, t := range list {
		values = append(values, t.Value)
	}
	return strings.Join(values, " ")
}
