package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/nyudlts/go-medialog/database"
)

func GetIndex(c *gin.Context) {

	if !isLoggedIn(c) {
		throwError(http.StatusUnauthorized, UNAUTHORIZED, c)
		return
	}

	p := 0
	pagination := database.Pagination{Limit: 10, Offset: 0, Sort: "updated_at desc"}

	entries, err := database.FindPaginatedEntries(pagination)
	if err != nil {
		throwError(http.StatusBadRequest, err.Error(), c)
		return
	}

	entryCount := database.GetCountOfEntriesInDB()

	sessionCookies, err := getSessionCookies(c)
	if err != nil {
		throwError(http.StatusInternalServerError, err.Error(), c)
		return
	}

	user, err := database.GetRedactedUser(sessionCookies.UserID)
	if err != nil {
		throwError(http.StatusBadRequest, err.Error(), c)
		return
	}

	repositoryMap, err := database.GetRepositoryMap()
	if err != nil {
		throwError(http.StatusBadRequest, err.Error(), c)
		return
	}

	c.HTML(http.StatusOK, "index.html", gin.H{
		"entries":       entries,
		"isAdmin":       sessionCookies.IsAdmin,
		"page":          p,
		"repositoryMap": repositoryMap,
		"entryCount":    entryCount,
		"isLoggedIn":    isLoggedIn,
		"user":          user,
	})
}
