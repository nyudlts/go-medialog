package test

import "github.com/nyudlts/go-medialog/database"

var dbLocat string

func init() {
	dbLocat := "../database/medialog-test.db"

	if err := database.ConnectDatabase(dbLocat); err != nil {
		panic(err)
	}
}
