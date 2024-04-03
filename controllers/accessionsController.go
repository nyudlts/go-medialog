package controllers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/nyudlts/go-medialog/database"
	"github.com/nyudlts/go-medialog/utils"
)

func GetAccessions(c *gin.Context) {
	if err := checkSession(c); err != nil {
		c.Redirect(302, "/")
		return
	}
	isAdmin := getCookie("is-admin", c)

	accessions := database.FindAccessions()
	c.HTML(200, "accessions-index.html", gin.H{
		"accessions":      accessions,
		"isAuthenticated": true,
		"isAdmin":         isAdmin,
	})
}

func GetAccession(c *gin.Context) {
	if err := checkSession(c); err != nil {
		c.Redirect(302, "/")
		return
	}
	isAdmin := getCookie("is-admin", c)

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

	//pagination
	var p = 0
	page := c.Request.URL.Query()["page"]

	if len(page) > 0 {
		p, err = strconv.Atoi(page[0])
		if err != nil {
			c.JSON(http.StatusBadRequest, err.Error())
			return
		}
	}

	pagination := utils.Pagination{Limit: 10, Offset: (p * 10), Sort: "media_id"}

	entries, err := database.FindEntriesByAccessionID(accession.ID, pagination)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	c.HTML(http.StatusOK, "accessions-show.html", gin.H{
		"accession":       accession,
		"resource":        resource,
		"repository":      repository,
		"entries":         entries,
		"isAuthenticated": true,
		"isAdmin":         isAdmin,
		"page":            p,
	})
}
