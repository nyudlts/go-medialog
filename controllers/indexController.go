package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/nyudlts/go-medialog/database"
	"github.com/nyudlts/go-medialog/utils"
)

func GetIndex(c *gin.Context) {

	if !isLoggedIn(c) {
		c.Redirect(302, "/401")
		return
	}

	p := 0
	pagination := utils.Pagination{Limit: 10, Offset: 0, Sort: "updated_at desc"}

	entries, err := database.FindPaginatedEntries(pagination)
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	isAdmin := getCookie("is-admin", c)

	repositoryMap, err := database.GetRepositoryMap()
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	c.HTML(http.StatusOK, "index.html", gin.H{
		"entries":       entries,
		"isAdmin":       isAdmin,
		"page":          p,
		"repositoryMap": repositoryMap,
	})
}
