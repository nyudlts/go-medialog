package controllers

import (
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/nyudlts/go-medialog/database"
	"github.com/nyudlts/go-medialog/models"
)

func checkSession(c *gin.Context) error {
	sessionKey, err := c.Cookie("session-key")
	if err != nil {
		return fmt.Errorf("no session key - must authenticate")
	}

	session, err := database.FindSessionByKey(sessionKey)
	if err != nil {
		return fmt.Errorf("no session key - must authenticate")

	}

	if session.SessionKey != sessionKey {
		return fmt.Errorf("invalid session key - must authenticate")
	}

	return nil

}

func removeSession(c *gin.Context) error {
	sessionKey, err := c.Cookie("session-key")
	if err != nil {
		return err
	}

	c.SetCookie("session-key", "", -1, "/", "localhost", false, true)
	log.Println(sessionKey)
	err = database.DropSession(sessionKey)
	if err != nil {
		return err
	}

	return nil

}

func NewSession(id uint, c *gin.Context) error {
	session := models.Session{}
	sessionKey := GenerateStringRunes(32)
	session.SessionKey = sessionKey
	session.UserId = int(id)
	session.IsActive = true

	_, err := c.Cookie("session-key")
	if err != nil {
		c.SetCookie("session-key", sessionKey, 3600, "/", "localhost", false, true)
	}

	if err := database.InsertSesssion(session); err != nil {
		return err
	}
	return nil
}
