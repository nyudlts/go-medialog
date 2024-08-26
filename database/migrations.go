package database

import (
	"fmt"

	"github.com/go-gormigrate/gormigrate/v2"
	"github.com/nyudlts/go-medialog/models"
	"gorm.io/gorm"
)

func AutoMigrate(dbc models.DatabaseConfig) error {
	if err := ConnectMySQL(dbc, true); err != nil {
		return err
	}

	if err := db.AutoMigrate(&models.Repository{}, &models.Resource{}, &models.Accession{}, &models.Entry{}, &models.User{}); err != nil {
		return err
	}
	return nil
}

func MigrateDatabase(rollback bool, dbc models.DatabaseConfig) error {
	if err := ConnectMySQL(dbc, true); err != nil {
		return err
	}

	migrations := []*gormigrate.Migration{
		{
			ID:       "20240710 - Adding First Name",
			Migrate:  func(tx *gorm.DB) error { return tx.Migrator().AddColumn(&models.User{}, "FirstName") },
			Rollback: func(tx *gorm.DB) error { return tx.Migrator().DropColumn(&models.User{}, "FirstName") },
		},
		{
			ID:       "20240711 - Adding Last Name",
			Migrate:  func(tx *gorm.DB) error { return tx.Migrator().AddColumn(&models.User{}, "LastName") },
			Rollback: func(tx *gorm.DB) error { return tx.Migrator().DropColumn(&models.User{}, "LastName") },
		},
		{
			ID:       "20240717 - Adding location to entry",
			Migrate:  func(tx *gorm.DB) error { return tx.Migrator().AddColumn(&models.Entry{}, "Location") },
			Rollback: func(tx *gorm.DB) error { return tx.Migrator().DropColumn(&models.Entry{}, "Location") },
		},
		{
			ID:       "20240805 - Adding API Access to User",
			Migrate:  func(tx *gorm.DB) error { return tx.Migrator().AddColumn(&models.User{}, "CanAccessAPI") },
			Rollback: func(tx *gorm.DB) error { return tx.Migrator().DropColumn(&models.User{}, "CanAccessAPI") },
		},
		{

			ID:       "20240806 - Adding Token table",
			Migrate:  func(tx *gorm.DB) error { return tx.Migrator().CreateTable(&models.Token{}) },
			Rollback: func(tx *gorm.DB) error { return tx.Migrator().DropTable(&models.Token{}) },
		},
		{
			ID:       "20240816 - Adding Token Type to Tokens",
			Migrate:  func(tx *gorm.DB) error { return tx.Migrator().AddColumn(&models.Token{}, "Type") },
			Rollback: func(tx *gorm.DB) error { return tx.Migrator().DropColumn(&models.Token{}, "Type") },
		},
	}

	m := gormigrate.New(db, gormigrate.DefaultOptions, migrations)

	if rollback {
		if err := m.RollbackLast(); err != nil {
			return err
		}
		return nil
	} else {
		if err := m.Migrate(); err != nil {
			return err
		}
		dbMigrations := []string{}
		if err := db.Table("migrations").Find(&dbMigrations).Error; err != nil {
			fmt.Println(err.Error())
		} else {
			fmt.Println("DB Migration complete")
			fmt.Println("Migrations currently applied")
			for _, m := range dbMigrations {
				fmt.Printf(" %s\n", m)
			}
		}
		return nil
	}
}
