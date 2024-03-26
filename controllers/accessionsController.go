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

	accession, err := database.FindAccession(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	resource, err := database.FindResource(uint(accession.CollectionID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	repository, err := database.FindRepository(uint(resource.RepositoryID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	entries, err := database.FindEntriesByAccessionID(accession.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	c.HTML(http.StatusOK, "accessions-show.html", gin.H{
		"accession":  accession,
		"resource":   resource,
		"repository": repository,
		"entries":    entries,
	})
}
