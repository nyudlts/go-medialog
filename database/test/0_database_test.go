package test

import (
	"flag"
	"testing"

	"github.com/nyudlts/go-medialog/database"
	"github.com/nyudlts/go-medialog/models"
	router "github.com/nyudlts/go-medialog/router"
	"gorm.io/gorm"
)

var (
	environment   string
	configuration string
	env           models.Environment
	db            *gorm.DB
)

func init() {
	flag.StringVar(&environment, "environment", "", "")
	flag.StringVar(&configuration, "config", "", "")
}

func TestDatabase(t *testing.T) {

	t.Run("Test get the environment", func(t *testing.T) {
		var err error
		env, err = router.GetEnvironment(configuration, environment)
		if err != nil {
			t.Error(err)
		}
	})

	t.Run("Test Connect Database", func(t *testing.T) {
		if err := database.ConnectMySQL(env.DatabaseConfig, true); err != nil {
			t.Error(err)
		}
	})

	t.Run("Test get the database", func(t *testing.T) {
		db = database.GetDB()
		t.Logf("%v", db)
	})
}
