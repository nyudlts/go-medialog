package controllers

import (
	"fmt"
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

var userkey = "user"
var isAdmin = "admin"

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

func DumpSession(c *gin.Context) {
	session := sessions.Default(c)
	userCookie := session.Get(userkey)
	if userCookie == nil {
		c.JSON(http.StatusOK, gin.H{"info": fmt.Sprintf("UserID: nil")})
	} else {
		userID := userCookie.(int)
		c.JSON(http.StatusOK, gin.H{"info": fmt.Sprintf("UserID: %d", userID)})

	}

}
