package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/nyudlts/go-medialog/database"
)

const ContextKeySessionCookies = "sessionCookies"
const ContextKeyUser = "user"

func RequireAuth(c *gin.Context) {
	if err := isLoggedIn(c); err != nil {
		ThrowError(http.StatusUnauthorized, err.Error(), c, false)
		c.Abort()
		return
	}

	sessionCookies, err := getSessionCookies(c)
	if err != nil {
		ThrowError(http.StatusInternalServerError, err.Error(), c, false)
		c.Abort()
		return
	}

	user, err := database.GetRedactedUser(sessionCookies.UserID)
	if err != nil {
		ThrowError(http.StatusBadRequest, err.Error(), c, false)
		c.Abort()
		return
	}

	c.Set(ContextKeySessionCookies, sessionCookies)
	c.Set(ContextKeyUser, user)
	c.Next()
}
