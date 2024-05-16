package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"text/template"

	"github.com/gin-contrib/sessions"
	gormsessions "github.com/gin-contrib/sessions/gorm"
	"github.com/gin-gonic/gin"
	"github.com/nyudlts/go-medialog/controllers"
	database "github.com/nyudlts/go-medialog/database"
	routes "github.com/nyudlts/go-medialog/routes"
	utils "github.com/nyudlts/go-medialog/utils"
	"gopkg.in/yaml.v2"
)

var (
	router      *gin.Engine
	config      string
	environment string
	env         Environment
)

type Environment struct {
	LogLocation     string `yaml:"log"`
	DatbaseLocation string `yaml:"database"`
}

const version = "v0.1.0-alpha"

func init() {
	flag.StringVar(&environment, "environment", "", "")
	flag.StringVar(&config, "config", "", "")
}

func getEnvironment(environment string, configBytes []byte) (Environment, error) {
	envMap := map[string]Environment{}

	err := yaml.Unmarshal(configBytes, &envMap)
	if err != nil {
		return Environment{}, err
	}

	for k, v := range envMap {
		if environment == k {
			return v, nil
		}
	}

	return Environment{}, fmt.Errorf("Error")
}

func main() {
	//parse cli flags
	flag.Parse()

	//check config exists
	if _, err := os.Stat(config); os.IsNotExist(err) {
		panic(err)
	}

	//read the config file
	bytes, err := os.ReadFile(config)
	if err != nil {
		panic(err)
	}

	//get the environment variables
	env, err := getEnvironment(environment, bytes)
	if err != nil {
		panic(err)
	}

	//configure logger
	gin.DisableConsoleColor()
	f, _ := os.Create(env.LogLocation)
	defer f.Close()
	gin.DefaultWriter = io.MultiWriter(f)

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
	if err := database.ConnectDatabase(env.DatbaseLocation); err != nil {
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
