package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/nyudlts/go-medialog/database"
)

func GetIndex(c *gin.Context) {
	var isAuthenticated bool
	if err := checkSession(c); err != nil {
		isAuthenticated = false
	} else {
		isAuthenticated = true
	}

	entries, err := database.FindEntriesSorted(10)
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	c.HTML(http.StatusOK, "index.html", gin.H{
		"entries":         entries,
		"isAuthenticated": isAuthenticated,
	})
}
