package controllers

import (
	"log"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

func throwError(code int, msg string, c *gin.Context) {
	session := sessions.Default(c)
	session.AddFlash(msg, "WARNING")
	var ila bool
	if err := isLoggedIn(c); err != nil {
		ila = false
	} else {
		ila = true
	}
	log.Printf("[ERROR] %d %s", code, msg)
	c.HTML(code, "error.html", gin.H{"flash": session.Flashes("WARNING"), "code": code, "isLoggedIn": ila})
	session.Save()
}

func ThrowError(code int, msg string, c *gin.Context) {
	session := sessions.Default(c)
	session.AddFlash(msg, "WARNING")
	ila := isLoggedIn(c)
	c.HTML(code, "error.html", gin.H{"flash": session.Flashes("WARNING"), "code": code, "isLoggedIn": ila})
	session.Save()
}

func TestError(c *gin.Context) {
	session := sessions.Default(c)
	session.AddFlash("Internal Server Error", "WARNING")
	c.HTML(500, "error.html", gin.H{"flash": session.Flashes("WARNING"), "code": 500})
	session.Save()
}
