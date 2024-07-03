//go:build exclude

package main

import (
	"crypto/sha512"
	"encoding/hex"
	"flag"
	"fmt"
	"log"

	"github.com/nyudlts/go-medialog/config"
	"github.com/nyudlts/go-medialog/controllers"
	"github.com/nyudlts/go-medialog/database"
	"github.com/nyudlts/go-medialog/models"
)

const version = "v0.2.6-beta"

var (
	test          bool
	migrateTables bool
	migrateData   bool
	migrateTable  string
	createAdmin   bool
	environment   string
	conf          string
	sqlite        bool
	compare       bool
	clearSessions bool
)

func init() {
	flag.BoolVar(&sqlite, "sqlite", false, "")
	flag.StringVar(&environment, "environment", "", "")
	flag.StringVar(&conf, "config", "", "")
	flag.BoolVar(&migrateTables, "migrate-tables", false, "migrate tables")
	flag.BoolVar(&migrateData, "migrate-data", false, "migrate data from legacy psql")
	flag.StringVar(&migrateTable, "migrate-table", "", "migrate a table")
	flag.BoolVar(&createAdmin, "create-admin", false, "")
	flag.BoolVar(&compare, "compare-data", false, "")
	flag.BoolVar(&clearSessions, "clear-sessions", false, "")
}

func main() {
	fmt.Println("go-medialog migration tool", version)
	fmt.Println("  * Parsing flags")
	flag.Parse()

	if migrateData || compare {
		if err := database.ConnectPGSQL(); err != nil {
			panic(err)
		}
		fmt.Println("  * Connected to Postgres DB")
	} else {
		fmt.Println("  * Skipping connecting to postgres db")
	}

	if sqlite {
		env, err := config.GetSQlite(conf, environment)
		if err != nil {
			panic(err)
		}

		fmt.Println("Migrating:", env.DatabaseLocation)
		if err := database.ConnectSQDatabase(env, false); err != nil {
			panic(err)
		}
		fmt.Println("  * Connected to SQLite3 database")

	} else {

		var err error
		env, err := config.GetEnvironment(conf, environment)
		if err != nil {
			panic(err)
		}

		if err := database.ConnectMySQL(env.DatabaseConfig, false); err != nil {
			panic(err)
		}

		fmt.Println("  * Connected to MySQL database")
	}

	db := database.GetDB()

	if clearSessions {
		fmt.Println("Clearing sessions")
	}

	if migrateTable != "" {

		switch migrateTable {
		case "repositories":
			{
				fmt.Print("  * Migrating repositories table: ")
				if err := db.AutoMigrate(models.Repository{}); err != nil {
					fmt.Printf("ERROR %s\n", err.Error())
				} else {
					fmt.Println("OK")
				}
			}

		case "users":
			{
				fmt.Print("  * Migrating users table")
				if err := db.AutoMigrate(models.User{}); err != nil {
					fmt.Printf("ERROR %s ", err.Error())
				}
			}

		case "entries":
			{
				fmt.Println("  * Migrating entries table: ")
				if err := db.AutoMigrate(models.Entry{}); err != nil {
					fmt.Printf("ERROR %s ", err.Error())
				} else {
					fmt.Println("OK")
				}
			}
		case "resources":
			{
				fmt.Println("  * Migrating resources table")
				if err := db.AutoMigrate(models.Resource{}); err != nil {
					fmt.Printf("ERROR %s ", err.Error())
				}
			}
		case "accessions":
			{
				fmt.Println("  * Migrating accession table")
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

	if compare {
		if err := compareDbs(); err != nil {
			panic(err)
		}
	}

}

func compareDbs() error {
	fmt.Println("  * comparing databases")
	fmt.Println("\nTable\t\tMySQL\t\tPostgreSQL")
	fmt.Println("-----\t\t-----\t\t----------")
	mAccessionCount := database.CountAccessions()
	pAccessionCount := database.CountAccessionsPG()
	fmt.Printf("Accessions\t%d\t\t%d\n", mAccessionCount, pAccessionCount)
	mEntryCount := database.GetCountOfEntriesInDB()
	pEntryCount := database.CountEntriesPG()
	fmt.Printf("Entries\t\t%d\t\t%d\n", mEntryCount, pEntryCount)
	mResourceCount := database.CountResources()
	pResourceCount := database.CountResourcesPG()
	fmt.Printf("Resources\t%d\t\t%d\n", mResourceCount, pResourceCount)
	mUserCount := database.CountUsers()
	pUserCount := database.CountUsersPG()
	fmt.Printf("Users\t\t%d\t\t%d\n", mUserCount, pUserCount)
	fmt.Println()
	return nil
}

func migrateDBTables() error {
	db := database.GetDB()
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

	usersPG, err := database.GetUsersPG()
	if err != nil {
		return err
	}

	for _, userPG := range usersPG {
		u := userPG.ToGormModel()
		if _, err := database.InsertUser(&u); err != nil {
			log.Println("[ERROR] %s", err.Error())
		}
	}

	return nil
}

func migrateAccessionsToGorm() error {
	accessionsPG, err := database.GetAccessionsPG()
	if err != nil {
		return err
	}

	for _, accessionPG := range accessionsPG {
		a := accessionPG.ToGormModel()
		if _, err := database.InsertAccession(&a); err != nil {
			log.Println("[ERROR] %s", err.Error())
			continue
		}
	}
	return nil
}
func migrateEntriesToGorm() error {

	mlogEntryPGs, err := database.GetEntriesPG()
	if err != nil {
		return err
	}

	for _, entryPG := range mlogEntryPGs {
		e := entryPG.ToGormModel()
		c, err := database.FindResource(e.ResourceID)
		if err != nil {
			log.Println("[ERROR] %s", err.Error())
			continue
		}

		e.RepositoryID = c.RepositoryID

		if _, err := database.InsertEntry(&e); err != nil {
			log.Println("[ERROR] %s", err.Error())
			continue
		}
	}
	return nil
}

func migrateCollectionsToGorm() error {

	collectionsPG, err := database.GetCollectionsPG()
	if err != nil {
		log.Println("[ERROR] %s", err.Error())
	}

	for _, collectionPG := range collectionsPG {
		c := collectionPG.ToGormModel()

		switch c.PartnerCode {
		case "tamwag":
			c.RepositoryID = 2
		case "fales":
			c.RepositoryID = 3
		case "nyuarchives":
			c.RepositoryID = 6
		}

		if _, err := database.InsertResource(&c); err != nil {
			log.Println("[ERROR] %s", err.Error())
			continue
		}
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
		if _, err := database.CreateRepository(&repo); err != nil {
			log.Println("[ERROR] %s", err.Error())
		}
	}

	return nil
}

func createAdminUser() error {
	user := models.User{}
	user.Email = "admin@medialog.com"
	user.IsActive = true
	user.IsAdmin = true
	password := controllers.GenerateStringRunes(12)
	log.Printf("[INFO] password set is `%s`", password)
	user.Salt = controllers.GenerateStringRunes(16)
	hash := sha512.Sum512([]byte(password + user.Salt))
	user.EncryptedPassword = hex.EncodeToString(hash[:])

	if _, err := database.InsertUser(&user); err != nil {
		return err
	}

	fmt.Println("    * Admin User Created")
	return nil
}
