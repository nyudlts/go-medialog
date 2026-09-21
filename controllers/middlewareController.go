package controllers

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/nyudlts/go-medialog/database"
)

const ContextKeySessionCookies = "sessionCookies"
const ContextKeyUser = "user"

func RequireAPIAuth(c *gin.Context) {
	token := c.Request.Header.Get("X-Medialog-Token")
	userID, err := database.FindUserIDByToken(token)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"error": "invalid token",
		})
		return
	}

	c.Set("userID", userID)
	c.Next()

}

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

func CheckToken(c *gin.Context) (string, error) {
	ExpireTokens()
	token := c.Request.Header.Get("X-Medialog-Token")

	if token == "" {
		return "", fmt.Errorf("no `X-Medialog-Token` set in request header")
	}

	apiToken, err := database.FindToken(token)
	if err != nil {
		return "", fmt.Errorf("could not find supplied token: %s", token)
	}

	if !apiToken.IsValid {
		return "", fmt.Errorf("invalid token - please reauthenticate")
	}

	return token, nil
}
