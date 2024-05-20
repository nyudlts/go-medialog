package database

import (
	"github.com/glebarez/sqlite"
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

func ConnectDatabase(dbLoc string) error {
	var err error
	db, err = gorm.Open(sqlite.Open(dbLoc), &gorm.Config{})
	if err != nil {
		return err
	}

	return nil

}

func GetDB() *gorm.DB { return db }
