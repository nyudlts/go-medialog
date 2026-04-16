package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/nyudlts/go-medialog/database"
	"github.com/nyudlts/go-medialog/models"
)

func GlobalSearch(c *gin.Context) {
	sessionCookies := c.MustGet(ContextKeySessionCookies).(SessionCookies)
	user := c.MustGet(ContextKeyUser).(models.User)
	query := c.Query("query")

	//get Entry matches
	entries, err := database.SearchEntries(query)
	if err != nil {
		ThrowError(http.StatusInternalServerError, err.Error(), c, true)
		return
	}

	c.HTML(200, "results.html", gin.H{
		"user":       user,
		"isLoggedIn": true,
		"isAdmin": sessionCookies.IsAdmin,
		"query":      query,
		"entries":    entries,
	})

}
