package database

import (
	"fmt"

	"github.com/glebarez/sqlite"
	"github.com/nyudlts/go-medialog/config"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var db *gorm.DB
var DatabaseLoc = "medialog.db"
var TestDatabaseLoc = "../database/medialog-test.db"

type Pagination struct {
	Limit  int    `json:"limit"`
	Offset int    `json:"offset"`
	Sort   string `json:"sort"`
}

func ConnectMySQL(dbconfig config.DatabaseConfig) error {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local", dbconfig.Username, dbconfig.Password, dbconfig.URL, dbconfig.Port, dbconfig.DatabaseName)
	var err error
	db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return err
	}
	return nil

}

func ConnectDatabase(dbLoc string) error {
	var err error
	db, err = gorm.Open(sqlite.Open(dbLoc), &gorm.Config{})
	if err != nil {
		return err
	}
	return nil

}

func GetDB() *gorm.DB { return db }
