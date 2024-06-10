package database

import (
	"fmt"

	"github.com/glebarez/sqlite"
	"github.com/nyudlts/go-medialog/config"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var db *gorm.DB

type Pagination struct {
	Limit  int    `json:"limit"`
	Offset int    `json:"offset"`
	Sort   string `json:"sort"`
}

func ConnectMySQL(dbconfig config.DatabaseConfig, gormDebug bool) error {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local", dbconfig.Username, dbconfig.Password, dbconfig.URL, dbconfig.Port, dbconfig.DatabaseName)
	var err error
	db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return err
	}

	if gormDebug {
		db.Debug()
	}
	return nil

}

func ConnectSQDatabase(env config.SQLiteEnv, gormDebug bool) error {
	var err error
	db, err = gorm.Open(sqlite.Open(env.DatabaseLocation), &gorm.Config{})
	if err != nil {
		return err
	}

	if gormDebug {
		db.Debug()
	}

	return nil

}

func GetDB() *gorm.DB { return db }
