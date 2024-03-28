package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/nyudlts/go-medialog/database"
)

func GetUsers(c *gin.Context) {
	users, err := database.FindUsers()
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	c.HTML(http.StatusOK, "users-index.html", gin.H{
		"users": users,
	})
}
