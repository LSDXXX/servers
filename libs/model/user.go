package model

type User struct {
	Id       int
	UserName string `gorm:"column:user_name"`
	Password string `gorm:"column:password"`
}
