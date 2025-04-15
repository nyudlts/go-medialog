package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/nyudlts/go-medialog/database"
)

func GlobalSearch(c *gin.Context) {

	user, err := checkLogin(c)
	if err != nil {
		ThrowError(http.StatusBadRequest, err.Error(), c, false)
		return
	}
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
		"isAdmin":    user.IsAdmin,
		"query":      query,
		"entries":    entries,
	})

}
