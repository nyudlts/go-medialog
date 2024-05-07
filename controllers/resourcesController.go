package controllers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/nyudlts/go-medialog/database"
	"github.com/nyudlts/go-medialog/utils"
)

func GetResource(c *gin.Context) {
	if !isLoggedIn(c) {
		c.Redirect(302, "/error")
		return
	}

	isAdmin := getCookie("is-admin", c)

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

	summary, err := database.GetSummaryByResource(resource.ID)
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	accessions, err := database.FindAccessionsByResourceID(resource.ID)
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
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

	entries, err := database.FindEntriesByResourceID(resource.ID, pagination)
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	users, err := getUserEmailMap([]int{resource.CreatedBy, resource.UpdatedBy})
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	c.HTML(http.StatusOK, "resources-show.html", gin.H{
		"resource":        resource,
		"accessions":      accessions,
		"entries":         entries,
		"isAdmin":         isAdmin,
		"isAuthenticated": true,
		"page":            p,
		"summary":         summary,
		"totals":          summary.GetTotals(),
		"users":           users,
	})
}

func GetResources(c *gin.Context) {
	if !isLoggedIn(c) {
		c.Redirect(302, "/error")
		return
	}

	isAdmin := getCookie("is-admin", c)

	resources, err := database.FindResources()
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	repositoryMap, err := database.GetRepositoryMap()
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	c.HTML(http.StatusOK, "resources-index.html", gin.H{
		"resources":       resources,
		"isAuthenticated": true,
		"isAdmin":         isAdmin,
		"repositoryMap":   repositoryMap,
	})
}
