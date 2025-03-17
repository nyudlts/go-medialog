package controllers

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/nyudlts/go-medialog/database"
)

func GetIndex(c *gin.Context) {

	if err := isLoggedIn(c); err != nil {
		ThrowError(http.StatusUnauthorized, err.Error(), c, false)
		return
	}

	isLoggedIn := true

	pagination := database.Pagination{Limit: 10, Offset: 0, Sort: "updated_at desc", Page: 0}
	pagination.TotalRecords = database.GetCountOfEntriesInDB()
	totalPages := pagination.TotalRecords / int64(pagination.Limit)
	if pagination.TotalRecords%int64(pagination.Limit) > 0 {
		totalPages++
	}
	pagination.TotalPages = int(totalPages)

	entries, err := database.FindPaginatedEntries(pagination)
	if err != nil {
		ThrowError(http.StatusBadRequest, err.Error(), c, isLoggedIn)
		return
	}

	sessionCookies, err := getSessionCookies(c)
	if err != nil {
		ThrowError(http.StatusInternalServerError, err.Error(), c, isLoggedIn)
		return
	}

	user, err := database.GetRedactedUser(sessionCookies.UserID)
	if err != nil {
		ThrowError(http.StatusBadRequest, err.Error(), c, isLoggedIn)
		return
	}

	repositoryMap, err := database.GetRepositoryMap()
	if err != nil {
		ThrowError(http.StatusBadRequest, err.Error(), c, isLoggedIn)
		return
	}

	c.HTML(http.StatusOK, "index.html", gin.H{
		"entries":       entries,
		"isAdmin":       sessionCookies.IsAdmin,
		"pagination":    pagination,
		"repositoryMap": repositoryMap,
		"isLoggedIn":    isLoggedIn,
		"user":          user,
		"limitValues":   LimitValues,
	})
}

func NoRoute(c *gin.Context) {

	if err := isLoggedIn(c); err != nil {
		ThrowError(http.StatusUnauthorized, err.Error(), c, false)
		return
	}

	ThrowError(http.StatusNotFound, fmt.Sprintf("The requested page, %s, does not exist", c.Request.RequestURI), c, true)

}
