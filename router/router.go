package router

import (
	"crypto/sha512"
	"encoding/hex"
	"fmt"
	"io"
	"log"
	"os"
	"text/template"
	"time"

	"github.com/gin-contrib/sessions"
	gormsessions "github.com/gin-contrib/sessions/gorm"
	"github.com/gin-gonic/gin"
	"github.com/nyudlts/go-medialog/controllers"
	"github.com/nyudlts/go-medialog/database"
	"github.com/nyudlts/go-medialog/models"
	"gopkg.in/yaml.v2"
)

func SetupRouter(env models.Environment, gormDebug bool, prod bool) (*gin.Engine, error) {

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
	r.Static("/public", "./public")
	r.SetTrustedProxies([]string{"127.0.0.1"})

	log.Println("[INFO] Connecting to database")
	//connect the database
	if err := database.ConnectMySQL(env.DatabaseConfig, gormDebug); err != nil {
		os.Exit(2)
	}

	if prod {
		log.Println("[INFO] Expiring session tokens")
		if err := database.ExpireAllTokens(); err != nil {
			os.Exit(3)
		}
	}

	//configure session parameters
	log.Println("[INFO] Configuring sessions")
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

	//load api routes
	log.Println("[INFO] Loading API")
	LoadAPI(r)

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

func Iterate(count int) []int {
	var i int
	var Items []int
	for i = 0; i < count; i++ {
		Items = append(Items, i)
	}
	return Items
}

func SetGlobalFuncs(router *gin.Engine) {
	router.SetFuncMap(template.FuncMap{
		"formatAsDate":          FormatAsDate,
		"add":                   Add,
		"subtract":              Subtract,
		"multiply":              Multiply,
		"getMediatype":          controllers.GetMediaType,
		"getMediatypes":         controllers.GetMediatypes,
		"multAndAdd":            MultAndAdd,
		"storageLocations":      controllers.GetStorageLocations,
		"getStorageLocation":    controllers.GetStorageLocation,
		"entryStatuses":         controllers.GetEntryStatuses,
		"getEntryStatus":        controllers.GetEntryStatus,
		"getOpticalContentType": controllers.GetOpticalContentType,
		"iterate":               Iterate,
	})
}

func GetEnvironment(config string, environment string) (models.Environment, error) {

	//check that the config exists
	if _, err := os.Stat(config); os.IsNotExist(err) {
		panic(err)
	}

	//read the config
	configBytes, err := os.ReadFile(config)
	if err != nil {
		return models.Environment{}, err
	}

	envMap := map[string]models.Environment{}
	if err := yaml.Unmarshal(configBytes, &envMap); err != nil {
		return models.Environment{}, err
	}

	for k, v := range envMap {
		if environment == k {
			return v, nil
		}
	}

	return models.Environment{}, fmt.Errorf("environment %s does not exist in configuration", environment)
}
