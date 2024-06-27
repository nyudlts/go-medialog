package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v2"
)

type Environment struct {
	LogLocation    string         `yaml:"log"`
	DatabaseConfig DatabaseConfig `yaml:"database"`
}

type SQLiteEnv struct {
	LogLocation      string `yaml:"log"`
	DatabaseLocation string `yaml:"database"`
}

type DatabaseConfig struct {
	Username     string `yaml:"username"`
	Password     string `yaml:"password"`
	URL          string `yaml:"url"`
	Port         string `yaml:"port"`
	DatabaseName string `yaml:"database_name"`
}

func GetEnvironment(config string, environment string) (Environment, error) {

	//check that the config exists
	if _, err := os.Stat(config); os.IsNotExist(err) {
		panic(err)
	}

	//read the config
	configBytes, err := os.ReadFile(config)
	if err != nil {
		return Environment{}, err
	}

	envMap := map[string]Environment{}
	if err := yaml.Unmarshal(configBytes, &envMap); err != nil {
		return Environment{}, err
	}

	for k, v := range envMap {
		if environment == k {
			return v, nil
		}
	}

	return Environment{}, fmt.Errorf("environment %s does not exist in configuration", environment)
}

func GetSQlite(config string, environment string) (SQLiteEnv, error) {

	//check that the config exists
	if _, err := os.Stat(config); os.IsNotExist(err) {
		return SQLiteEnv{}, err
	}

	//read the config
	configBytes, err := os.ReadFile(config)
	if err != nil {
		return SQLiteEnv{}, err
	}

	envMap := map[string]SQLiteEnv{}
	if err := yaml.Unmarshal(configBytes, &envMap); err != nil {
		return SQLiteEnv{}, err
	}

	for k, v := range envMap {
		if environment == k {
			return v, nil
		}
	}

	return SQLiteEnv{}, fmt.Errorf("envirnoment %s does not exist in configuration", environment)
}
