package main

import (
	"github.com/glebarez/sqlite"
	"github.com/nyudlts/go-medialog/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var pgdb *gorm.DB
var sqdb *gorm.DB

func main() {
	var err error
	pgdb, err = gorm.Open(postgres.New(postgres.Config{
		DSN:                  "host=localhost user=medialog password=medialog dbname=medialog port=5432 sslmode=disable",
		PreferSimpleProtocol: true,
	}), &gorm.Config{})
	if err != nil {
		panic(err)
	}

	sqdb, err = gorm.Open(sqlite.Open("medialog.db"), &gorm.Config{})
	if err != nil {
		panic(err)
	}

	if err := migrateUsersToGorm(); err != nil {
		panic(err)
	}

}

func migrateUsersToGorm() error {
	if err := sqdb.AutoMigrate(&models.User{}); err != nil {
		return err
	}

	usersPG := []models.UserPG{}
	pgdb.Find(&usersPG)
	for _, userPG := range usersPG {
		u := userPG.ToGormModel()
		sqdb.Create(&u)
	}

	return nil
}
func migrateCollectionsToGorm() error {
	if err := sqdb.AutoMigrate(&models.Collection{}); err != nil {
		return err
	}

	collectionsPG := []models.CollectionPG{}
	pgdb.Find(&collectionsPG)
	for _, collectionPG := range collectionsPG {
		c := collectionPG.ToGormModel()
		sqdb.Create(&c)
	}

	return nil
}

func MigrateAccessionToGorm() {
	sqdb.AutoMigrate(&models.Accession{})
	accessionsPG := []models.AccessionPG{}
	pgdb.Find(&accessionsPG)

	for _, accessionPG := range accessionsPG {
		a := accessionPG.ToGormModel()
		sqdb.Create(&a)
	}

}
