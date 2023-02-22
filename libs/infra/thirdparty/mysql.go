package thirdparty

import (
	"fmt"
	"time"

	"github.com/LSDXXX/libs/config"
	"github.com/LSDXXX/libs/constant"
	"github.com/LSDXXX/libs/pkg/log"
	"github.com/LSDXXX/sql"
	"github.com/pkg/errors"

	// mysql
	_ "github.com/go-sql-driver/mysql"
	cache "github.com/patrickmn/go-cache"
	"github.com/spf13/cast"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// NewMysqlDB create db
//
//	@param conf
//	@return *gorm.DB
//	@return error
func NewMysqlDB(conf config.MysqlConfig) (*gorm.DB, error) {

	sourceName := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=%s&timeout=5s",
		conf.User, conf.Password, conf.Host, conf.Port, conf.DBName, "Asia%2fShanghai")
	fmt.Println(sourceName)
	db, err := gorm.Open(mysql.Open(sourceName), &gorm.Config{
		Logger: log.NewGormLog(),
	})
	return db, err
}

func NewPgSqlDB(conf config.MysqlConfig) (*gorm.DB, error) {
	sourceName := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%d sslmode=disable TimeZone=Asia/Shanghai connect_timeout=5",
		conf.Host, conf.User, conf.Password, conf.DBName, conf.Port)
	fmt.Println(sourceName)
	db, err := gorm.Open(postgres.Open(sourceName), &gorm.Config{
		Logger: log.NewGormLog(),
	})
	return db, err
}

func SetupDatabase(conf config.MysqlConfig) {
	database := conf.DBName
	conf.DBName = "mysql"
	db, err := NewMysqlDB(conf)
	if err != nil {
		panic(err)
	}
	if err = db.Exec("CREATE DATABASE IF NOT EXISTS " + database).Error; err != nil {
		panic(err)
	}
	conf.DBName = database

	db, err = NewMysqlDB(conf)
	if err != nil {
		panic(err)
	}
	sqlDB, _ := db.DB()
	fmt.Println(sql.Up(sqlDB, db.Dialector.Name()))
}

func NewDBPool() *DBPool {

	return &DBPool{
		cache: cache.New(5*time.Minute, 10*time.Minute),
	}
}

type DBPool struct {
	cache *cache.Cache
}

func (p *DBPool) GetOrCreateDB(catalog string, conf config.MysqlConfig) (*gorm.DB, error) {
	key := conf.Host + "_" + cast.ToString(conf.Port) + "_" + conf.DBName
	v, ok := p.cache.Get(key)
	// v, ok := p.pool.Load(key)
	if ok {
		db := v.(*gorm.DB)
		return db, nil
	}
	var db *gorm.DB
	var err error
	switch catalog {
	case constant.Mysql, constant.Tdsql:
		db, err = NewMysqlDB(conf)
	case constant.Pgsql:
		db, err = NewPgSqlDB(conf)
	default:
		return nil, errors.New("wrong catalog type")
	}
	if err != nil {
		return nil, err
	}
	p.cache.Set(key, db, cache.DefaultExpiration)
	// p.pool.Store(key, db)
	return db, nil
}
