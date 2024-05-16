package test

import (
	"fmt"
	"os"

	"github.com/nyudlts/go-medialog/database"
	"gopkg.in/yaml.v2"
)

var dbLocat string
var config string

type Environment struct {
	LogLocation     string `yaml:"log"`
	DatbaseLocation string `yaml:"database"`
}

func init() {

	bytes, err := os.ReadFile("../config/go-medialog.conf")
	if err != nil {
		panic(err)
	}

	env, err := getEnvironment("test", bytes)
	if err != nil {
		panic(err)
	}

	if err := database.ConnectDatabase(env.DatbaseLocation); err != nil {
		panic(err)
	}
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
