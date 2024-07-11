package database

import (
	"github.com/go-gormigrate/gormigrate/v2"
	"github.com/nyudlts/go-medialog/config"
	"github.com/nyudlts/go-medialog/models"
	"gorm.io/gorm"
)

func AutoMigrate(dbc config.DatabaseConfig) error {
	if err := ConnectMySQL(dbc, true); err != nil {
		return err
	}

	if err := db.AutoMigrate(&models.Repository{}, &models.Resource{}, &models.Accession{}, &models.Entry{}, &models.User{}); err != nil {
		return err
	}
	return nil
}

func MigrateDatabase(rollback bool, dbc config.DatabaseConfig) error {
	if err := ConnectMySQL(dbc, true); err != nil {
		return err
	}
	m := gormigrate.New(db, gormigrate.DefaultOptions, []*gormigrate.Migration{
		{
			ID: "20240710 - Adding First Name",
			Migrate: func(tx *gorm.DB) error {
				return tx.Migrator().AddColumn(&models.User{}, "FirstName")
			},
			Rollback: func(tx *gorm.DB) error {
				return tx.Migrator().DropColumn(&models.User{}, "FirstName")
			},
		},
		{
			ID: "20240711 - Adding Last Name",
			Migrate: func(tx *gorm.DB) error {
				return tx.Migrator().AddColumn(&models.User{}, "LastName")
			},
			Rollback: func(tx *gorm.DB) error {
				return tx.Migrator().DropColumn(&models.User{}, "LastName")
			},
		},
	})

	if rollback {
		if err := m.RollbackLast(); err != nil {
			return err
		}
		return nil
	} else {
		if err := m.Migrate(); err != nil {
			return err
		}
		return nil
	}
}
