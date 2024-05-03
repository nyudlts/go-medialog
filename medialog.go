package main

import (
	"log"
	"os"
	"text/template"

	"github.com/gin-contrib/sessions"
	gormsessions "github.com/gin-contrib/sessions/gorm"
	"github.com/gin-gonic/gin"
	"github.com/nyudlts/go-medialog/controllers"
	database "github.com/nyudlts/go-medialog/database"
	routes "github.com/nyudlts/go-medialog/routes"
	utils "github.com/nyudlts/go-medialog/utils"
)

var router *gin.Engine

func main() {
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

	if err := database.ConnectDatabase(false); err != nil {
		os.Exit(1)
	}

	store := gormsessions.NewStore(database.GetDB(), true, []byte("secret"))
	options := sessions.Options{}
	options.HttpOnly = true
	options.Domain = "127.0.0.1"
	log.Println(options)

	router.Use(sessions.Sessions("mysession", store))
	routes.LoadRoutes(router)

	if err := router.Run(); err != nil {
		panic(err)
	}

}
