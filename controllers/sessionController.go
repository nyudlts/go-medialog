package controllers

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

func checkSession(c *gin.Context) error {
	session := sessions.Default(c)
	sessionKey := session.Get("session-key")
	if sessionKey == nil {
		return fmt.Errorf("no session key found")
	}

	return nil
}

func getCookie(key string, c *gin.Context) interface{} {
	session := sessions.Default(c)
	return session.Get(key)
}

func setCookie(name string, value interface{}, c *gin.Context) {
	session := sessions.Default(c)
	session.Set(name, value)
	session.Save()
}

func newSession(c *gin.Context) error {
	session := sessions.Default(c)
	log.Println(session)
	sessionKey := GenerateStringRunes(32)
	session.Set("session-key", sessionKey)
	//session.Options(sessions.Options{MaxAge: 3600 * 4})
	session.Save()
	return nil
}

func LogoutUser(c *gin.Context) {
	fmt.Println("Hello")
	session := sessions.Default(c)
	session.Delete("session-key")
	session.Clear()
	session.Options(sessions.Options{MaxAge: -1})
	session.Save()
	c.Redirect(http.StatusPermanentRedirect, "/sessions/dump")
}

func DumpSession(c *gin.Context) {
	session := sessions.Default(c)
	key := session.Get("session-key")
	c.JSON(200, fmt.Sprintf("%v\n%v", session, key))
}
