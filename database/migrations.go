package database

import (
	"log"

	"github.com/go-gormigrate/gormigrate/v2"
	"github.com/nyudlts/go-medialog/models"
	"gorm.io/gorm"
)

func MigrateDatabase() error {
	m := gormigrate.New(db, gormigrate.DefaultOptions, []*gormigrate.Migration{
		{
			ID: "20040709 - Adding First and Last Name",
			Migrate: func(tx *gorm.DB) error {
				type User struct {
					FirstName string
					LastName  string
				}
				if err := tx.Migrator().AddColumn(&models.User{}, "FirstName"); err != nil {
					return err
				}
				if err := tx.Migrator().AddColumn(&models.User{}, "LastName"); err != nil {
					return err
				}
				return nil
			},
			Rollback: func(tx *gorm.DB) error {
				type User struct {
					FirstName string
					LastName  string
				}
				if err := db.Migrator().DropColumn(&models.User{}, "FirstName"); err != nil {
					return err
				}
				if err := tx.Migrator().DropColumn(&models.User{}, "LastName"); err != nil {
					return err
				}
				return nil
			},
		},
	})

	if err := m.Migrate(); err != nil {
		log.Println("Migration did run successfully")
		return err
	}

	return nil
}
