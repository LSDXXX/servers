package crudgen

const (
	NotEditMarkTemplate = `
// Code generated by crudgen DO NOT EDIT.
// Code generated by crudgen DO NOT EDIT.
// Code generated by crudgen DO NOT EDIT.
	`

	ImportHeaderTemplate = `
package {{.Package}}
import (
	"gorm.io/gorm"

	"github.com/LSDXXX/libs/pkg/container"
	"{{.ModelPackage}}"
	"{{.InterfacePackage}}"
	"{{.HelperPackage}}"
)
	`

	StructDefineTemplate = `

func init() {
	AppendInitFunc(func() {
		_ = container.Singleton(New{{.InterfaceName}})
	})
}

type {{.StructName}} struct {
	db *gorm.DB {{.ContainerTag}} 
	table string
}

func New{{.InterfaceName}}() {{.InterfacePackageTail}}.{{.InterfaceName}} {
	out := &{{.StructName}} {
		table: "{{.TableName}}",
	}
	err := container.Fill(out)
	if err != nil {
		panic(err)
	}
	return out
}

func (d *{{.StructName}}) DB() *gorm.DB {
	return d.db
}

func (d *{{.StructName}}) Table() string {
	return d.table 
}

func (d *{{.StructName}}) WithTable() *gorm.DB {
	return d.db.Table(d.table) 
}

func (d *{{.StructName}}) WithDB(db *gorm.DB) {{.InterfacePackageTail}}.{{.InterfaceName}} {
	return &{{.StructName}} {
		db: db,
		table: d.table,
	}
}

func (d *{{.StructName}}) Page(page, pageSize int, order string, conds ...{{.Model.FullType}}) (result []{{.Model.FullType}}, count int64, err error) { 

	db := d.db.Table(d.table)
	err = db.Count(&count).Error
	if err != nil {
		return 
	}

	if len(conds) > 0 {
		db = db.Where(conds[0])
	}
	db = db.Limit(pageSize).Offset((page-1)*pageSize)
	if len(order) > 0 {
		db = db.Order(order)
	}
	err = db.Find(&result).Error
	return 
}


func (d *{{.StructName}}) Find(conds {{.Model.FullType}}) (result []{{.Model.FullType}}, err error) {
	err = d.db.Table(d.table).Where(conds).Find(&result).Error
	return
}

func (d *{{.StructName}}) Take(order string, conds ...{{.Model.FullType}}) (result {{.Model.FullType}}, err error) {
	db := d.db.Table(d.table).Where(conds)
	if len(order) > 0 {
		db = db.Order(order)
	}
	if len(conds) > 0 {
		db = db.Where(conds[0])
	}
	err = db.Take(&result).Error
	return
}

func (d *{{.StructName}}) Count(conds ...{{.Model.FullType}}) (count int64, err error) {
	db := d.db.Table(d.table)
	if len(conds) > 0 {
		db = db.Where(conds[0])
	}
	err = db.Count(&count).Error
	return
}

func (d *{{.StructName}}) Insert(items ...*{{.Model.FullType}}) error {
	return d.db.Table(d.table).Create(&items).Error
} 

func (d *{{.StructName}}) InsertInBatches(items []*{{.Model.FullType}}, size int) error {
	return d.db.Table(d.table).CreateInBatches(&items, size).Error
} 

func (d *{{.StructName}}) UpdateOrCreate(update *{{.Model.FullType}}, conds {{.Model.FullType}}) error {
	return d.DB().Table(d.table).
		Where(conds).
		Assign(*update).
		FirstOrCreate(update).Error
}

func (d *{{.StructName}}) Updates(updates *{{.Model.FullType}}, conds {{.Model.FullType}}) (rowsAffected int64, err error) {
	res := d.db.Table(d.table).Where(conds).Updates(updates)
	rowsAffected = res.RowsAffected
	err = res.Error
	return
}

func (d *{{.StructName}}) FirstOrCreate(insert *{{.Model.FullType}}, conds {{.Model.FullType}}) (rowsAffected int64, err error) { 
	res := d.db.Table(d.table).
		Where(conds).
		Attrs(*insert).
		FirstOrCreate(insert)
	rowsAffected = res.RowsAffected
	err = res.Error
	return
}

func (d *{{.StructName}}) Delete(conds {{.Model.FullType}}) (rowsAffected int64, err error) {
	res := d.db.Table(d.table).Where(conds).Delete(&{{.Model.FullType}}{})
	rowsAffected = res.RowsAffected
	err = res.Error
	return
}

	`

	UserDefinedMethodTemplate = `
func (d *{{.StructName}}) {{.MethodName}}({{.GetParamInTmpl}})({{.GetResultsInTmpl}}) {
	{{if .HasSqlData}}params := map[string]interface{} { {{range $index,$data:= .SqlData}}
		"{{$data.Name}}": {{$data.Value}}, {{end}}
	}
	{{end}}{{if .HasNeedGenerateSql}}var generateSQL string {{range $line:=.SqlTmpList}}{{$line}}
	{{end}}{{end}}{{if.HasWhereConditions}}var whereConditions string {{range $line:=.GetWhereConditionTmp}}{{$line}}
	{{end}}{{end}}{{if .HasNeedNewResult}}{{.ResultData.Name}} = {{if .ResultData.IsMap}}make{{else}}new{{end}}({{if ne .ResultData.Package ""}}{{.ResultData.Package}}.{{end}}{{.ResultData.Type}}){{end}}
	{{if or .HasResultRowsAffected .HasResultError}}executeSQL:{{else}}_{{end}}= d.DB().{{.GetGORMChainTmp}}{{if .HasResultData}}.{{.GormRunMethodName}}({{if .HasGotPoint}}&{{end}}{{.ResultData.Name}}){{end}}
	{{if .HasResultRowsAffected}}rowsAffected = executeSQL.RowsAffected
	{{end}}{{if .HasResultError}}{{.ResultErrorName}} = executeSQL.Error
	{{end}}return
}
	`
	PageMethodTemplate = `

	`
)
