package controllers

import (
	"log"
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/nyudlts/go-medialog/database"
	"github.com/nyudlts/go-medialog/models"
)

func ThrowError(code int, msg string, c *gin.Context, loggedIn bool) {
	session := sessions.Default(c)
	session.AddFlash(msg, "WARNING")
	var user = models.User{}
	if loggedIn {
		sessionCookies, err := getSessionCookies(c)
		if err != nil {
			ThrowError(http.StatusInternalServerError, err.Error(), c, loggedIn)
			return
		}
		user, err = database.GetRedactedUser(sessionCookies.UserID)
		if err != nil {
			ThrowError(http.StatusBadRequest, err.Error(), c, loggedIn)
			return
		}
	}

	log.Printf("[ERROR] %d %s", code, msg)
	c.HTML(code, "error.html", gin.H{
		"flash":      session.Flashes("WARNING"),
		"code":       code,
		"isLoggedIn": loggedIn,
		"isAdmin":    user.IsAdmin,
		"user":       user,
	})
	session.Save()
}

func TestError(c *gin.Context) {
	session := sessions.Default(c)
	session.AddFlash("Internal Server Error", "WARNING")
	c.HTML(500, "error.html", gin.H{"flash": session.Flashes("WARNING"), "code": 500})
	session.Save()
}
