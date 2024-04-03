package controllers

import (
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/nyudlts/go-medialog/database"
	"github.com/nyudlts/go-medialog/utils"
)

func GetEntry(c *gin.Context) {
	if err := checkSession(c); err != nil {
		c.Redirect(302, "/")
		return
	}
	isAdmin := getCookie("is-admin", c)

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
		"entry":           entry,
		"accession":       accession,
		"resource":        resource,
		"repository":      repository,
		"isAuthenticated": true,
		"isAdmin":         isAdmin,
	})
}

func GetEntries(c *gin.Context) {
	if err := checkSession(c); err != nil {
		c.Redirect(302, "/")
		return
	}
	isAdmin := getCookie("is-admin", c)

	//pagination
	var p = 0
	var err error
	page := c.Request.URL.Query()["page"]

	if len(page) > 0 {
		p, err = strconv.Atoi(page[0])
		if err != nil {
			c.JSON(http.StatusBadRequest, err.Error())
			return
		}

	}

	pagination := utils.Pagination{Limit: 10, Offset: (p * 10), Sort: "updated_at desc"}

	entries, err := database.FindPaginatedEntries(pagination)
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	c.HTML(http.StatusOK, "entries-index.html", gin.H{
		"entries":         entries,
		"isAuthenticated": true,
		"isAdmin":         isAdmin,
		"page":            p,
	})
}
