package database

import (
	"flag"

	"github.com/nyudlts/go-medialog/config"
)

var (
	environment   string
	configuration string
	env           config.Environment
)

func init() {
	flag.StringVar(&environment, "environment", "", "")
	flag.StringVar(&configuration, "config", "", "")
}
