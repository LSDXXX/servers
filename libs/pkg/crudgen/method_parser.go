package crudgen

import (
	"bytes"
	"errors"
	"fmt"
	"log"
	"strconv"
	"strings"
	"text/template"

	"github.com/LSDXXX/libs/pkg/crudgen/helper"
	"github.com/spf13/cast"
)

// Clause a symbol of clause, it can be sql condition clause, if clause, where clause, set clause and else cluase
type Clause interface {
	String() string
}

var (
	_ Clause = new(SQLClause)
	_ Clause = new(IfClause)
	_ Clause = new(ElseClause)
	_ Clause = new(WhereClause)
	_ Clause = new(SetClause)
)

var (
	userDefinedMethod = template.Must(template.New("userDefinedMethod").Parse(UserDefinedMethodTemplate))
)

type clause struct {
	VarName string
	Type    StatementType
}

// SQLClause sql condition clause
type SQLClause struct {
	clause
	Value []string
}

func (s SQLClause) String() string {
	return strings.ReplaceAll(strings.Join(s.Value, "+"), `"+"`, "")
}

// IfClause if clause
type IfClause struct {
	clause
	Cond  string
	Value []Clause
	Else  []Clause
}

func (i IfClause) String() string {
	return fmt.Sprintf("helper.IfClause(%s)", i.VarName)
}

// ElseClause else clause
type ElseClause struct {
	clause
	Cond  string
	Value []Clause
}

func (e ElseClause) String() (res string) {
	condList := make([]string, len(e.Value))
	for i, v := range e.Value {
		condList[i] = v.String()
	}
	return strings.ReplaceAll(strings.Join(condList, "+"), `"+"`, "")
}

// WhereClause where clause
type WhereClause struct {
	clause
	Value []Clause
}

func (w WhereClause) String() string {
	return fmt.Sprintf("helper.WhereClause(%s)", w.VarName)
}

// SetClause set clause
type SetClause struct {
	clause
	Value []Clause
}

func (w SetClause) String() string {
	return fmt.Sprintf("helper.SetClause(%s)", w.VarName)
}

type StatementType int

const (
	UNKNOWN StatementType = iota
	SQL
	DATA
	VARIABLE
	WHERE
	IF
	SET
	ELSE
	ELSEIF
	END
	BOOL
	INT
	STRING
	TIME
	OTHER
	UNCERTAIN
	EXPRESSION
	LOGICAL
	NIL
)

type statement struct {
	Type   StatementType
	Value  string
	Origin string
}

type statements struct {
	data         []statement
	tmpl         []string
	currentIndex int
	Names        map[StatementType]int
}

// Len return length of s.statements
func (s *statements) Len() int {
	return len(s.data)
}

// Next return next slice and increase index by 1
func (s *statements) Next() statement {
	s.currentIndex++
	return s.data[s.currentIndex]
}

// SubIndex take index one step back
func (s *statements) SubIndex() {
	s.currentIndex--
}

// HasMore whether has more slice
func (s *statements) HasMore() bool {
	return s.currentIndex < len(s.data)-1
}

// IsNull whether slice is empty
func (s *statements) IsNull() bool {
	return len(s.data) == 0
}

func (s *statements) appendIfCond(name, cond, result string) {
	s.tmpl = append(s.tmpl, fmt.Sprintf("%s = append(%s, helper.Cond{Cond: %s, Result: func()string {return %s}})", name, name, cond, result))
}

func (s *statements) appendSetValue(name, result string) {
	s.tmpl = append(s.tmpl, fmt.Sprintf("%s = append(%s,  %s)", name, name, strings.TrimSpace(result)))
}

// CreateIf create if clause code
func (s *statements) CreateIf(name string) {
	s.tmpl = append(s.tmpl, fmt.Sprintf("%s := make([]helper.Cond, 0, 100)", name))
}

// CreateStringSet create string set
func (s *statements) CreateStringSet(name string) {
	s.tmpl = append(s.tmpl, fmt.Sprintf("%s := make([]string, 0, 100)", name))
}

// Current return current slice
func (s *statements) Current() statement {
	return s.data[s.currentIndex]
}

func (s *statements) GetName(status StatementType) string {
	switch status {
	case IF:
		defer func() { s.Names[IF]++ }()
		return fmt.Sprintf("ifCond%d", s.Names[IF])
	case WHERE:
		defer func() { s.Names[WHERE]++ }()
		return fmt.Sprintf("whereCond%d", s.Names[WHERE])
	case SET:
		defer func() { s.Names[SET]++ }()
		return fmt.Sprintf("setCond%d", s.Names[SET])
	default:
		return fmt.Sprintf("Cond%d", s.currentIndex)
	}
}

// parse slice and append result to tmpl, return a Clause array
func (s *statements) parse() ([]Clause, error) {
	if s.IsNull() {
		return nil, nil
	}

	name := "generateSQL"
	res := make([]Clause, 0, s.Len())
	for slice := s.Current(); ; slice = s.Next() {
		s.tmpl = append(s.tmpl, "")
		switch slice.Type {
		case SQL, DATA, VARIABLE:
			sqlClause := s.parseSQL(name)
			res = append(res, sqlClause)
			s.tmpl = append(s.tmpl, fmt.Sprintf("%s+=%s", name, sqlClause.String()))
		case IF:
			ifClause, err := s.parseIF()
			if err != nil {
				return nil, err
			}
			res = append(res, ifClause)
			s.tmpl = append(s.tmpl, fmt.Sprintf("%s+=helper.IfClause(%s)", name, ifClause.VarName))
		case WHERE:
			whereClause, err := s.parseWhere()
			if err != nil {
				return nil, err
			}
			res = append(res, whereClause)
			s.tmpl = append(s.tmpl, fmt.Sprintf("%s+=helper.WhereClause(%s)", name, whereClause.VarName))
		case SET:
			setClause, err := s.parseSet()
			if err != nil {
				return nil, err
			}
			res = append(res, setClause)
			s.tmpl = append(s.tmpl, fmt.Sprintf("%s+=helper.SetClause(%s)", name, setClause.VarName))
		case END:
		default:
			return nil, fmt.Errorf("unknow clause:%s", slice.Origin)
		}

		if !s.HasMore() {
			break
		}
	}
	return res, nil
}

// parseIF parse if clause
func (s *statements) parseIF() (res IfClause, err error) {
	slice := s.Current()
	name := s.GetName(slice.Type)
	s.CreateIf(name)

	res.Type = slice.Type
	res.Cond = slice.Value
	res.VarName = name
	cond := []string{res.Cond}
	for s.HasMore() {
		n := s.Next()
		switch n.Type {
		case SQL, DATA, VARIABLE:
			str := s.parseSQL(name)
			res.Value = append(res.Value, str)
			s.appendIfCond(name, res.Cond, str.String())
		case IF:
			var ifClause IfClause
			ifClause, err = s.parseIF()
			if err != nil {
				return
			}
			res.Value = append(res.Value, ifClause)
			s.appendIfCond(name, res.Cond, ifClause.String())
		case WHERE:
			var whereClause WhereClause
			whereClause, err = s.parseWhere()
			if err != nil {
				return
			}
			res.Value = append(res.Value, whereClause)
			s.appendIfCond(name, res.Cond, whereClause.String())
		case SET:
			var setClause SetClause
			setClause, err = s.parseSet()
			if err != nil {
				return
			}
			res.Value = append(res.Value, setClause)
			s.appendIfCond(name, res.Cond, setClause.String())
		case ELSEIF:
			elseClause := s.parseElSE(name)
			elseCond := elseClause.Cond
			elseClause.Cond = fmt.Sprintf("!(%s) && %s", strings.Join(cond, " || "), elseCond)
			res.Else = append(res.Else, elseClause)
			s.appendIfCond(name, elseClause.Cond, elseClause.String())
			cond = append(cond, elseCond)
		case ELSE:
			elseClause := s.parseElSE(name)
			elseClause.Cond = fmt.Sprintf("!(%s)", strings.Join(cond, " || "))
			res.Else = append(res.Else, elseClause)
			s.appendIfCond(name, elseClause.Cond, elseClause.String())
		case END:
			return
		default:
			err = fmt.Errorf("unknow clause : %s", n.Origin)
			return
		}
	}
	if s.Current().Type == END {
		return
	}
	err = fmt.Errorf("incomplete SQL,if not end")
	return
}

// parseElSE parse else clause, the clause' type must be one of if, where, set, SQL condition
func (s *statements) parseElSE(name string) (res ElseClause) {
	slice := s.Current()
	res.Cond = slice.Value
	res.VarName = name
	res.Type = slice.Type

	if !s.HasMore() {
		return
	}
	for n := s.Next(); s.HasMore(); n = s.Next() {
		switch n.Type {
		case SQL, DATA, VARIABLE:
			res.Value = append(res.Value, s.parseSQL(name))
		case IF:
			ifClause, err := s.parseIF()
			if err != nil {
				return
			}
			res.Value = append(res.Value, ifClause)
		case WHERE:
			whereClause, err := s.parseWhere()
			if err != nil {
				return
			}
			res.Value = append(res.Value, whereClause)
		case SET:
			setClause, err := s.parseSet()
			if err != nil {
				return
			}
			res.Value = append(res.Value, setClause)
		default:
			s.SubIndex()
			return
		}
	}
	return
}

// parseWhere parse where clause, the clause' type must be one of if, SQL condition
func (s *statements) parseWhere() (res WhereClause, err error) {
	slice := s.Current()
	name := s.GetName(slice.Type)
	s.CreateStringSet(name)

	res.VarName = name
	res.Type = slice.Type
	for s.HasMore() {
		n := s.Next()
		switch n.Type {
		case SQL, DATA, VARIABLE:
			strClause := s.parseSQL(name)
			res.Value = append(res.Value, strClause)
			s.appendSetValue(name, strClause.String())
		case IF:
			var ifClause IfClause
			ifClause, err = s.parseIF()
			if err != nil {
				return
			}
			res.Value = append(res.Value, ifClause)
			s.appendSetValue(name, ifClause.String())
		case END:
			return
		default:
			err = fmt.Errorf("unknow clause : %s", n.Origin)
			return
		}
	}
	if s.Current().Type == END {
		return
	}
	err = fmt.Errorf("incomplete SQL,where not end")
	return
}

// parseSet parse set clause, the clause' type must be one of if, SQL condition
func (s *statements) parseSet() (res SetClause, err error) {
	slice := s.Current()
	name := s.GetName(slice.Type)
	s.CreateStringSet(name)

	res.VarName = name
	res.Type = slice.Type
	for s.HasMore() {
		n := s.Next()
		switch n.Type {
		case SQL, DATA, VARIABLE:
			strClause := s.parseSQL(name)
			res.Value = append(res.Value, strClause)
			s.appendSetValue(name, strClause.String())
		case IF:
			var ifClause IfClause
			ifClause, err = s.parseIF()
			if err != nil {
				return
			}
			res.Value = append(res.Value, ifClause)
			s.appendSetValue(name, ifClause.String())
		case END:
			return
		default:
			err = fmt.Errorf("unknow clause : %s", n.Origin)
			return
		}
	}
	if s.Current().Type == END {
		return
	}
	err = fmt.Errorf("incomplete SQL,set not end")
	return
}

// parseSQL parse sql condition, the clause' type must be one of SQL condition, VARIABLE, Data
func (s *statements) parseSQL(name string) (res SQLClause) {
	res.VarName = name
	res.Type = SQL
	for slice := s.Current(); ; slice = s.Next() {
		switch slice.Type {
		case SQL, VARIABLE, DATA:
			res.Value = append(res.Value, slice.Value)
		default:
			s.SubIndex()
			return
		}

		if !s.HasMore() {
			return
		}
	}
}

type SQLBuffer struct{ bytes.Buffer }

func (s *SQLBuffer) WriteSql(b byte) {
	switch b {
	case '\n', '\t', ' ':
		if s.Len() == 0 || s.Bytes()[s.Len()-1] != ' ' {
			_ = s.WriteByte(' ')
		}
	default:
		_ = s.WriteByte(b)
	}
}

func (s *SQLBuffer) Dump() string {
	defer s.Reset()
	return s.String()
}

func findParamByName(params []param, name string) (param, bool) {
	for _, p := range params {
		if p.Name == name {
			return p, true
		}
	}
	return param{}, false
}

const (
	argTypeParam = iota + 1
	argTypeNumber
	argTypeString
)

type argInfo struct {
	t     int
	value string
}

type gormOptions struct {
	args []argInfo
	name string
}

type SqlDataInfo struct {
	Name  string
	Value string
}

type ResultDataInfo struct {
	param
	pos string
}

type MethodParser struct {
	MethodName      string
	StructName      string
	FuncDefine      string
	Params          []param
	Results         []param
	ResultData      ResultDataInfo
	ResultIsParam   bool
	RowsAffected    *param
	SqlData         []SqlDataInfo
	Doc             string
	GormOptions     []gormOptions
	Table           string
	SqlTmpList      []string
	WhereConditions []string
	MethodTemplate  *template.Template
}

func (m *MethodParser) HasSqlData() bool {
	return m.SqlData != nil
}

func (m *MethodParser) GetWhereConditionTmp() []string {
	out := []string{""}
	for i, cond := range m.WhereConditions {
		var line string
		line += "whereConditions += \""
		if i > 0 {
			line += " and "
		}
		line += "("
		line += cond
		line += ")\""
		out = append(out, line)
	}
	return out
}

func (m *MethodParser) HasNeedGenerateSql() bool {
	return len(m.SqlTmpList) > 0
}

func (m *MethodParser) GetGORMChainTmp() string {
	notSupportMulti := func() {
		if len(m.GormOptions) > 1 {
			log.Fatalf("Raw not support multiple annotations")
		}
	}
	if len(m.GormOptions) == 0 && m.HasWhereConditions() {
		if m.HasSqlData() {
			return "Table(d.table)." + "Where(whereConditions, params)"
		} else {
			return "Table(d.table)." + "Where(whereConditions)"
		}
	}
	for _, op := range m.GormOptions {
		switch op.name {
		case "Raw":
			notSupportMulti()
			if m.HasSqlData() {
				return "Raw(generateSQL, params)"
			} else {
				return "Raw(generateSQL)"
			}
		case "Exec":
			notSupportMulti()
			if m.HasSqlData() {
				return "Exec(generateSQL, params)"
			} else {
				return "Exec(generateSQL)"
			}
		case "Create":
			//@Create(model)
			//db.Table(table).Create(model)
			notSupportMulti()
			if len(op.args) != 1 {
				log.Fatal("Create annotation annotation only supports one parameter")
			}
			return "Table(d.table).Create(" + op.args[0].value + ")"
		case "UpdateOrCreate":
			//@UpdateOrCreate(update)
			//@Where(conditions)
			//@Result(model)
			if m.ResultData.IsNull() {
				log.Fatal("UpdateOrCreate need result")
			}
			if len(m.WhereConditions) == 0 {
				log.Fatal("UpdateOrCreate need conditions")
			}
			if len(op.args) != 1 {
				log.Fatal("Create annotation only supports one parameter")
			}

			out := "Table(d.table)"
			if m.HasSqlData() {
				out += ".Where(whereConditions, params)"
			} else {
				out += ".Where(whereConditions)"
			}
			out += (".Assign(" + op.args[0].value + ").FirstOrCreate(" + m.ResultData.Name + ")")
			m.ResultData = ResultDataInfo{}
			return out
		}
	}
	return ""
}

func (m *MethodParser) HasWhereConditions() bool {
	return len(m.WhereConditions) > 0
}

func (m *MethodParser) HasGotPoint() bool {
	return !(!m.ResultData.IsArray && (m.ResultData.IsPointer || m.ResultData.IsMap()))
}

func (m *MethodParser) HasNeedNewResult() bool {
	return m.ResultData.pos == "return" && !m.ResultData.IsArray && (m.ResultData.IsPointer || m.ResultData.IsMap())
}

func (m *MethodParser) HasResultData() bool {
	return !m.ResultData.IsNull()
}

func (m *MethodParser) HasResultRowsAffected() bool {
	return m.RowsAffected != nil
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

func (m *MethodParser) GormRunMethodName() string {
	if m.ResultData.IsArray {
		return "Find"
	}
	return "Take"
}

func (m *MethodParser) GetResultsInTmpl() string {
	return paramToString(m.Results)
}

func paramToString(params []param) string {
	var res []string
	for _, param := range params {
		name := param.Name
		if len(name) == 0 && param.Type == "error" {
			name = "err"
		}
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

func (m *MethodParser) isParamExist(paramName string) bool {
	for _, param := range m.SqlData {
		if param.Name == paramName {
			return true
		}
	}
	return false
}

func (m *MethodParser) methodParams(param string, s StatementType) (result statement, err error) {
	for _, p := range m.Params {
		if strings.HasPrefix(param, p.Name+".") || p.Name == param {
			var str string
			switch s {
			case DATA:
				str = fmt.Sprintf("\"@%s\"", param)
				if !m.isParamExist(param) {
					m.SqlData = append(m.SqlData, SqlDataInfo{
						Name:  param,
						Value: param,
					})
				}
			case VARIABLE:
				// if p.Type != "string" {
				// 	err = fmt.Errorf("variable name must be string :%s type is %s", param, p.Type)
				// }
				// str = fmt.Sprintf("%s.Quote(%s)", m.S, param)
			}
			result = statement{
				Type:  s,
				Value: str,
			}
			return
		}
	}
	if param == "table" {
		result = statement{
			Type:  SQL,
			Value: strconv.Quote(m.Table),
		}
		return
	}
	if param == "maxId" {
		result = statement{
			Type:  DATA,
			Value: strconv.Quote(strconv.Itoa(helper.SystemMaxId)),
		}
		return
	}
	if m.isParamExist(param) {
		result = statement{
			Type:  DATA,
			Value: fmt.Sprintf("\"@%s\"", param),
		}
		return
	}
	return result, fmt.Errorf("unknow variable param:%s", param)
}

func (m *MethodParser) getSQLDocString() string {
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

func (m *MethodParser) parseDoc() string {
	docString := strings.TrimSpace(m.getSQLDocString())
	lines := strings.Split(strings.ReplaceAll(docString, "\n\r", "\n"), "\n")
	var outLines []string
	for _, line := range lines {
		line = strings.TrimSpace(line)
		switch {
		case strings.HasPrefix(line, "@Result("):
			end := strings.Index(line, ")")
			if end == -1 {
				log.Fatalf("incomplete sql @Result define, struct: %s,method: %s",
					m.StructName, m.MethodName)
			}
			name := strings.TrimSpace(line[8:end])
			p, ok := findParamByName(m.Params, name)
			if ok {
				m.ResultData = ResultDataInfo{
					param: p,
					pos:   "param",
				}
				m.ResultIsParam = true
			}
			p, ok = findParamByName(m.Results, name)
			if ok {
				m.ResultData = ResultDataInfo{
					param: p,
					pos:   "return",
				}
			}
			if m.ResultData.IsNull() {
				log.Fatalf("@Result defined result not found, struct: %s, method: %s",
					m.StructName, m.MethodName)
			}
		case strings.HasPrefix(line, "@RowsAffected("):
			end := strings.Index(line, ")")
			if end == -1 {
				log.Fatalf("incomplete sql @RowsAffected define, struct: %s,method: %s",
					m.StructName, m.MethodName)
			}
			name := strings.TrimSpace(line[12:end])
			p, ok := findParamByName(m.Params, name)
			if ok {
				m.RowsAffected = &p
			}
			p, ok = findParamByName(m.Results, name)
			if ok {
				m.RowsAffected = &p
			}
			if m.RowsAffected == nil {
				log.Fatalf("@RowsAffected defined result not found, struct: %s, method: %s",
					m.StructName, m.MethodName)
			} else if m.RowsAffected.Type != "int64" {
				log.Fatalf("@RowsAffected param type must be int64, struct: %s, method: %s ",
					m.StructName, m.MethodName)
			}
		case strings.HasPrefix(line, "@Where("):
			end := strings.Index(line, ")")
			if end == -1 {
				log.Fatalf("incomplete sql @Where define, struct: %s,method: %s",
					m.StructName, m.MethodName)
			}
			line := strings.TrimSpace(line[7:end])
			sql, err := m.parseWhereStatement(line)
			if err != nil {
				log.Fatalf("parse where condition error: %s ", err.Error())
			}
			m.WhereConditions = append(m.WhereConditions, sql)
		case strings.HasPrefix(line, "@AddParam("):
			_, line, ok := parseAnnotation(line)
			if !ok {
				log.Fatalf("incomplete sql @Where define, struct: %s,method: %s",
					m.StructName, m.MethodName)
			}
			line = strings.TrimSpace(line)
			err := m.parseAddParam(line)
			if err != nil {
				log.Fatalf("pars AddParam error: %s", err.Error())
			}
		default:
			key, value, ok := parseAnnotation(line)
			if !ok || key == "Sql" {
				outLines = append(outLines, line)
			} else {
				err := m.processPreDefinedOp(key, value)
				if err != nil {
					log.Fatalf("process predefined op: %v", err)
				}
			}
		}
	}

	docString = strings.Join(outLines, "\n")

	switch {
	case strings.HasPrefix(docString, "@Sql("):
		docString = docString[5 : len(docString)-1]
		option := "Raw"
		if m.ResultData.IsNull() {
			option = "Exec"
		}
		m.GormOptions = append(m.GormOptions, gormOptions{
			name: option,
		})
	default:
		if len(m.GormOptions) == 0 && len(m.WhereConditions) == 0 {
			log.Fatalf("undefined notations")
		}
		// matches := annotationRegexp.FindStringSubmatch(docString)
	}

	if strings.HasPrefix(docString, `"`) && strings.HasSuffix(docString, `"`) {
		docString = docString[1 : len(docString)-1]
	}
	return docString
}

func (m *MethodParser) parseAddParam(line string) error {
	args := strings.Split(line, ",")
	if len(args) != 2 {
		return errors.New("AddParam annotation invalid param count")
	}
	buf := bytes.NewBuffer(nil)
	sqlString := args[1]
	for i := 0; i < len(sqlString); i++ {
		b := sqlString[i]
		switch b {
		case '"':
			buf.WriteByte('\\')
			_ = buf.WriteByte(sqlString[i])
			for i++; ; i++ {
				if strOutrange(i, sqlString) {
					return fmt.Errorf("incomplete SQL:%s", sqlString)
				}
				if sqlString[i] == '"' && sqlString[i-1] != '\\' {
					buf.WriteByte('\\')
					_ = buf.WriteByte(sqlString[i])
					break
				}
				_ = buf.WriteByte(sqlString[i])
			}
		case '@':
			i++
			status := DATA
			if sqlString[i] == '@' {
				i++
				status = VARIABLE
			}
			builder := bytes.NewBuffer(nil)
			for ; ; i++ {
				if strOutrange(i, sqlString) || isEnd(sqlString[i]) {
					varString := builder.String()
					_, err := m.methodParams(varString, status)
					if err != nil {
						return fmt.Errorf("sql [%s] varable %s err:%s", sqlString, varString, err)
					}
					buf.WriteString(varString)
					i--
					break
				}
				builder.WriteByte(sqlString[i])
			}
		default:
			buf.WriteByte(b)
		}
	}
	m.SqlData = append(m.SqlData, SqlDataInfo{
		Name:  args[0],
		Value: buf.String(),
	})
	return nil
}

func (m *MethodParser) processPreDefinedOp(opName string, argStr string) error {
	argsVal := strings.Split(argStr, ",")
	var args []argInfo
	for _, arg := range argsVal {
		arg = strings.TrimSpace(arg)
		info := argInfo{value: arg}
		if _, ok := findParamByName(m.Params, arg); ok {
			info.t = argTypeParam
		} else if len(arg) > 0 && arg[0] == '"' && arg[len(arg)-1] == '"' {
			info.t = argTypeString
		} else if _, err := cast.ToFloat64E(arg); len(arg) > 0 && err == nil {
			info.t = argTypeNumber
		} else {
			return errors.New(fmt.Sprintf("invalid param, name: %s,annotation: %s", arg, opName))
		}
		args = append(args, info)
	}
	switch opName {
	case "Create":
		m.GormOptions = append(m.GormOptions, gormOptions{
			args: args,
			name: "Create",
		})
	case "UpdateOrCreate":
		m.GormOptions = append(m.GormOptions, gormOptions{
			args: args,
			name: "UpdateOrCreate",
		})
	default:
		return errors.New("unknown notation: " + opName)
	}
	return nil
}

func strOutrange(index int, str string) bool {
	return index >= len(str)
}

func (m *MethodParser) Parse() error {
	m.MethodTemplate = userDefinedMethod
	res := m.parseDoc()
	l, err := m.parseSql(res)
	if err != nil {
		return err
	}
	m.SqlTmpList = l
	return nil
}

func (m *MethodParser) parseWhereStatement(sqlString string) (string, error) {
	buf := bytes.NewBuffer(nil)
	for i := 0; i < len(sqlString); i++ {
		b := sqlString[i]
		switch b {
		case '"':
			buf.WriteByte('\\')
			_ = buf.WriteByte(sqlString[i])
			for i++; ; i++ {
				if strOutrange(i, sqlString) {
					return "", fmt.Errorf("incomplete SQL:%s", sqlString)
				}
				if sqlString[i] == '"' && sqlString[i-1] != '\\' {
					buf.WriteByte('\\')
					_ = buf.WriteByte(sqlString[i])
					break
				}
				_ = buf.WriteByte(sqlString[i])
			}
		case '@':
			buf.WriteByte(sqlString[i])
			i++
			status := DATA
			if sqlString[i] == '@' {
				i++
				status = VARIABLE
			}
			builder := bytes.NewBuffer(nil)
			for ; ; i++ {
				if strOutrange(i, sqlString) || isEnd(sqlString[i]) {
					varString := builder.String()
					_, err := m.methodParams(varString, status)
					if err != nil {
						return "", fmt.Errorf("sql [%s] varable %s err:%s", sqlString, varString, err)
					}
					buf.WriteString(varString)
					i--
					break
				}
				builder.WriteByte(sqlString[i])
			}
		default:
			buf.WriteByte(b)
		}
	}
	return buf.String(), nil
}

func (m *MethodParser) parseSql(sqlString string) ([]string, error) {

	result := &statements{Names: make(map[StatementType]int)}
	var buf SQLBuffer
	for i := 0; !strOutrange(i, sqlString); i++ {
		b := sqlString[i]
		switch b {
		case '"':
			_ = buf.WriteByte(sqlString[i])
			for i++; ; i++ {
				if strOutrange(i, sqlString) {
					return nil, fmt.Errorf("incomplete SQL:%s", sqlString)
				}
				_ = buf.WriteByte(sqlString[i])
				if sqlString[i] == '"' && sqlString[i-1] != '\\' {
					break
				}
			}
		case '{', '@':
			if sqlClause := buf.Dump(); strings.TrimSpace(sqlClause) != "" {
				result.data = append(result.data, statement{
					Type:  SQL,
					Value: strconv.Quote(sqlClause),
				})
			}

			if strOutrange(i+1, sqlString) {
				return nil, fmt.Errorf("incomplete SQL:%s", sqlString)
			}
			if b == '{' && sqlString[i+1] == '{' {
				for i += 2; ; i++ {
					if strOutrange(i, sqlString) {
						return nil, fmt.Errorf("incomplete SQL:%s", sqlString)
					}
					if sqlString[i] == '"' {
						_ = buf.WriteByte(sqlString[i])
						for i++; ; i++ {
							if strOutrange(i, sqlString) {
								return nil, fmt.Errorf("incomplete SQL:%s", sqlString)
							}
							_ = buf.WriteByte(sqlString[i])
							if sqlString[i] == '"' && sqlString[i-1] != '\\' {
								break
							}
						}
						i++
					}

					if strOutrange(i+1, sqlString) {
						return nil, fmt.Errorf("incomplete SQL:%s", sqlString)
					}
					if sqlString[i] == '}' && sqlString[i+1] == '}' {
						i++

						sqlClause := buf.Dump()
						part, err := checkTemplate(sqlClause, m.Params)
						if err != nil {
							return nil, fmt.Errorf("sql [%s] dynamic template %s err:%w", sqlString, sqlClause, err)
						}
						result.data = append(result.data, part)
						break
					}
					buf.WriteSql(sqlString[i])
				}
			}
			if b == '@' {
				i++
				status := DATA
				if sqlString[i] == '@' {
					i++
					status = VARIABLE
				}
				for ; ; i++ {
					if strOutrange(i, sqlString) || isEnd(sqlString[i]) {
						varString := buf.Dump()
						params, err := m.methodParams(varString, status)
						if err != nil {
							return nil, fmt.Errorf("sql [%s] varable %s err:%s", sqlString, varString, err)
						}
						result.data = append(result.data, params)
						i--
						break
					}
					buf.WriteSql(sqlString[i])
				}
			}
		default:
			buf.WriteSql(b)
		}
	}

	if sqlClause := buf.Dump(); strings.TrimSpace(sqlClause) != "" {
		result.data = append(result.data, statement{
			Type:  SQL,
			Value: strconv.Quote(sqlClause),
		})
	}

	_, err := result.parse()
	if err != nil {
		return nil, fmt.Errorf("sql [%s] parser err:%w", sqlString, err)
	}
	return result.tmpl, nil
}

func isEnd(b byte) bool {
	switch {
	case b >= 'a' && b <= 'z':
		return false
	case b >= 'A' && b <= 'Z':
		return false
	case b >= '0' && b <= '9':
		return false
	case b == '-' || b == '_' || b == '.':
		return false
	default:
		return true
	}
}
