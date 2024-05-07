package main

import (
	"crypto/sha512"
	"encoding/hex"
	"flag"
	"fmt"

	"github.com/glebarez/sqlite"
	"github.com/nyudlts/go-medialog/controllers"
	"github.com/nyudlts/go-medialog/database"
	"github.com/nyudlts/go-medialog/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var pgdb *gorm.DB
var sqdb *gorm.DB
var test bool
var dbLoc string
var migrateTables bool
var migrateData bool

func init() {
	flag.BoolVar(&test, "test", false, "load the test database")
	flag.BoolVar(&migrateTables, "migrate-tables", false, "migrate tables")
	flag.BoolVar(&migrateData, "migrate-data", false, "migrate data from legacy psql")

}

func main() {
	flag.Parse()
	fmt.Println("go-medialog migration tool")

	if test {
		dbLoc = "medialog-test.db"
	} else {
		dbLoc = "medialog.db"
	}

	if migrateData {
		var err error
		pgdb, err = gorm.Open(postgres.New(postgres.Config{
			DSN:                  "host=localhost user=medialog password=medialog dbname=medialog port=5432 sslmode=disable",
			PreferSimpleProtocol: true,
		}), &gorm.Config{})
		if err != nil {
			panic(err)
		}
	}

	var err error
	sqdb, err = gorm.Open(sqlite.Open(dbLoc), &gorm.Config{})
	if err != nil {
		panic(err)
	}

	if migrateTables {
		if err := migrateDBTables(); err != nil {
			panic(err)
		}
	}

	if migrateData {
		if err := migrateLegacyData(); err != nil {
			panic(err)
		}
	}

}

func migrateDBTables() error {
	fmt.Println("migrating database tables")
	if err := sqdb.AutoMigrate(&models.Repository{}, &models.Accession{}, &models.Collection{}, &models.User{}, &models.Entry{}); err != nil {
		return err
	}
	return nil
}

func migrateLegacyData() error {
	fmt.Println("migrating database data")

	fmt.Println("  * migrating users")
	if err := migrateUsersToGorm(); err != nil {
		return err
	}

	fmt.Println("  * creating repositories")
	if err := populateRepos(); err != nil {
		return err
	}

	fmt.Println("  * migrating resources")
	if err := migrateCollectionsToGorm(); err != nil {
		return err
	}

	fmt.Println("  * migrating accessions")
	if err := migrateAccessionsToGorm(); err != nil {
		return err
	}

	fmt.Println("  * migrating entries")
	if err := migrateEntriesToGorm(); err != nil {
		return err
	}

	return nil
}

func migrateUsersToGorm() error {
	usersPG := []models.UserPG{}
	pgdb.Find(&usersPG)
	for _, userPG := range usersPG {
		u := userPG.ToGormModel()
		sqdb.Create(&u)
	}

	createAdminUser()

	return nil
}

func migrateEntriesToGorm() error {

	mlog_EntryPGs := []models.Mlog_EntryPG{}
	pgdb.Find(&mlog_EntryPGs)
	for _, entryPG := range mlog_EntryPGs {
		e := entryPG.ToGormModel()
		c := models.Collection{}
		if err := sqdb.Where("id = ?", e.CollectionID).First(&c).Error; err != nil {
			return err
		}
		e.RepositoryID = c.RepositoryID

		if err := sqdb.Create(&e).Error; err != nil {
			return err
		}
	}
	return nil
}

func migrateCollectionsToGorm() error {

	collectionsPG := []models.CollectionPG{}
	pgdb.Find(&collectionsPG)
	for _, collectionPG := range collectionsPG {
		c := collectionPG.ToGormModel()
		if c.PartnerCode == "tamwag" {
			c.RepositoryID = 2
		} else if c.PartnerCode == "fales" {
			c.RepositoryID = 3
		} else if c.PartnerCode == "nyuarchives" {
			c.RepositoryID = 6
		}
		sqdb.Create(&c)
	}

	return nil
}

func migrateAccessionsToGorm() error {
	accessionsPG := []models.AccessionPG{}
	pgdb.Find(&accessionsPG)

	for _, accessionPG := range accessionsPG {
		a := accessionPG.ToGormModel()
		sqdb.Create(&a)
	}
	return nil
}

func populateRepos() error {

	tamwag := models.Repository{}
	tamwag.ID = 2
	tamwag.Slug = "tamwag"
	tamwag.Title = "Tamiment Library and Robert F. Wagner Labor Archives"

	fales := models.Repository{}
	fales.ID = 3
	fales.Slug = "fales"
	fales.Title = "Fales Library & Special Collections"

	nyuarchives := models.Repository{}
	nyuarchives.ID = 6
	nyuarchives.Slug = "nyuarchives"
	nyuarchives.Title = "NYU University Archives"

	for _, repo := range []models.Repository{fales, tamwag, nyuarchives} {
		if err := sqdb.Create(&repo).Error; err != nil {
			return err
		}
	}

	return nil
}

func createAdminUser() {
	user := models.User{}
	user.Email = "admin@nyu.edu"
	user.IsActive = true
	user.IsAdmin = true
	password := "test"
	user.Salt = controllers.GenerateStringRunes(16)
	hash := sha512.Sum512([]byte(password + user.Salt))
	user.EncryptedPassword = hex.EncodeToString(hash[:])

	if err := database.ConnectDatabase(false); err != nil {
		panic(err)
	}

	if err := database.InsertUser(&user); err != nil {
		panic(err)
	}

	fmt.Println("    * Admin User Created")
}
