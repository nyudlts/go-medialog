package main

import (
	"flag"
	"os"
	"path/filepath"
	"text/template"

	"github.com/gin-contrib/sessions"
	gormsessions "github.com/gin-contrib/sessions/gorm"
	"github.com/gin-gonic/gin"
	"github.com/nyudlts/go-medialog/controllers"
	database "github.com/nyudlts/go-medialog/database"
	routes "github.com/nyudlts/go-medialog/routes"
	utils "github.com/nyudlts/go-medialog/utils"
)

var (
	router *gin.Engine
	test   bool
)

func init() {
	flag.BoolVar(&test, "test", false, "run application against test db")
}

func main() {
	//parse cli flags
	flag.Parse()

	//initialize the router
	router = gin.Default()

	//add global funcs
	router.SetFuncMap(template.FuncMap{
		"formatAsDate": utils.FormatAsDate,
		"add":          utils.Add,
		"subtract":     utils.Subtract,
		"getMediatype": controllers.GetMediaType,
	})

	//configure the router
	router.LoadHTMLGlob("templates/**/*.html")
	router.StaticFile("/favicon.ico", "./public/favicon.ico")
	router.Static("/public", "./public")
	router.SetTrustedProxies([]string{"127.0.0.1"})

	//connect the database
	if err := database.ConnectDatabase(filepath.Join("database", "medialog-test.db")); err != nil {
		os.Exit(2)
	}

	//configure session parametes
	store := gormsessions.NewStore(database.GetDB(), true, []byte("secret"))
	options := sessions.Options{}
	options.HttpOnly = true
	options.Domain = "127.0.0.1"
	options.MaxAge = 3600
	router.Use(sessions.Sessions("mysession", store))

	//load applicatin routes
	routes.LoadRoutes(router)

	//start the application
	if err := router.Run(); err != nil {
		os.Exit(1)
	}

}
