package main

import (
	"flag"
	"io"
	"os"

	"github.com/gin-contrib/sessions"
	gormsessions "github.com/gin-contrib/sessions/gorm"
	"github.com/gin-gonic/gin"
	config "github.com/nyudlts/go-medialog/config"
	database "github.com/nyudlts/go-medialog/database"
	routes "github.com/nyudlts/go-medialog/routes"
	"github.com/nyudlts/go-medialog/utils"
)

var (
	router        *gin.Engine
	environment   string
	configuration string
	env           config.Environment
)

const version = "v0.1.0-alpha"

func init() {

	flag.StringVar(&environment, "environment", "", "")
	flag.StringVar(&configuration, "config", "", "")
}

func main() {
	//parse cli flags
	flag.Parse()

	//set the environment variables
	var err error
	env, err = config.GetEnvironment(configuration, environment)
	if err != nil {
		panic(err)
	}

	//get a router
	router, err = setupRouter("templates/**/*.html")
	if err != nil {
		panic(err)
	}

	//start the application
	if err := router.Run(":8080"); err != nil {
		os.Exit(1)
	}

}

func setupRouter(templateLoc string) (*gin.Engine, error) {
	//configure logger
	gin.DisableConsoleColor()
	f, _ := os.Create(env.LogLocation)
	defer f.Close()
	gin.DefaultWriter = io.MultiWriter(f)

	//initialize the router
	r := gin.Default()

	//add global funcs
	utils.SetGlobalFuncs(r)

	//configure the router
	r.LoadHTMLGlob(templateLoc)
	r.StaticFile("/favicon.ico", "./public/favicon.ico")
	r.Static("/public", "./public")
	r.SetTrustedProxies([]string{"127.0.0.1"})

	//connect the database
	if err := database.ConnectDatabase(env.DatabaseLocation); err != nil {
		os.Exit(2)
	}

	//configure session parametes
	store := gormsessions.NewStore(database.GetDB(), true, []byte("secret"))
	options := sessions.Options{}
	options.HttpOnly = true
	options.Domain = "127.0.0.1"
	options.MaxAge = 3600
	r.Use(sessions.Sessions("mysession", store))

	//load applicatin routes
	routes.LoadRoutes(r)

	return r, nil
}
