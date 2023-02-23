package repo

import (
	"github.com/LSDXXX/libs/model"
	"github.com/LSDXXX/libs/pkg/crudgen/helper"
)

//go:generate crudgentool -f $GOFILE -op ../infra
//@Table(user)
type UserMapper interface {
	helper.DAO[UserMapper, model.User]

	//@Sql(select * from @@table
	//	where name = @name
	//)
	//@Result(res)
	GetByUserName(name string) (res model.User, err error)
}
