package database

import (
	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
)

var db *gorm.DB
var DatabaseLoc = "medialog.db"
var TestDatabaseLoc = "medialog-test.db"

func ConnectDatabase(test bool) error {
	var dbLoc string
	if test {
		dbLoc = DatabaseLoc
	} else {
		dbLoc = TestDatabaseLoc
	}
	var err error
	db, err = gorm.Open(sqlite.Open(dbLoc), &gorm.Config{})
	if err != nil {
		return err
	}
	return nil

}

func GetDB() *gorm.DB { return db }
