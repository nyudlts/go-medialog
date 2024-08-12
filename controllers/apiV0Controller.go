package controllers

import (
	"crypto/sha512"
	"encoding/hex"
	"fmt"
	"net/http"
	"runtime"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/nyudlts/go-medialog/database"
	"github.com/nyudlts/go-medialog/models"
)

type MedialogInfo struct {
	Version       string
	GinVersion    string
	GolangVersion string
	APIVersion    string
}

const UNAUTHORIZED = "Please authenticate to access this service"

func TestAPI(c *gin.Context) {
	sessionCookies, err := getSessionCookies(c)
	if err != nil {
		c.JSON(403, UNAUTHORIZED)
		return
	}
	c.JSON(200, sessionCookies)
}

func APILogin(c *gin.Context) {
	expireTokens()
	email := c.Param("user")
	passwd := c.Query("password")

	user, err := database.FindUserByEmail(email)
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
	}

	hash := sha512.Sum512([]byte(passwd + user.Salt))
	userSHA512 := hex.EncodeToString(hash[:])

	if userSHA512 != user.EncryptedPassword {
		c.JSON(http.StatusBadRequest, "password was incorrect")
		return
	}

	if !user.CanAccessAPI {
		c.JSON(http.StatusUnauthorized, "not authorized to access api")
	}

	token := GenerateStringRunes(24)
	tkHash := sha512.Sum512([]byte(token))
	token = hex.EncodeToString(tkHash[:])

	user.EncryptedPassword = "####"
	user.Salt = "####"

	apiToken := models.Token{
		Token:   token,
		UserID:  user.ID,
		IsValid: true,
		Expires: time.Now().Add(time.Hour * 2),
		User:    user,
	}

	//expire users other tokens
	if err := database.ExpireTokensByUserID(user.ID); err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	//add token to api db
	if err := database.InsertToken(&apiToken); err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(200, apiToken)
}

func checkToken(c *gin.Context) error {
	expireTokens()
	token := c.Request.Header.Get("X-Medialog-Token")
	apiToken, err := database.FindToken(token)
	if err != nil {
		return fmt.Errorf("could not find supplied token: %s", token)
	}

	if !apiToken.IsValid {
		return fmt.Errorf("invalid token - please reauthenticate")
	}

	return nil

}

func GetV0Index(c *gin.Context) {
	medialogInfo := MedialogInfo{
		Version:       "v1.0.4",
		GolangVersion: runtime.Version(),
		GinVersion:    gin.Version,
		APIVersion:    "0.1.1",
	}

	c.JSON(http.StatusOK, medialogInfo)
}

func GetResourcesV0(c *gin.Context) {

	err := checkToken(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, err.Error())
		return
	}

	resources, err := database.FindResources()
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error)
		return
	}

	c.JSON(http.StatusOK, resources)
}

func GetResourceV0(c *gin.Context) {

	if err := checkToken(c); err != nil {
		c.JSON(http.StatusUnauthorized, err.Error())
		return
	}

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error)
	}

	resource, err := database.FindResource(uint(id))
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error)
		return
	}

	c.JSON(http.StatusOK, resource)
}

func GetRepositoriesV0(c *gin.Context) {

	if err := checkToken(c); err != nil {
		c.JSON(http.StatusUnauthorized, err.Error())
		return
	}

	repositories, err := database.FindRepositories()
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error)
		return
	}

	c.JSON(http.StatusOK, repositories)
}

func GetRepositoryV0(c *gin.Context) {

	if err := checkToken(c); err != nil {
		c.JSON(http.StatusUnauthorized, err.Error())
		return
	}

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error)
	}

	repository, err := database.FindRepository(uint(id))
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error)
		return
	}

	c.JSON(http.StatusOK, repository)
}

func GetAccessionsV0(c *gin.Context) {

	if err := checkToken(c); err != nil {
		c.JSON(http.StatusUnauthorized, err.Error())
		return
	}

	accessions, err := database.FindAccessions()
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error)
		return
	}

	c.JSON(http.StatusOK, accessions)
}

func GetAccessionV0(c *gin.Context) {

	if err := checkToken(c); err != nil {
		c.JSON(http.StatusUnauthorized, err.Error())
		return
	}

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error)
	}
	accession, err := database.FindAccession(uint(id))
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error)
		return
	}

	c.JSON(http.StatusOK, accession)
}
