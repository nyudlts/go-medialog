package controllers

import (
	"crypto/sha512"
	"encoding/hex"
	"log"
	"math/rand"
	"net/http"
	"strconv"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/nyudlts/go-medialog/database"
	"github.com/nyudlts/go-medialog/models"
)

func GetUsers(c *gin.Context) {
	if err := checkSession(c); err != nil {
		c.Redirect(302, "/")
		return
	}
	isAdmin := getCookie("is-admin", c)

	users, err := database.FindUsers()
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	c.HTML(http.StatusOK, "users-index.html", gin.H{
		"users":           users,
		"isAuthenticated": true,
		"isAdmin":         isAdmin,
	})
}

func NewUser(c *gin.Context) {
	if err := checkSession(c); err != nil {
		c.Redirect(302, "/")
		return
	}
	isAdmin := getCookie("is-admin", c)

	c.HTML(http.StatusOK, "users-new.html", gin.H{
		"isAdmin":         isAdmin,
		"isAuthenticated": true,
	})
}

type UserForm struct {
	ID        int    `form:"id"`
	Password1 string `form:"password_1"`
	Password2 string `form:"password_2"`
	Email     string `form:"email"`
}

func CreateUser(c *gin.Context) {
	var createUser = UserForm{}
	if err := c.Bind(&createUser); err != nil {
		log.Printf("\t[ERROR]\t[MEDIALOG] %s", err.Error())
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	if createUser.Password1 != createUser.Password2 {
		c.JSON(http.StatusBadRequest, "passwords do not match")
		return
	}

	user := models.User{}
	user.Email = createUser.Email
	user.Salt = GenerateStringRunes(16)
	hash := sha512.Sum512([]byte(createUser.Password1 + user.Salt))
	user.EncryptedPassword = hex.EncodeToString(hash[:])

	if err := database.InsertUser(&user); err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	c.Redirect(http.StatusMovedPermanently, "/users")
}

func AuthenticateUser(c *gin.Context) {

	var authUser = UserForm{}
	if err := c.Bind(&authUser); err != nil {
		log.Println(err.Error())
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	user, err := database.FindUserByEmail(authUser.Email)
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	hash := sha512.Sum512([]byte(authUser.Password1 + user.Salt))
	userSHA512 := hex.EncodeToString(hash[:])

	if userSHA512 != user.EncryptedPassword {
		c.JSON(http.StatusBadRequest, "password was incorrect")
		return
	}

	if err := login(int(user.ID), c); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save session"})
		return
	}
	log.Printf("[INFO] Successfully authenticated user")

	if user.IsAdmin {
		setCookie("is-admin", true, c)
	} else {
		setCookie("is-admin", false, c)
	}

	c.Redirect(http.StatusMovedPermanently, "/")
}

func ResetUserPassword(c *gin.Context) {
	if err := checkSession(c); err != nil {
		c.Redirect(302, "/")
		return
	}
	isAdmin := getCookie("is-admin", c)

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	user, err := database.FindUserByID(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	c.HTML(http.StatusOK, "users-reset-password.html", gin.H{
		"user":            user,
		"isAdmin":         isAdmin,
		"isAuthenticated": true,
	})

}

func ResetPassword(c *gin.Context) {
	var resetUser = UserForm{}
	if err := c.Bind(&resetUser); err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	if resetUser.Password1 != resetUser.Password2 {
		c.JSON(http.StatusBadRequest, "passwords do not match")
		return
	}

	user, err := database.FindUserByID(resetUser.ID)
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	user.Salt = GenerateStringRunes(16)
	hash := sha512.Sum512([]byte(resetUser.Password1 + user.Salt))
	user.EncryptedPassword = hex.EncodeToString(hash[:])

	if err := database.UpdateUser(user); err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	c.Redirect(http.StatusPermanentRedirect, "/users")
}

func DeactivateUser(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	user, err := database.FindUserByID(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	user.IsActive = false

	if err := database.UpdateUser(user); err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	c.Redirect(http.StatusPermanentRedirect, "/users")

}

func ReactivateUser(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	user, err := database.FindUserByID(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	user.IsActive = true

	if err := database.UpdateUser(user); err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	c.Redirect(http.StatusPermanentRedirect, "/users")

}

func MakeUserAdmin(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	user, err := database.FindUserByID(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	user.IsAdmin = true

	if err := database.UpdateUser(user); err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	c.Redirect(http.StatusPermanentRedirect, "/users")
}

func RemoveUserAdmin(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	user, err := database.FindUserByID(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	user.IsAdmin = false

	if err := database.UpdateUser(user); err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	c.Redirect(http.StatusPermanentRedirect, "/users")
}

func LoginUser(c *gin.Context) { c.HTML(http.StatusOK, "users-login.html", gin.H{}) }

func LogoutUser(c *gin.Context) {
	session := sessions.Default(c)
	user := session.Get(userkey)
	if user == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid session token"})
		return
	}
	session.Delete(userkey)
	if err := session.Save(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save session"})
		return
	}
	log.Printf("[INFO] Successfully authenticated user")

	c.Redirect(302, "/")
}

var runes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890!@#$%^&*()_+{}[]:;<>,.?/")

func GenerateStringRunes(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = runes[rand.Intn(len(runes))]
	}
	return string(b)
}

func getUserEmailMap(ids []int) (map[int]string, error) {
	users := map[int]string{}
	for _, id := range ids {
		if id == 0 {
			users[id] = "unknown"
		} else {
			email, err := database.FindUserEmailByID(id)
			if err != nil {
				return users, err
			}
			users[id] = email
		}
	}
	return users, nil
}
