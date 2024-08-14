package router

import (
	"crypto/sha512"
	"encoding/hex"
	"io"
	"log"
	"os"
	"text/template"
	"time"

	"github.com/gin-contrib/sessions"
	gormsessions "github.com/gin-contrib/sessions/gorm"
	"github.com/gin-gonic/gin"
	"github.com/nyudlts/go-medialog/config"
	"github.com/nyudlts/go-medialog/controllers"
	"github.com/nyudlts/go-medialog/database"
)

func SetupRouter(env config.Environment, gormDebug bool, prod bool) (*gin.Engine, error) {

	log.Println("Medialog starting up")

	if prod {
		gin.SetMode(gin.ReleaseMode)
		log.Println("[INFO] Gin logger")
		//configure logger
		gin.DisableConsoleColor()
		f, _ := os.OpenFile(env.LogLocation, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0755)
		defer f.Close()
		gin.DefaultWriter = io.MultiWriter(f)
	}

	log.Println("[INFO] Setting up router")
	//initialize the router
	r := gin.Default()

	//add global funcs
	SetGlobalFuncs(r)

	//configure the router
	r.LoadHTMLGlob("templates/**/*.html")
	r.StaticFile("/favicon.ico", "./public/favicon.ico")
	r.StaticFile("/test.css", "/public/test.css")
	r.Static("/public", "./public")
	r.SetTrustedProxies([]string{"127.0.0.1"})

	log.Println("[INFO] Connecting to database")
	//connect the database
	if err := database.ConnectMySQL(env.DatabaseConfig, gormDebug); err != nil {
		os.Exit(2)
	}

	log.Println("[INFO] Expiring session tokens")
	if err := database.ExpireAllTokens(); err != nil {
		os.Exit(3)
	}

	log.Println("[INFO] Configuring sessions")
	//configure session parameters

	runes := controllers.GenerateStringRunes(24)
	hash := sha512.Sum512([]byte(runes))
	secret := hex.EncodeToString(hash[:])

	store := gormsessions.NewStore(database.GetDB(), true, []byte(secret))
	options := sessions.Options{}
	options.HttpOnly = true
	options.Domain = "127.0.0.1"
	options.MaxAge = 3600
	r.Use(sessions.Sessions("medialog-sessions", store))

	//load application routes
	log.Println("[INFO] Loading routes")
	LoadRoutes(r)

	return r, nil
}

// Global Functions
func Add(a int, b int) int { return a + b }

func Subtract(a int, b int) int { return a - b }

func Multiply(a int, b int) int { return a * b }

func MultAndAdd(page int, mult int, add int) int {
	return Add(add, Multiply(page, mult))
}

func FormatAsDate(t time.Time) string { return t.Format("2006-01-02") }

func SetGlobalFuncs(router *gin.Engine) {
	router.SetFuncMap(template.FuncMap{
		"formatAsDate":       FormatAsDate,
		"add":                Add,
		"subtract":           Subtract,
		"multiply":           Multiply,
		"getMediatype":       controllers.GetMediaType,
		"getMediatypes":      controllers.GetMediatypes,
		"multAndAdd":         MultAndAdd,
		"storageLocations":   controllers.GetStorageLocations,
		"getStorageLocation": controllers.GetStorageLocation,
	})
}
