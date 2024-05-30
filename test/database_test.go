//go:build exclude

package test

import (
	"flag"

	"github.com/nyudlts/go-medialog/config"
)

var (
	test_environment   string
	test_configuration string
	test_env           config.Environment
)

func init() {
	flag.StringVar(&environment, "environment", "", "")
	flag.StringVar(&configuration, "config", "", "")
}
