package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	config "github.com/nyudlts/go-medialog/config"
	router "github.com/nyudlts/go-medialog/router"
)

var (
	environment   string
	configuration string
	sqlite        bool
	gormDebug     bool
	vers          bool
	prod          bool
)

const version = "v0.2.3-beta"

func init() {

	flag.StringVar(&environment, "environment", "", "")
	flag.StringVar(&configuration, "config", "", "")
	flag.BoolVar(&sqlite, "sqlite", false, "")
	flag.BoolVar(&gormDebug, "gorm-debug", false, "")
	flag.BoolVar(&vers, "version", false, "")
	flag.BoolVar(&prod, "prod", false, "")
}

func main() {
	//parse cli flags
	flag.Parse()

	if vers {
		fmt.Printf("{ \"version\": \"%s\"}", version)
		os.Exit(0)
	}

	var r *gin.Engine
	if sqlite {
		env, err := config.GetSQlite(configuration, environment)
		if err != nil {
			panic(err)
		}
		r, err = router.SetupSQRouter(env, gormDebug)
		if err != nil {
			panic(err)
		}
	} else {

		log.Println("Getting Configuration")
		env, err := config.GetEnvironment(configuration, environment)
		if err != nil {
			panic(err)
		}

		log.Println("Setting Up Router")

		r, err = router.SetupRouter(env, gormDebug, prod)
		if err != nil {
			panic(err)
		}
	}

	//start the application
	log.Printf("Running Go-Medialog %s", version)

	if err := r.Run(":8080"); err != nil {
		os.Exit(1)
	}

}
