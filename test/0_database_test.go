package test

import (
	"flag"
	"testing"

	config "github.com/nyudlts/go-medialog/config"
	"github.com/nyudlts/go-medialog/database"
	"gorm.io/gorm"
)

var (
	environment   string
	configuration string
	env           config.Environment
	db            *gorm.DB
)

func init() {
	flag.StringVar(&environment, "environment", "", "")
	flag.StringVar(&configuration, "config", "", "")
}

func TestDatabase(t *testing.T) {
	t.Run("Test get the environment", func(t *testing.T) {
		var err error
		env, err = config.GetEnvironment(configuration, environment)
		if err != nil {
			t.Error(err)
		}
		t.Logf("%v", env)
	})

	t.Run("Test Connect Database", func(t *testing.T) {
		if err := database.ConnectMySQL(env.DatabaseConfig); err != nil {
			t.Error(err)
		}
	})

	t.Run("Test get the database", func(t *testing.T) {
		db := database.GetDB()
		t.Logf("%v", db)
	})
}
