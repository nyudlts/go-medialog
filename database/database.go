package database

import (
	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
)

var db *gorm.DB
var DatabaseLoc = "medialog.db"

func ConnectDatabase() error {

	var err error
	db, err = gorm.Open(sqlite.Open(DatabaseLoc), &gorm.Config{})
	if err != nil {
		return err
	}
	return nil
}

func GetDB() *gorm.DB { return db }
