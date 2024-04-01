package controllers

import (
	"crypto/sha512"
	"encoding/hex"
	"log"
	"math/rand"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/nyudlts/go-medialog/database"
	"github.com/nyudlts/go-medialog/models"
)

func GetUsers(c *gin.Context) {
	users, err := database.FindUsers()
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	c.HTML(http.StatusOK, "users-index.html", gin.H{
		"users": users,
	})
}

func NewUser(c *gin.Context) {
	c.HTML(http.StatusOK, "users-new.html", gin.H{})
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
		log.Printf("\t[ERROR]\t[DATABASE]\t%s", err.Error())
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}
	GetUsers(c)
}

func AuthenticateUser(c *gin.Context) {
	var authUser = UserForm{}
	if err := c.Bind(&authUser); err != nil {
		log.Printf("\t[ERROR]\t[MEDIALOG]\t%s", err.Error())
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

	if err := NewSession(user.ID, c); err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
	}

	log.Println(c.Request.Referer())
	c.Redirect(302, "/")
}

func LogoutUser(c *gin.Context) {
	removeSession(c)
	GetIndex(c)
}

func ResetUserPassword(c *gin.Context) {
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

	c.HTML(http.StatusOK, "users-reset-password.html", gin.H{"user": user})

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

	GetUsers(c)
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

	GetUsers(c)

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

	GetUsers(c)

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

	GetUsers(c)
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

	GetUsers(c)
}

func LoginUser(c *gin.Context) { c.HTML(http.StatusOK, "users-login.html", gin.H{}) }

var runes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890!@#$%^&*()_+{}[]:;<>,.?/")

func GenerateStringRunes(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = runes[rand.Intn(len(runes))]
	}
	return string(b)
}
