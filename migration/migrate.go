//go:build exclude

package main

import (
	"crypto/sha512"
	"encoding/hex"
	"flag"
	"fmt"

	"github.com/nyudlts/go-medialog/config"
	"github.com/nyudlts/go-medialog/controllers"
	"github.com/nyudlts/go-medialog/database"
	"github.com/nyudlts/go-medialog/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var (
	pgdb *gorm.DB
	//sqdb          *gorm.DB
	db            *gorm.DB
	test          bool
	migrateTables bool
	migrateData   bool
	migrateTable  string
	createAdmin   bool
	environment   string
	conf          string
	sqlite        bool
)

func init() {
	flag.BoolVar(&sqlite, "sqlite", false, "")
	flag.StringVar(&environment, "environment", "", "")
	flag.StringVar(&conf, "config", "", "")
	flag.BoolVar(&migrateTables, "migrate-tables", false, "migrate tables")
	flag.BoolVar(&migrateData, "migrate-data", false, "migrate data from legacy psql")
	flag.StringVar(&migrateTable, "migrate-table", "", "migrate a table")
	flag.BoolVar(&createAdmin, "create-admin", false, "")
}

func main() {
	fmt.Println("go-medialog migration tool")
	fmt.Println("parsing flags")
	flag.Parse()

	if migrateData {
		var err error
		pgdb, err = gorm.Open(postgres.New(postgres.Config{
			DSN:                  "host=172.27.16.1 user=medialog password=medialog dbname=medialog port=5432 sslmode=disable",
			PreferSimpleProtocol: true,
		}), &gorm.Config{})
		if err != nil {
			panic(err)
		}
	} else {
		fmt.Println("Skipping connecting to postgres db")
	}

	if sqlite {
		env, err := config.GetSQlite(conf, environment)
		if err != nil {
			panic(err)
		}

		fmt.Println("Migrating:", env.DatabaseLocation)
		if err := database.ConnectDatabase(env.DatabaseLocation); err != nil {
			panic(err)
		}
		fmt.Println("Connected to database")

	} else {

		var err error
		env, err := config.GetEnvironment(conf, environment)
		if err != nil {
			panic(err)
		}

		fmt.Println("Migrating:", env.DatabaseConfig.DatabaseName)

		if err := database.ConnectMySQL(env.DatabaseConfig); err != nil {
			panic(err)
		}

		fmt.Println("Connected to database")
	}

	db = database.GetDB()

	if migrateTable != "" {
		switch migrateTable {
		case "repositories":
			{
				fmt.Print("Migrating repositories table: ")
				if err := db.AutoMigrate(models.Repository{}); err != nil {
					fmt.Printf("ERROR %s\n", err.Error())
				} else {
					fmt.Println("OK")
				}
			}

		case "users":
			{
				fmt.Print("Migrating users table")
				if err := db.AutoMigrate(models.User{}); err != nil {
					fmt.Printf("ERROR %s ", err.Error())
				}
			}

		case "entries":
			{
				fmt.Println("Migrating entries table: ")
				if err := db.AutoMigrate(models.Entry{}); err != nil {
					fmt.Printf("ERROR %s ", err.Error())
				} else {
					fmt.Println("OK")
				}
			}
		case "resources":
			{
				fmt.Println("Migrating resourcess table")
				if err := db.AutoMigrate(models.Resource{}); err != nil {
					fmt.Printf("ERROR %s ", err.Error())
				}
			}
		case "accessions":
			{
				fmt.Println("Migrating accession table")
				if err := db.AutoMigrate(models.Accession{}); err != nil {
					fmt.Printf("ERROR %s ", err.Error())
				}
			}
		default:
			fmt.Printf("ERROR %s is not a valid table to migrate", migrateTable)

		}
	}

	if createAdmin {
		if err := createAdminUser(); err != nil {
			panic(err)
		}
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
	if err := db.AutoMigrate(&models.Repository{}, &models.Accession{}, &models.Resource{}, &models.User{}, &models.Entry{}); err != nil {
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
		db.Create(&u)
	}

	return nil
}

func migrateEntriesToGorm() error {

	mlog_EntryPGs := []models.Mlog_EntryPG{}
	pgdb.Find(&mlog_EntryPGs)
	for _, entryPG := range mlog_EntryPGs {
		e := entryPG.ToGormModel()
		c := models.Resource{}
		if err := db.Where("id = ?", e.ResourceID).First(&c).Error; err != nil {
			fmt.Printf("ERROR: %s", err.Error())
			continue
		}
		e.RepositoryID = c.RepositoryID

		if err := db.Create(&e).Error; err != nil {
			fmt.Printf("ERROR: %s", err.Error())
			continue
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
		db.Create(&c)
	}

	return nil
}

func migrateAccessionsToGorm() error {
	accessionsPG := []models.AccessionPG{}
	pgdb.Find(&accessionsPG)

	for _, accessionPG := range accessionsPG {
		a := accessionPG.ToGormModel()
		db.Create(&a)
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
		if err := db.Create(&repo).Error; err != nil {
			return err
		}
	}

	return nil
}

func createAdminUser() error {
	user := models.User{}
	user.Email = "admin@nyu.edu"
	user.IsActive = true
	user.IsAdmin = true
	password := "test"
	user.Salt = controllers.GenerateStringRunes(16)
	hash := sha512.Sum512([]byte(password + user.Salt))
	user.EncryptedPassword = hex.EncodeToString(hash[:])

	if _, err := database.InsertUser(&user); err != nil {
		return err
	}

	fmt.Println("    * Admin User Created")
	return nil
}
