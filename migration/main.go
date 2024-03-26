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

	if err := AutoMigrate(); err != nil {
		panic(err)
	}

	if err := migrateReposToGorm(); err != nil {
		panic(err)
	}

	if err := migrateCollectionsToGorm(); err != nil {
		panic(err)
	}

	if err := migrateAccessionsToGorm(); err != nil {
		panic(err)
	}

	if err := migrateEntriesToGorm(); err != nil {
		panic(err)
	}

	if err := migrateUsersToGorm(); err != nil {
		panic(err)
	}

}

func AutoMigrate() error {
	if err := sqdb.AutoMigrate(&models.Repository{}, &models.Accession{}, &models.Collection{}, &models.User{}, &models.Entry{}); err != nil {
		return err
	}
	return nil
}

func migrateReposToGorm() error {
	if err := sqdb.AutoMigrate(&models.Repository{}); err != nil {
		return err
	}
	return nil
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

func migrateEntriesToGorm() error {
	if err := sqdb.AutoMigrate(&models.Entry{}); err != nil {
		return err
	}

	mlog_EntryPGs := []models.Mlog_EntryPG{}
	pgdb.Find(&mlog_EntryPGs)
	for _, entryPG := range mlog_EntryPGs {
		e := entryPG.ToGormModel()
		sqdb.Create(&e)
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

func migrateAccessionsToGorm() error {
	if err := sqdb.AutoMigrate(&models.Accession{}); err != nil {
		return err
	}
	accessionsPG := []models.AccessionPG{}
	pgdb.Find(&accessionsPG)

	for _, accessionPG := range accessionsPG {
		a := accessionPG.ToGormModel()
		sqdb.Create(&a)
	}
	return nil
}
