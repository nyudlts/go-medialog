package controllers

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	database "github.com/nyudlts/go-medialog/database"
)

func GetEntry(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	log.Println(id)

	entry, err := database.FindEntry(id.String())
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	accession, err := database.FindAccession(entry.AccessionID)
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	resource, err := database.FindResource(uint(accession.CollectionID))
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	repository, err := database.FindRepository(uint(resource.RepositoryID))
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	c.HTML(http.StatusOK, "entries-show.html", gin.H{
		"entry":      entry,
		"accession":  accession,
		"resource":   resource,
		"repository": repository,
	})
}

func GetEntries(c *gin.Context) {
	entries, err := database.FindEntries()
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	c.HTML(http.StatusOK, "entries-index.html", gin.H{
		"entries": entries,
	})
}
