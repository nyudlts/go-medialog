package main

import (
	"log"
	"os"

	"github.com/gin-contrib/sessions"
	gormsessions "github.com/gin-contrib/sessions/gorm"
	"github.com/gin-gonic/gin"
	database "github.com/nyudlts/go-medialog/database"
	routes "github.com/nyudlts/go-medialog/routes"
)

var router *gin.Engine

func main() {
	router = gin.Default()

	//configure the router
	router.LoadHTMLGlob("templates/**/*.html")
	router.StaticFile("/favicon.ico", "./public/favicon.ico")
	router.Static("/public", "./public")
	router.SetTrustedProxies([]string{"127.0.0.1"})

	if err := database.ConnectDatabase(); err != nil {
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
