package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/nyudlts/go-medialog/controllers"
	"github.com/nyudlts/go-medialog/database"
	"github.com/nyudlts/go-medialog/models"
	router "github.com/nyudlts/go-medialog/router"
)

var (
	environment   string
	configuration string
	gormDebug     bool
	vers          bool
	prod          bool
	migrate       bool
	rollback      bool
	automigrate   bool
	createAdmin   bool
)

const version = "v1.0.9"

func init() {
	flag.StringVar(&environment, "environment", "", "")
	flag.StringVar(&configuration, "config", "", "")
	flag.BoolVar(&gormDebug, "gorm-debug", false, "")
	flag.BoolVar(&vers, "version", false, "")
	flag.BoolVar(&prod, "prod", false, "")
	flag.BoolVar(&migrate, "migrate", false, "")
	flag.BoolVar(&automigrate, "automigrate", false, "")
	flag.BoolVar(&rollback, "rollback", false, "")
	flag.BoolVar(&createAdmin, "create-admin", false, "")
}

var r *gin.Engine
var env models.Environment

func main() {
	//parse cli flags
	flag.Parse()

	if vers {
		fmt.Printf("{ \"version\": \"%s\"}", version)
		os.Exit(0)
	}

	var err error
	env, err = router.GetEnvironment(configuration, environment)
	if err != nil {
		panic(err)
	}

	if migrate || rollback {
		fmt.Println("[]running migrations")

		if err := database.MigrateDatabase(rollback, env.DatabaseConfig); err != nil {
			panic(err)
		}

		os.Exit(0)
	}

	if automigrate {
		fmt.Println("auto-migrating database")
		if err := database.AutoMigrate(env.DatabaseConfig); err != nil {
			panic(err)
		}
		os.Exit(0)
	}

	if prod {
		logFile, err := os.OpenFile(env.LogLocation, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0755)
		if err != nil {
			panic(err)
		}
		defer logFile.Close()

		log.SetOutput(logFile)
		log.Println("[INFO] Medialog starting up")
		log.Printf("[INFO] Logging to %s", env.LogLocation)
		log.Println("[INFO] Setting Up Router")
	}

	r, err = router.SetupRouter(env, gormDebug, prod)
	if err != nil {
		log.Fatal(err)
	}

	if createAdmin {
		password, err := controllers.CreateAdminUser()
		if err != nil {
			panic(err)
		}

		fmt.Printf("admin user create with password `%s`", password)
		os.Exit(0)
	}

	//start the application
	log.Printf("[INFO] Running Go-Medialog %s", version)

	if err := r.Run(":8080"); err != nil {
		log.Fatal(err)
	}

}
