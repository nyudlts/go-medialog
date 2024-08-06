package controllers

import (
	"crypto/sha512"
	"encoding/hex"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/nyudlts/go-medialog/database"
	"github.com/nyudlts/go-medialog/models"
)

func GetUsers(c *gin.Context) {
	isLoggedIn := isLoggedIn(c)
	if !isLoggedIn {
		throwError(http.StatusUnauthorized, UNAUTHORIZED, c)
		return
	}

	sessionCookies, err := getSessionCookies(c)
	if err != nil {
		throwError(http.StatusUnauthorized, UNAUTHORIZED, c)
		return
	}

	if !sessionCookies.IsAdmin {
		throwError(http.StatusUnauthorized, "Must be logged in as an admin to access users management", c)
		return
	}

	user, err := database.GetRedactedUser(sessionCookies.UserID)
	if err != nil {
		throwError(http.StatusUnauthorized, UNAUTHORIZED, c)
		return
	}

	users, err := database.FindUsers()
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	c.HTML(http.StatusOK, "users-index.html", gin.H{
		"users":           users,
		"user":            user,
		"isAuthenticated": true,
		"isAdmin":         isAdmin,
		"isLoggedIn":      isLoggedIn,
	})
}

func GetUser(c *gin.Context) {
	isLoggedIn := isLoggedIn(c)
	if !isLoggedIn {
		throwError(http.StatusUnauthorized, UNAUTHORIZED, c)
		return
	}

	sessionCookies, err := getSessionCookies(c)
	if err != nil {
		throwError(http.StatusUnauthorized, UNAUTHORIZED, c)
		return
	}

	if !sessionCookies.IsAdmin {
		throwError(http.StatusUnauthorized, "Must be logged in as an admin to access users management", c)
		return
	}

	user, err := database.GetRedactedUser(sessionCookies.UserID)
	if err != nil {
		throwError(http.StatusUnauthorized, UNAUTHORIZED, c)
		return
	}

	uuserID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		throwError(http.StatusBadRequest, err.Error(), c)
		return
	}

	cookieId, err := getUserkey(c)
	if err != nil {
		throwError(http.StatusExpectationFailed, "User is not logged in / Unauthorized", c)
	}

	if (uuserID != cookieId) && !sessionCookies.IsAdmin {
		throwError(http.StatusUnauthorized, "Logged in as different user", c)
		return
	}

	uuser, err := database.GetRedactedUser(uuserID)
	if err != nil {
		throwError(http.StatusBadRequest, err.Error(), c)
		return
	}

	cookies, _ := getSessionCookies(c)

	c.HTML(200, "users-show.html", gin.H{
		"isLoggedIn": isLoggedIn,
		"isAdmin":    isAdmin,
		"uuser":      uuser,
		"user":       user,
		"cookies":    cookies,
	})
}

func NewUser(c *gin.Context) {
	isLoggedIn := isLoggedIn(c)
	if !isLoggedIn {
		throwError(http.StatusUnauthorized, UNAUTHORIZED, c)
		return
	}

	sessionCookies, err := getSessionCookies(c)
	if err != nil {
		throwError(http.StatusUnauthorized, UNAUTHORIZED, c)
		return
	}

	if !sessionCookies.IsAdmin {
		throwError(http.StatusUnauthorized, "Must be logged in as an admin to access users management", c)
		return
	}

	user, err := database.GetRedactedUser(sessionCookies.UserID)
	if err != nil {
		throwError(http.StatusBadRequest, err.Error(), c)
		return
	}

	c.HTML(http.StatusOK, "users-new.html", gin.H{
		"isAdmin":         sessionCookies.IsAdmin,
		"isAuthenticated": true,
		"isLoggedIn":      isLoggedIn,
		"user":            user,
	})
}

type UserForm struct {
	ID        int    `form:"id"`
	Password1 string `form:"password_1"`
	Password2 string `form:"password_2"`
	Email     string `form:"email"`
	FirstName string `form:"first_name"`
	LastName  string `form:"last_name"`
}

func CreateUser(c *gin.Context) {
	isLoggedIn := isLoggedIn(c)
	if !isLoggedIn {
		throwError(http.StatusUnauthorized, UNAUTHORIZED, c)
		return
	}

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
	user.FirstName = createUser.FirstName
	user.LastName = createUser.LastName
	user.IsActive = true
	user.Salt = GenerateStringRunes(16)
	hash := sha512.Sum512([]byte(createUser.Password1 + user.Salt))
	user.EncryptedPassword = hex.EncodeToString(hash[:])

	if _, err := database.InsertUser(&user); err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	c.Redirect(http.StatusFound, "/users")
}

func EditUser(c *gin.Context) {

	isLoggedIn := isLoggedIn(c)
	if !isLoggedIn {
		throwError(http.StatusUnauthorized, UNAUTHORIZED, c)
		return
	}

	sessionCookies, err := getSessionCookies(c)
	if err != nil {
		throwError(http.StatusInternalServerError, err.Error(), c)
		return
	}

	user, err := database.GetRedactedUser(sessionCookies.UserID)
	if err != nil {
		throwError(http.StatusBadRequest, err.Error(), c)
		return
	}

	userID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		throwError(http.StatusBadRequest, err.Error(), c)
		return
	}

	updateUser, err := database.GetRedactedUser(userID)
	if err != nil {
		throwError(http.StatusBadRequest, err.Error(), c)
		return
	}

	c.HTML(http.StatusOK, "users-edit.html", gin.H{
		"isLoggedIn": isLoggedIn,
		"isAdmin":    sessionCookies.IsAdmin,
		"user":       user,
		"updateUser": updateUser,
	})

}

func UpdateUser(c *gin.Context) {
	isLoggedIn := isLoggedIn(c)
	if !isLoggedIn {
		throwError(http.StatusUnauthorized, UNAUTHORIZED, c)
		return
	}

	var updateUser = UserForm{}
	if err := c.Bind(&updateUser); err != nil {
		throwError(http.StatusBadRequest, err.Error(), c)
		return
	}

	user, err := database.FindUser(uint(updateUser.ID))
	if err != nil {
		throwError(http.StatusBadRequest, err.Error(), c)
		return
	}

	user.Email = updateUser.Email
	user.FirstName = updateUser.FirstName
	user.LastName = updateUser.LastName

	if err := database.UpdateUser(&user); err != nil {
		throwError(http.StatusBadRequest, err.Error(), c)
		return
	}

	c.Redirect(http.StatusFound, fmt.Sprintf("/users/%d/show", user.ID))
}

func AuthenticateUser(c *gin.Context) {

	var authUser = UserForm{}
	if err := c.Bind(&authUser); err != nil {
		throwError(http.StatusBadRequest, err.Error(), c)
		return
	}

	user, err := database.FindUserByEmail(authUser.Email)
	if err != nil {
		throwError(http.StatusUnauthorized, err.Error(), c)
		return
	}

	if !user.IsActive {
		throwError(http.StatusUnauthorized, fmt.Sprintf("User %s is not active, contact a system administrator", user.Email), c)
	}

	hash := sha512.Sum512([]byte(authUser.Password1 + user.Salt))
	userSHA512 := hex.EncodeToString(hash[:])

	if userSHA512 != user.EncryptedPassword {
		throwError(http.StatusBadRequest, "password was incorrect", c)
		return
	}

	if err := login(int(user.ID), c); err != nil {
		throwError(http.StatusInternalServerError, "Failed to save session", c)
		return
	}

	if user.IsAdmin {
		setCookie("is-admin", true, c)
	} else {
		setCookie("is-admin", false, c)
	}

	if user.CanAccessAPI {
		setCookie("can-access-api", true, c)
	} else {
		setCookie("can-access-api", false, c)
	}

	sessionToken := GenerateStringRunes(24)
	setCookie("token", sessionToken, c)

	token := models.Token{
		Token:   sessionToken,
		UserID:  user.ID,
		Expires: time.Now().Add(time.Hour * 3),
		IsValid: true,
	}

	if err := database.InsertToken(&token); err != nil {
		throwError(http.StatusInternalServerError, "could not save session token", c)
	}

	user.SignInCount = user.SignInCount + 1
	user.PreviousIPAddress = user.CurrentIPAddress
	user.CurrentIPAddress = c.ClientIP()
	if err := database.UpdateUser(&user); err != nil {
		throwError(http.StatusInternalServerError, "failed to update user", c)
	}

	c.Redirect(http.StatusFound, "/")
}

func ResetUserPassword(c *gin.Context) {
	isLoggedIn := isLoggedIn(c)
	if !isLoggedIn {
		throwError(http.StatusUnauthorized, UNAUTHORIZED, c)
		return
	}

	isAdmin := getCookie("is-admin", c)

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		throwError(http.StatusBadRequest, err.Error(), c)
		return
	}

	user, err := database.FindUserByID(id)
	if err != nil {
		throwError(http.StatusBadRequest, err.Error(), c)
		return
	}

	c.HTML(http.StatusOK, "users-reset-password.html", gin.H{
		"user":            user,
		"isAdmin":         isAdmin,
		"isAuthenticated": true,
		"isLoggedIn":      isLoggedIn,
	})

}

func ResetPassword(c *gin.Context) {
	isLoggedIn := isLoggedIn(c)
	if !isLoggedIn {
		throwError(http.StatusUnauthorized, UNAUTHORIZED, c)
		return
	}
	var resetUser = UserForm{}
	if err := c.Bind(&resetUser); err != nil {
		throwError(http.StatusBadRequest, err.Error(), c)
		return
	}

	if resetUser.Password1 != resetUser.Password2 {
		throwError(http.StatusBadRequest, "passwords do not match", c)
		return
	}

	user, err := database.FindUserByID(resetUser.ID)
	if err != nil {
		throwError(http.StatusBadRequest, err.Error(), c)
		return
	}

	user.Salt = GenerateStringRunes(16)
	hash := sha512.Sum512([]byte(resetUser.Password1 + user.Salt))
	user.EncryptedPassword = hex.EncodeToString(hash[:])

	if err := database.UpdateUser(&user); err != nil {
		throwError(http.StatusBadRequest, err.Error(), c)
		return
	}

	c.Redirect(http.StatusFound, "/users")
}

func DeactivateUser(c *gin.Context) {
	isLoggedIn := isLoggedIn(c)
	if !isLoggedIn {
		throwError(http.StatusUnauthorized, UNAUTHORIZED, c)
		return
	}

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		throwError(http.StatusBadRequest, err.Error(), c)
		return
	}

	user, err := database.FindUserByID(id)
	if err != nil {
		throwError(http.StatusBadRequest, err.Error(), c)
		return
	}

	user.IsActive = false

	if err := database.UpdateUser(&user); err != nil {
		throwError(http.StatusBadRequest, err.Error(), c)
		return
	}

	c.Redirect(http.StatusFound, "/users")

}

func ReactivateUser(c *gin.Context) {
	isLoggedIn := isLoggedIn(c)
	if !isLoggedIn {
		throwError(http.StatusUnauthorized, UNAUTHORIZED, c)
		return
	}

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		throwError(http.StatusBadRequest, err.Error(), c)
		return
	}

	user, err := database.FindUserByID(id)
	if err != nil {
		throwError(http.StatusBadRequest, err.Error(), c)
		return
	}

	user.IsActive = true

	if err := database.UpdateUser(&user); err != nil {
		throwError(http.StatusBadRequest, err.Error(), c)
		return
	}

	c.Redirect(http.StatusFound, "/users")

}

func MakeUserAdmin(c *gin.Context) {
	isLoggedIn := isLoggedIn(c)
	if !isLoggedIn {
		throwError(http.StatusUnauthorized, UNAUTHORIZED, c)
		return
	}

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		throwError(http.StatusBadRequest, err.Error(), c)
		return
	}

	user, err := database.FindUserByID(id)
	if err != nil {
		throwError(http.StatusBadRequest, err.Error(), c)
		return
	}

	user.IsAdmin = true

	if err := database.UpdateUser(&user); err != nil {
		throwError(http.StatusBadRequest, err.Error(), c)
		return
	}

	c.Redirect(http.StatusFound, "/users")
}

func RemoveUserAdmin(c *gin.Context) {
	isLoggedIn := isLoggedIn(c)
	if !isLoggedIn {
		throwError(http.StatusUnauthorized, UNAUTHORIZED, c)
		return
	}

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		throwError(http.StatusBadRequest, err.Error(), c)
		return
	}

	user, err := database.FindUserByID(id)
	if err != nil {
		throwError(http.StatusBadRequest, err.Error(), c)
		return
	}

	user.IsAdmin = false

	if err := database.UpdateUser(&user); err != nil {
		throwError(http.StatusBadRequest, err.Error(), c)
		return
	}

	c.Redirect(http.StatusFound, "/users")
}

func LoginUser(c *gin.Context) { c.HTML(http.StatusOK, "users-login.html", gin.H{}) }

func LogoutUser(c *gin.Context) {
	isLoggedIn := isLoggedIn(c)
	if !isLoggedIn {
		throwError(http.StatusInternalServerError, "not currently logged in -- cannot log out", c)
		return
	}

	logout(c)
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

func DeleteUser(id uint) error {
	if err := database.DeleteUser(id); err != nil {
		return err
	}
	return nil
}
