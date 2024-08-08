package controllers

import (
	"log"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

func ThrowError(code int, msg string, c *gin.Context, loggedIn bool) {
	session := sessions.Default(c)
	session.AddFlash(msg, "WARNING")
	log.Printf("[ERROR] %d %s", code, msg)
	c.HTML(code, "error.html", gin.H{"flash": session.Flashes("WARNING"), "code": code, "isLoggedIn": loggedIn})
	session.Save()
}

func TestError(c *gin.Context) {
	session := sessions.Default(c)
	session.AddFlash("Internal Server Error", "WARNING")
	c.HTML(500, "error.html", gin.H{"flash": session.Flashes("WARNING"), "code": 500})
	session.Save()
}
