package main

import (
	"crypto/sha512"
	"encoding/hex"
	"flag"
	"fmt"
	"os"

	"github.com/nyudlts/go-medialog/controllers"
	"github.com/nyudlts/go-medialog/database"
	"github.com/nyudlts/go-medialog/models"
	"gopkg.in/yaml.v2"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var (
	pgdb          *gorm.DB
	sqdb          *gorm.DB
	test          bool
	migrateTables bool
	migrateData   bool
	migrateTable  string
	createAdmin   bool
	environment   string
	config        string
	env           Environment
)

type Environment struct {
	LogLocation     string `yaml:"log"`
	DatbaseLocation string `yaml:"database"`
}

func init() {
	flag.StringVar(&environment, "environment", "", "")
	flag.StringVar(&config, "config", "", "")
	flag.BoolVar(&migrateTables, "migrate-tables", false, "migrate tables")
	flag.BoolVar(&migrateData, "migrate-data", false, "migrate data from legacy psql")
	flag.StringVar(&migrateTable, "migrate-table", "", "migrate a table")
	flag.BoolVar(&createAdmin, "create-admin", false, "")
}

func getEnvironment(environment string, configBytes []byte) (Environment, error) {
	envMap := map[string]Environment{}

	err := yaml.Unmarshal(configBytes, &envMap)
	if err != nil {
		return Environment{}, err
	}

	for k, v := range envMap {
		if environment == k {
			return v, nil
		}
	}

	return Environment{}, fmt.Errorf("Error")
}

func main() {
	fmt.Println("go-medialog migration tool")
	fmt.Println("parsing flags")
	flag.Parse()

	if migrateData {
		var err error
		pgdb, err = gorm.Open(postgres.New(postgres.Config{
			DSN:                  "host=localhost user=medialog password=medialog dbname=medialog port=5432 sslmode=disable",
			PreferSimpleProtocol: true,
		}), &gorm.Config{})
		if err != nil {
			panic(err)
		}
	} else {
		fmt.Println("Skipping connecting to postgres db")
	}

	bytes, err := os.ReadFile(config)
	if err != nil {
		panic(err)
	}

	env, err := getEnvironment(environment, bytes)
	if err != nil {
		panic(err)
	}

	fmt.Println("Migrating", env.DatbaseLocation)

	if err := database.ConnectDatabase(env.DatbaseLocation); err != nil {
		panic(err)
	} else {
		fmt.Println("Connected to database")
	}

	sqdb = database.GetDB()

	if migrateTable != "" {
		switch migrateTable {
		case "repositories":
			{
				fmt.Print("Migrating repositories table: ")
				if err := sqdb.AutoMigrate(models.Repository{}); err != nil {
					fmt.Printf("ERROR %s\n", err.Error())
				} else {
					fmt.Println("OK")
				}
			}

		case "users":
			{
				fmt.Print("Migrating users table")
				if err := sqdb.AutoMigrate(models.User{}); err != nil {
					fmt.Printf("ERROR %s ", err.Error())
				}
			}

		case "entries":
			{
				fmt.Println("Migrating entries table: ")
				if err := sqdb.AutoMigrate(models.Entry{}); err != nil {
					fmt.Printf("ERROR %s ", err.Error())
				} else {
					fmt.Println("OK")
				}
			}
		case "collections":
			{
				fmt.Println("Migrating collections table")
				if err := sqdb.AutoMigrate(models.Collection{}); err != nil {
					fmt.Printf("ERROR %s ", err.Error())
				}
			}
		case "accessions":
			{
				fmt.Println("Migrating accession table")
				if err := sqdb.AutoMigrate(models.Accession{}); err != nil {
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