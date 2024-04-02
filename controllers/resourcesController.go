package controllers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/nyudlts/go-medialog/database"
)

func GetResource(c *gin.Context) {
	if err := checkSession(c); err != nil {
		c.Redirect(302, "/")
		return
	}

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	resource, err := database.FindResource(uint(id))
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	repository, err := database.FindRepository(uint(resource.RepositoryID))
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	accessions, err := database.FindAccessionsByResourceID(resource.ID)
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	entries, err := database.FindEntriesByResourceID(resource.ID)
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	c.HTML(http.StatusOK, "resources-show.html", gin.H{
		"resource":   resource,
		"repository": repository,
		"accessions": accessions,
		"entries":    entries,
	})
}

func GetResources(c *gin.Context) {
	if err := checkSession(c); err != nil {
		c.Redirect(302, "/")
		return
	}

	resources, err := database.FindResources()
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	c.HTML(http.StatusOK, "resources-index.html", gin.H{
		"resources":       resources,
		"isAuthenticated": true,
	})
}
