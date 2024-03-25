package main

import (
	"log"
	"os"

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
	routes.LoadRoutes(router)

	if err := database.ConnectDatabase(); err != nil {
		log.Printf("\t[FATAL]\t[DATABASE]\tdatabase connection failed")
		os.Exit(1)
	}

	if err := router.Run(); err != nil {
		panic(err)
	}

}
