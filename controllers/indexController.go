package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/nyudlts/go-medialog/database"
)

func GetIndex(c *gin.Context) {

	if err := checkLogin(c); err != nil {
		return
	}

	p := 0
	pagination := database.Pagination{Limit: 10, Offset: 0, Sort: "updated_at desc"}

	entries, err := database.FindPaginatedEntries(pagination)
	if err != nil {
		throwError(http.StatusBadRequest, err.Error(), c)
		return
	}

	isAdmin := getCookie("is-admin", c)

	repositoryMap, err := database.GetRepositoryMap()
	if err != nil {
		throwError(http.StatusBadRequest, err.Error(), c)
		return
	}

	c.HTML(http.StatusOK, "index.html", gin.H{
		"entries":       entries,
		"isAdmin":       isAdmin,
		"page":          p,
		"repositoryMap": repositoryMap,
	})
}
