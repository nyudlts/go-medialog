package controllers

import (
	"fmt"
	"log"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

func checkSession(c *gin.Context) error {
	session := sessions.Default(c)
	sessionKey := session.Get("session-key")
	if sessionKey == nil {
		session.AddFlash("INFO", "no session key found - must authenticate")
		session.Save()
		return fmt.Errorf("no session key found")
	}
	return nil
}

func removeSession(c *gin.Context) error {
	session := sessions.Default(c)
	session.Delete("session-key")
	session.Save()
	return nil
}

func NewSession(c *gin.Context) error {
	session := sessions.Default(c)
	log.Println(session)
	sessionKey := GenerateStringRunes(32)
	session.Set("session-key", sessionKey)
	session.Save()
	return nil
}
