package main

import (
	"flag"
	"os"

	config "github.com/nyudlts/go-medialog/config"
	router "github.com/nyudlts/go-medialog/router"
)

var (
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
	r, err := router.SetupRouter(env)
	if err != nil {
		panic(err)
	}

	//start the application
	if err := r.Run(":8080"); err != nil {
		os.Exit(1)
	}

}
