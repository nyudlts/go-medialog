package router

import (
	"io"
	"log"
	"os"

	"github.com/gin-contrib/sessions"
	gormsessions "github.com/gin-contrib/sessions/gorm"
	"github.com/gin-gonic/gin"
	"github.com/nyudlts/go-medialog/config"
	"github.com/nyudlts/go-medialog/database"
	"github.com/nyudlts/go-medialog/utils"
)

func SetupRouter(env config.Environment, gormDebug bool, prod bool) (*gin.Engine, error) {

	log.Println("Medialog starting up")

	if prod {
		gin.SetMode(gin.ReleaseMode)
		log.Println("  ** Configuring Gin logger")
		//configure logger
		gin.DisableConsoleColor()
		f, _ := os.OpenFile(env.LogLocation, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0755)
		defer f.Close()
		gin.DefaultWriter = io.MultiWriter(f)
	}

	log.Println("  ** Setting up router")
	//initialize the router
	r := gin.Default()

	//add global funcs
	utils.SetGlobalFuncs(r)

	//configure the router
	r.LoadHTMLGlob("templates/**/*.html")
	r.StaticFile("/favicon.ico", "./public/favicon.ico")
	r.StaticFile("/test.css", "/public/test.css")
	r.Static("/public", "./public")
	r.SetTrustedProxies([]string{"127.0.0.1"})

	log.Println("  ** Connecting to database")
	//connect the database
	if err := database.ConnectMySQL(env.DatabaseConfig, gormDebug); err != nil {
		os.Exit(2)
	}

	log.Println("  ** Configuring sessions")
	//configure session parameters

	store := gormsessions.NewStore(database.GetDB(), true, []byte("secret"))
	options := sessions.Options{}
	options.HttpOnly = true
	options.Domain = "127.0.0.1"
	options.MaxAge = 3600
	r.Use(sessions.Sessions("mysession", store))

	//load application routes
	log.Println("  ** Loading routes")
	LoadRoutes(r)

	return r, nil
}
