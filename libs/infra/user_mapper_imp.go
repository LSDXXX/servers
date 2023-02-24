// Code generated by crudgen DO NOT EDIT.
// Code generated by crudgen DO NOT EDIT.
// Code generated by crudgen DO NOT EDIT.

package infra

import (
	"gorm.io/gorm"

	"github.com/LSDXXX/libs/model"
	"github.com/LSDXXX/libs/pkg/container"
	"github.com/LSDXXX/libs/repo"
)

func init() {
	AppendInitFunc(func() {
		_ = container.Singleton(NewUserMapper)
	})
}

type UserMapperImp struct {
	db    *gorm.DB `container:"type"`
	table string
}

func NewUserMapper() repo.UserMapper {
	out := &UserMapperImp{
		table: "user",
	}
	err := container.Fill(out)
	if err != nil {
		panic(err)
	}
	return out
}

func (d *UserMapperImp) DB() *gorm.DB {
	return d.db
}

func (d *UserMapperImp) Table() string {
	return d.table
}

func (d *UserMapperImp) WithTable() *gorm.DB {
	return d.db.Table(d.table)
}

func (d *UserMapperImp) WithDB(db *gorm.DB) repo.UserMapper {
	return &UserMapperImp{
		db:    db,
		table: d.table,
	}
}

func (d *UserMapperImp) Page(page, pageSize int, order string, conds ...model.User) (result []model.User, count int64, err error) {

	db := d.db.Table(d.table)
	err = db.Count(&count).Error
	if err != nil {
		return
	}

	if len(conds) > 0 {
		db = db.Where(conds[0])
	}
	db = db.Limit(pageSize).Offset((page - 1) * pageSize)
	if len(order) > 0 {
		db = db.Order(order)
	}
	err = db.Find(&result).Error
	return
}

func (d *UserMapperImp) Find(conds model.User) (result []model.User, err error) {
	err = d.db.Table(d.table).Where(conds).Find(&result).Error
	return
}

func (d *UserMapperImp) Take(order string, conds ...model.User) (result model.User, err error) {
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

func (d *UserMapperImp) Count(conds ...model.User) (count int64, err error) {
	db := d.db.Table(d.table)
	if len(conds) > 0 {
		db = db.Where(conds[0])
	}
	err = db.Count(&count).Error
	return
}

func (d *UserMapperImp) Insert(items ...*model.User) error {
	return d.db.Table(d.table).Create(&items).Error
}

func (d *UserMapperImp) InsertInBatches(items []*model.User, size int) error {
	return d.db.Table(d.table).CreateInBatches(&items, size).Error
}

func (d *UserMapperImp) UpdateOrCreate(update *model.User, conds model.User) error {
	return d.DB().Table(d.table).
		Where(conds).
		Assign(*update).
		FirstOrCreate(update).Error
}

func (d *UserMapperImp) Updates(updates *model.User, conds model.User) (rowsAffected int64, err error) {
	res := d.db.Table(d.table).Where(conds).Updates(updates)
	rowsAffected = res.RowsAffected
	err = res.Error
	return
}

func (d *UserMapperImp) FirstOrCreate(insert *model.User, conds model.User) (rowsAffected int64, err error) {
	res := d.db.Table(d.table).
		Where(conds).
		Attrs(*insert).
		FirstOrCreate(insert)
	rowsAffected = res.RowsAffected
	err = res.Error
	return
}

func (d *UserMapperImp) Delete(conds model.User) (rowsAffected int64, err error) {
	res := d.db.Table(d.table).Where(conds).Delete(&model.User{})
	rowsAffected = res.RowsAffected
	err = res.Error
	return
}

func (d *UserMapperImp) GetByUserName(name string) (res model.User, err error) {
	params := map[string]interface{}{
		"name": name,
	}
	var generateSQL string
	generateSQL += "select * from user where name = @name"

	executeSQL := d.DB().Raw(generateSQL, params).Take(&res)
	err = executeSQL.Error
	return
}
