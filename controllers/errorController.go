package controllers

import (
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

func throwError(code int, msg string, c *gin.Context) {
	session := sessions.Default(c)
	session.AddFlash(msg, "WARNING")
	c.HTML(code, "error.html", gin.H{"flash": session.Flashes("WARNING"), "code": code})
	session.Save()
}

func TestError(c *gin.Context) {
	session := sessions.Default(c)
	session.AddFlash("Internal Server Error", "WARNING")
	c.HTML(500, "error.html", gin.H{"flash": session.Flashes("WARNING"), "code": 500})
	session.Save()
}
