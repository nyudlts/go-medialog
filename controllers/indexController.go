package controllers

import (
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/nyudlts/go-medialog/database"
	"github.com/nyudlts/go-medialog/utils"
)

func GetIndex(c *gin.Context) {
	var isAuthenticated bool
	if err := checkSession(c); err != nil {
		isAuthenticated = false
	} else {
		isAuthenticated = true
	}

	pagination := utils.Pagination{Limit: 10, Offset: 0, Sort: "updated_at desc"}

	entries, err := database.FindPaginatedEntries(pagination)
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	session := sessions.Default(c)
	if !isAuthenticated {
		session.AddFlash("Please authenticate to access this service", "WARNING")
		session.Save()
	}

	isAdmin := getCookie("is-admin", c)

	c.HTML(http.StatusOK, "index.html", gin.H{
		"entries":         entries,
		"isAuthenticated": isAuthenticated,
		"isAdmin":         isAdmin,
		"flash":           session.Flashes("WARNING"),
	})

	session.Flashes()
	session.Save()
}
