package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/nyudlts/go-medialog/database"
	"github.com/nyudlts/go-medialog/models"
)

var LimitValues = []int{10, 25, 50, 100}

func checkLogin(c *gin.Context) (models.User, error) {
	if err := isLoggedIn(c); err != nil {
		return models.User{}, err
	}

	sessionCookies, err := getSessionCookies(c)
	if err != nil {
		return models.User{}, err
	}

	user, err := database.GetRedactedUser(sessionCookies.UserID)
	if err != nil {
		return models.User{}, err
	}
	return user, nil
}
