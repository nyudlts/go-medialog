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
	sqlite        bool
	sqEnv         config.SQLiteEnv
)

func init() {
	flag.StringVar(&environment, "environment", "", "")
	flag.StringVar(&configuration, "config", "", "")
	flag.BoolVar(&sqlite, "sqlite", false, "")
}

func TestDatabase(t *testing.T) {

	if sqlite {
		t.Run("Test get a sqlite environment", func(t *testing.T) {
			var err error
			sqEnv, err = config.GetSQlite(configuration, environment)
			if err != nil {
				t.Error(err)
			}
		})

		t.Run("Test COnnect SQLite database", func(t *testing.T) {
			if err := database.ConnectSQDatabase(sqEnv, true); err != nil {
				t.Error(err)
			}
		})

	} else {
		t.Run("Test get the mysql environment", func(t *testing.T) {
			var err error
			env, err = config.GetEnvironment(configuration, environment)
			if err != nil {
				t.Error(err)
			}
			t.Logf("%v", env)
		})

		t.Run("Test Connect Database", func(t *testing.T) {
			if err := database.ConnectMySQL(env.DatabaseConfig, true); err != nil {
				t.Error(err)
			}
		})

	}

	t.Run("Test get the database", func(t *testing.T) {
		db = database.GetDB()
		t.Logf("%v", db)
	})

	t.Run("Test get a table", func(t *testing.T) {
		entries, err := database.FindEntries()
		if err != nil {
			t.Error(err)
		}

		t.Logf("%v", len(entries))
	})
}
