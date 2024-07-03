//go:build exclude

package main

import (
	"crypto/sha512"
	"encoding/hex"
	"flag"
	"fmt"
	"log"
	"os"

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
	createUsers   bool
	environment   string
	conf          string
	sqlite        bool
	compare       bool
)

func init() {
	flag.BoolVar(&sqlite, "sqlite", false, "")
	flag.StringVar(&environment, "environment", "", "")
	flag.StringVar(&conf, "config", "", "")
	flag.BoolVar(&migrateTables, "migrate-tables", false, "migrate tables")
	flag.BoolVar(&migrateData, "migrate-data", false, "migrate data from legacy psql")
	flag.StringVar(&migrateTable, "migrate-table", "", "migrate a table")
	flag.BoolVar(&createUsers, "create-users", false, "")
	flag.BoolVar(&compare, "compare-data", false, "")
}

func main() {
	fmt.Println("go-medialog migration tool", version)
	fmt.Println("  * Parsing flags")
	flag.Parse()

	logFile, _ := os.OpenFile("migration.log", os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0755)
	defer logFile.Close()
	log.SetOutput(logFile)

	if migrateData || compare {
		if err := database.ConnectPGSQL(); err != nil {
			panic(err)
		}
		fmt.Println("  * Connected to Postgres DB")
		log.Println("[INFO] connected to Postgres DB")
	} else {
		fmt.Println("  * Skipping connecting to postgres db")
		log.Println("[INFO] Skipping connecting to postgres db")
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
		log.Println("[INFO] Connected to MySQL db")
	}

	db := database.GetDB()

	if migrateTable != "" {

		switch migrateTable {
		case "repositories":
			{
				fmt.Print("  * Migrating repositories table: ")
				log.Println("[INFO] Migrating repositories table")
				if err := db.AutoMigrate(models.Repository{}); err != nil {
					fmt.Printf("ERROR %s\n", err.Error())
					log.Printf("[ERROR] %s", err.Error())
				} else {
					fmt.Println("OK")
					log.Println("[INFO] Repositories table migration complete")
				}
			}

		case "users":
			{
				fmt.Print("  * Migrating users table")
				log.Println("[INFO] Migrating users table")
				if err := db.AutoMigrate(models.User{}); err != nil {
					fmt.Printf("ERROR %s\n", err.Error())
					log.Printf("[ERROR] %s", err.Error())
				} else {
					fmt.Println("OK")
					log.Println("[INFO] Users table migration complete")
				}
			}

		case "entries":
			{
				fmt.Println("  * Migrating entries table: ")
				log.Println("[INFO] Migrating entries table")
				if err := db.AutoMigrate(models.Entry{}); err != nil {
					fmt.Printf("ERROR %s ", err.Error())
					log.Printf("[ERROR] %s", err.Error())
				} else {
					fmt.Println("OK")
					log.Println("[INFO] entries table migration complete")
				}
			}
		case "resources":
			{
				fmt.Println("  * Migrating resources table")
				if err := db.AutoMigrate(models.Resource{}); err != nil {
					fmt.Printf("ERROR %s ", err.Error())
					log.Printf("[ERROR] %s", err.Error())
				} else {
					fmt.Println("OK")
					log.Println("[INFO] Resources table migration complete")
				}
			}
		case "accessions":
			{
				fmt.Println("  * Migrating accession table")
				if err := db.AutoMigrate(models.Accession{}); err != nil {
					fmt.Printf("ERROR %s ", err.Error())
					log.Printf("[ERROR] %s", err.Error())
				} else {
					fmt.Println("OK")
					log.Println("[INFO] Accessions table migration complete")
				}
			}
		default:
			fmt.Printf("ERROR %s is not a valid table to migrate", migrateTable)
			log.Printf("[ERROR] %s", migrateTable)
		}
	}

	if createUsers {
		if err := createAdminUser(); err != nil {
			log.Printf("[ERROR] %s", err.Error())
		}

		if err := createUnknownUser(); err != nil {
			log.Printf("[ERROR] %s", err.Error())
		}
	}

	if migrateTables {
		if err := migrateDBTables(); err != nil {
			log.Printf("[ERROR] %s", err.Error())
		}
	}

	if migrateData {
		if err := migrateLegacyData(); err != nil {
			log.Printf("[ERROR] %s", err.Error())
		}
	}

	if compare {
		if err := compareDbs(); err != nil {
			log.Printf("[ERROR] %s", err.Error())
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
	log.Println("[INFO] migrating database tables")
	if err := db.AutoMigrate(&models.Repository{}, &models.Accession{}, &models.Resource{}, &models.User{}, &models.Entry{}); err != nil {
		fmt.Printf("ERROR %s ", err.Error())
		log.Printf("[ERROR] %s", err.Error())
		return err
	}
	return nil
}

func migrateLegacyData() error {
	fmt.Println("migrating database data")
	log.Println("[INFO] migrating database data")

	fmt.Println("  * migrating users")
	log.Println("[INFO] migrating users data")
	if err := migrateUsersToGorm(); err != nil {
		return err
	}

	fmt.Println("  * creating repositories")
	log.Println("[INFO] creating repositories data")
	if err := populateRepos(); err != nil {
		return err
	}

	fmt.Println("  * migrating resources")
	log.Println("[INFO] migrating resources data")
	if err := migrateCollectionsToGorm(); err != nil {
		return err
	}

	fmt.Println("  * migrating accessions")
	log.Println("[INFO] migrating accessions data")
	if err := migrateAccessionsToGorm(); err != nil {
		return err
	}

	fmt.Println("  * migrating entries")
	log.Println("[INFO] migrating entries data")
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

		if u.CreatedBy == 0 {
			u.CreatedBy = 101
		}

		if u.UpdatedBy == 0 {
			u.UpdatedBy = 101
		}

		if _, err := database.InsertUser(&u); err != nil {
			log.Printf("[ERROR] %s", err.Error())
			continue
		}
		log.Printf("[INFO] User %d %s created", u.ID, u.Email)
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

		if a.CreatedBy == 0 {
			a.CreatedBy = 101
		}

		if a.UpdatedBy == 0 {
			a.UpdatedBy = 101
		}

		if _, err := database.InsertAccession(&a); err != nil {
			log.Printf("[ERROR] %s", err.Error())
			continue
		}
		log.Printf("[INFO] Accession %d %s created", a.ID, a.AccessionNum)
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
			log.Printf("[ERROR] %s", err.Error())
			continue
		}

		e.RepositoryID = c.RepositoryID

		if e.CreatedBy == 0 {
			e.CreatedBy = 101
		}

		if e.UpdatedBy == 0 {
			e.UpdatedBy = 101
		}

		if _, err := database.InsertEntry(&e); err != nil {
			log.Printf("[ERROR] %s", err.Error())
			continue
		}
		log.Printf("[INFO] Entry %s created", e.ID.String())
	}
	return nil
}

func migrateCollectionsToGorm() error {

	collectionsPG, err := database.GetCollectionsPG()
	if err != nil {
		log.Printf("[ERROR] %s", err.Error())
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

		if c.CreatedBy == 0 {
			c.CreatedBy = 101
		}

		if c.UpdatedBy == 0 {
			c.UpdatedBy = 101
		}

		if _, err := database.InsertResource(&c); err != nil {
			log.Printf("[ERROR] %s", err.Error())
			continue
		}

		log.Printf("[INFO] Resource %d %s created", c.ID, c.CollectionCode)
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
		repo.CreatedBy = 100
		repo.UpdatedBy = 100
		if _, err := database.CreateRepository(&repo); err != nil {
			log.Println("[ERROR] %s", err.Error())
		}
		log.Printf("[INFO] Repository %d %s created", repo.ID, repo.Slug)
	}

	return nil
}

func createAdminUser() error {
	user := models.User{}
	user.Email = "admin@medialog.dlib.nyu.edu"
	user.IsActive = true
	user.IsAdmin = true
	user.ID = 100
	user.CreatedBy = 100
	user.UpdatedBy = 100
	password := controllers.GenerateStringRunes(12)
	log.Printf("[INFO] admin password set is `%s`", password)
	user.Salt = controllers.GenerateStringRunes(16)
	hash := sha512.Sum512([]byte(password + user.Salt))
	user.EncryptedPassword = hex.EncodeToString(hash[:])

	if _, err := database.InsertUser(&user); err != nil {
		return err
	}

	fmt.Println("    * Admin User Created")
	return nil
}

func createUnknownUser() error {
	user := models.User{}
	user.Email = "unknown@medialog.dlib.nyu.edu"
	user.IsActive = false
	user.IsAdmin = false
	user.ID = 101
	user.CreatedBy = 100
	user.UpdatedBy = 100
	password := controllers.GenerateStringRunes(12)
	log.Printf("[INFO] unknown password set is `%s`", password)
	user.Salt = controllers.GenerateStringRunes(16)
	hash := sha512.Sum512([]byte(password + user.Salt))
	user.EncryptedPassword = hex.EncodeToString(hash[:])

	if _, err := database.InsertUser(&user); err != nil {
		return err
	}

	fmt.Println("    * Unknown User Created")
	return nil
}
