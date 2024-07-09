package controllers

import (
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/nyudlts/go-medialog/database"
)

func MigrationsIndex(c *gin.Context) {

	isLoggedIn := isLoggedIn(c)
	if !isLoggedIn {
		throwError(http.StatusUnauthorized, UNAUTHORIZED, c)
		return
	}

	sessionCookies, err := getSessionCookies(c)
	if err != nil {
		throwError(http.StatusUnauthorized, err.Error(), c)
		return
	}

	if !sessionCookies.IsAdmin {
		throwError(http.StatusUnauthorized, "not authorized to access this page", c)
		return
	}

	user, err := database.GetRedactedUser(sessionCookies.UserID)
	if !sessionCookies.IsAdmin {
		throwError(http.StatusInternalServerError, err.Error(), c)
		return
	}

	c.HTML(http.StatusOK, "migrations-index.html", gin.H{
		"isAdmin":    sessionCookies.IsAdmin,
		"isLoggedIn": isLoggedIn,
		"user":       user,
	})

}

func MigrateDB(c *gin.Context) {
	isLoggedIn := isLoggedIn(c)
	if !isLoggedIn {
		throwError(http.StatusUnauthorized, UNAUTHORIZED, c)
		return
	}

	sessionCookies, err := getSessionCookies(c)
	if err != nil {
		throwError(http.StatusUnauthorized, err.Error(), c)
		return
	}

	if !sessionCookies.IsAdmin {
		throwError(http.StatusUnauthorized, "not authorized to access this page", c)
		return
	}

	user, err := database.GetRedactedUser(sessionCookies.UserID)
	if !sessionCookies.IsAdmin {
		throwError(http.StatusInternalServerError, err.Error(), c)
		return
	}

	session := sessions.Default(c)
	session.AddFlash("migrations successfully run", "INFO")

	c.HTML(http.StatusOK, "migrations-index.html", gin.H{
		"isAdmin":    sessionCookies.IsAdmin,
		"isLoggedIn": isLoggedIn,
		"user":       user,
		"flash":      session.Flashes("INFO"),
	})
	session.Save()
}
