package controllers

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/nyudlts/go-medialog/database"
)

var userkey = "user"
var isAdmin = "is-admin"
var canAccessAPI = "can-access-api"
var sessionToken = "token"

func ExpireTokens() {
	tokens := database.GetTokens()
	log.Printf("[INFO] expiring api tokens")
	for _, token := range tokens {
		if token.IsValid && time.Now().After(token.Expires) {
			//log.Printf("[INFO] Expiring token %d", token.ID)
			if err := database.ExpireToken(token.ID); err != nil {
				log.Printf("[ERROR] %s", err.Error())
			}
		}
	}
}

func isLoggedIn(c *gin.Context) error {

	ExpireTokens()

	session := sessions.Default(c)
	userIDCookie := session.Get(userkey)
	if userIDCookie == nil {
		return fmt.Errorf("please reauthenticate (no user key)")
	}

	userID := userIDCookie.(int)

	tokenCookie := session.Get(sessionToken)
	if tokenCookie == nil {
		return fmt.Errorf("please reauthenticate (no token)")
	}

	token := tokenCookie.(string)

	sessionToken, err := database.FindToken(token)
	if err != nil {
		return fmt.Errorf("please reauthenticate (token not found)")
	}

	if !sessionToken.IsValid {
		return fmt.Errorf("please reauthenticate (token not valid)")
	}

	if sessionToken.UserID != uint(userID) {
		return fmt.Errorf("please reauthenticate (session key not set for user)")
	}

	return nil
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
	c.Redirect(http.StatusFound, "/")
}

type SessionCookies struct {
	UserID       int    `json:"user_id"`
	IsAdmin      bool   `json:"is_admin"`
	CanAccessAPI bool   `json:"can_access_api"`
	SessionToken string `json:"session_token"`
}

func getSessionCookies(c *gin.Context) (SessionCookies, error) {
	session := sessions.Default(c)
	sessionCookies := SessionCookies{}
	userID := session.Get(userkey)
	if userID == nil {
		return sessionCookies, fmt.Errorf("no user key")
	}
	sessionCookies.UserID = userID.(int)

	adminCookie := session.Get(isAdmin)
	if adminCookie == nil {
		return sessionCookies, fmt.Errorf("user must be admin")
	}
	sessionCookies.IsAdmin = adminCookie.(bool)

	apiCookie := session.Get(canAccessAPI)
	if apiCookie == nil {
		return sessionCookies, fmt.Errorf("no api access cookie")
	}
	sessionCookies.CanAccessAPI = apiCookie.(bool)

	sessionToken := session.Get(sessionToken).(string)
	sessionCookies.SessionToken = sessionToken

	return sessionCookies, nil
}

func DumpSession(c *gin.Context) {
	session := sessions.Default(c)
	sessionCookies := SessionCookies{}
	userID := session.Get(userkey)
	if userID != nil {
		sessionCookies.UserID = userID.(int)
	}

	adminCookie := session.Get(isAdmin).(bool)
	sessionCookies.IsAdmin = adminCookie

	apiCookie := session.Get(canAccessAPI).(bool)
	sessionCookies.CanAccessAPI = apiCookie

	sessionToken := session.Get(sessionToken).(string)
	sessionCookies.SessionToken = sessionToken

	c.JSON(200, sessionCookies)
}

func TestSession(c *gin.Context) {
	c.JSON(200, "TBD")
}
