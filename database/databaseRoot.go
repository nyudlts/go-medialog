package database

import (
	"fmt"

	"github.com/nyudlts/go-medialog/models"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var db *gorm.DB

type Pagination struct {
	Limit  int    `json:"limit"`
	Offset int    `json:"offset"`
	Sort   string `json:"sort"`
}

type DateRange struct {
	StartYear    int `form:"start-year"`
	StartMonth   int `form:"start-month"`
	StartDay     int `form:"start-day"`
	EndYear      int `form:"end-year"`
	EndMonth     int `form:"end-month"`
	EndDay       int `form:"end-day"`
	RepositoryID int `form:"repository-id"`
}

func ConnectMySQL(dbconfig models.DatabaseConfig, gormDebug bool) error {
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

func GetDB() *gorm.DB { return db }
