package controllers

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/nyudlts/go-medialog/database"
)

func getPagination(c *gin.Context) (database.Pagination, error) {
	numEntries := c.Request.URL.Query()["num_entries"][0]

	limit, err := strconv.Atoi(numEntries)
	if err != nil {
		return database.Pagination{}, err
	}

	page := c.Request.URL.Query()["page"]
	offset, err := strconv.Atoi(page[0])
	if err != nil {
		return database.Pagination{}, err
	}

	if offset < 1 {
		offset = 1
	}

	offset = offset - 1

	return database.Pagination{Limit: limit, Offset: (offset * limit), Sort: "media_id"}, nil
}
