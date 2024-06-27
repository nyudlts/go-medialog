package controllers

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

var userkey = "user"
var isAdmin = "is-admin"

func isLoggedIn(c *gin.Context) bool {
	session := sessions.Default(c)
	sessionKey := session.Get(userkey)
	if sessionKey == nil {
		return false
	}
	return true
}

func getUserkey(c *gin.Context) (int, error) {
	session := sessions.Default(c)
	userKey := session.Get(userkey)
	if userKey == nil {
		return 0, fmt.Errorf("no user key found")
	}
	return userKey.(int), nil
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

func login(userid int, c *gin.Context) error {
	session := sessions.Default(c)
	session.Set(userkey, userid)
	if err := session.Save(); err != nil {
		return err
	}
	return nil
}

func logout(c *gin.Context) {
	session := sessions.Default(c)
	user := session.Get(userkey)
	if user == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid session token"})
		return
	}
	session.Delete(userkey)
	session.Delete(isAdmin)
	if err := session.Save(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save session"})
		return
	}
	c.Redirect(302, "/")
}

type SessionCookies struct {
	UserID  int  `json:"user_id"`
	IsAdmin bool `json:"is_admin"`
}

func DumpSession(c *gin.Context) {
	log.Println("HI")
	session := sessions.Default(c)
	sessionCookies := SessionCookies{}
	userID := session.Get(userkey)
	if userID != nil {
		sessionCookies.UserID = userID.(int)
	}

	adminCookie := session.Get(isAdmin).(bool)
	sessionCookies.IsAdmin = adminCookie
	c.JSON(200, sessionCookies)
	return
}
