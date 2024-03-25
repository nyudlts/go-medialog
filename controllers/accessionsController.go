package controllers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	database "github.com/nyudlts/go-medialog/database"
)

func GetAccessions(c *gin.Context) {
	accessions := database.FindAccessions()
	c.JSON(200, accessions)
}

func GetAccession(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	//pagination
	var p = 1
	page := c.Request.URL.Query()["page"]

	if len(page) > 0 {
		p, err = strconv.Atoi(page[0])
		if err != nil {
			c.JSON(http.StatusBadRequest, err.Error())
		}

		if p == 0 {
			p = 1
		}
	}

	accession, err := database.FindAccession(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	c.HTML(http.StatusOK, "accessions-show.html", gin.H{
		"accession": accession,
	})
}
