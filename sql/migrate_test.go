package sql

import (
	"fmt"
	"testing"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func TestMigrate(t *testing.T) {
	sourceName := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		"root", "SDXje)3497uB", "9.134.44.52", 3306, "mysql")
	fmt.Printf(sourceName)
	db, err := gorm.Open(mysql.Open(sourceName), &gorm.Config{})
	if err != nil {
		t.Fatal(err)
	}
	if err = db.Exec("CREATE DATABASE IF NOT EXISTS weiling_desktop_test").Error; err != nil {
		t.Fatal(err)
	}

	sourceName = fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		"root", "SDXje)3497uB", "9.134.44.52", 3306, "weiling_desktop_test")
	fmt.Printf(sourceName)
	db, err = gorm.Open(mysql.Open(sourceName), &gorm.Config{})
	if err != nil {
		t.Fatal(err)
	}
	sqlDB, _ := db.DB()

	err = Up(sqlDB, db.Dialector.Name())
	if err != nil {
		t.Fatal(err)
	}

}
