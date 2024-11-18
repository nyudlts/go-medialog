package controllers

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/nyudlts/go-medialog/database"
)

func getPagination(c *gin.Context) (database.Pagination, error) {
	num_entries := c.Request.URL.Query()["num_entries"]
	var limit = 0
	var err error
	if len(num_entries) > 0 {
		limit, err = strconv.Atoi(num_entries[0])
		if err != nil {
			return database.Pagination{}, err
		}
	}

	page := c.Request.URL.Query()["page"]
	var offset int
	if len(page) > 0 {
		offset, err = strconv.Atoi(page[0])
		if err != nil {
			return database.Pagination{}, err
		}
		if offset < 1 {
			offset = 1
		}
		offset = offset - 1
	}

	return database.Pagination{Limit: limit, Offset: (offset * limit), Sort: "media_id"}, nil
}
