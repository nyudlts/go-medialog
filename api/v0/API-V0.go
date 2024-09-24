package api

import (
	"crypto/sha512"
	"encoding/hex"
	"fmt"
	"net/http"
	"runtime"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/nyudlts/go-medialog/controllers"
	"github.com/nyudlts/go-medialog/database"
	"github.com/nyudlts/go-medialog/models"
)

type EntryResultSet struct {
	FirstPage int            `json:"first_page"`
	LastPage  int            `json:"last_page"`
	ThisPage  int            `json:"this_page"`
	Total     int64          `json:"total"`
	Results   []models.Entry `json:"results"`
}

type SummaryTotalsRepo struct {
	Repository string             `json:"repository"`
	Totals     database.Totals    `json:"totals"`
	Summaries  []database.Summary `json:"summaries"`
}

type SummaryTotalsResource struct {
	ResourceIdentifier string             `json:"resource_identifier"`
	ResourceTitle      string             `json:"resource_title"`
	Totals             database.Totals    `json:"totals"`
	Summaries          []database.Summary `json:"summaries"`
}

type SummaryTotalsAccession struct {
	AccessionIdentifier string             `json:"accession_identifier"`
	Totals              database.Totals    `json:"totals"`
	Summaries           []database.Summary `json:"summaries"`
}

const UNAUTHORIZED = "Please authenticate to access this service"
const apiVersion = "v0.1.4"
const medialogVersion = "v1.0.7"

var ACCESS_DENIED = map[string]string{"error": "access denied"}

type APIError struct {
	Message map[string][]string `json:"error"`
}

func APILogin(c *gin.Context) {
	controllers.ExpireTokens()
	email := c.Param("user")
	password := c.Query("password")

	if password == "" {
		apiError := APIError{}
		e := map[string][]string{"password": []string{"Parameter required but no value provided"}}
		apiError.Message = e
		c.JSON(http.StatusBadRequest, apiError)
		return
	}

	user, err := database.FindUserByEmail(email)
	if err != nil {
		c.JSON(http.StatusUnauthorized, map[string]string{"error": "login failed - user not found"})
		return
	}

	hash := sha512.Sum512([]byte(password + user.Salt))
	userSHA512 := hex.EncodeToString(hash[:])

	if userSHA512 != user.EncryptedPassword {
		c.JSON(http.StatusBadRequest, fmt.Sprintf("storedChecksum: %s, calculatedChecksum: %s", user.EncryptedPassword, userSHA512))
		return
	}

	if !user.CanAccessAPI {
		c.JSON(http.StatusUnauthorized, "login failed -- user not authorized to access api")
		return
	}

	token := controllers.GenerateStringRunes(24)
	tkHash := sha512.Sum512([]byte(token))
	token = hex.EncodeToString(tkHash[:])

	user.EncryptedPassword = "####"
	user.Salt = "####"

	apiToken := models.Token{
		Token:   token,
		UserID:  user.ID,
		IsValid: true,
		Expires: time.Now().Add(time.Hour * 3),
		User:    user,
		Type:    "api",
	}

	//expire users other tokens
	if err := database.ExpireAPITokensByUserID(user.ID); err != nil {
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

func APILogout(c *gin.Context) {
	token, err := checkToken(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, ACCESS_DENIED)
		return
	}

	if err := database.DeleteToken(token); err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, fmt.Sprintf("Logged Out"))
}

func GetV0Root(c *gin.Context) {
	medialogInfo := models.MedialogInfo{
		Version:       medialogVersion,
		GolangVersion: runtime.Version(),
		GinVersion:    gin.Version,
		APIVersion:    apiVersion,
	}

	c.JSON(http.StatusOK, medialogInfo)
}

func DeleteSessionsV0(c *gin.Context) {
	_, err := checkToken(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, err.Error())
		return
	}

	if err := database.DeleteSessions(); err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
	}

	c.JSON(http.StatusOK, "sessions deleted")
}

func checkToken(c *gin.Context) (string, error) {
	controllers.ExpireTokens()
	token := c.Request.Header.Get("X-Medialog-Token")

	if token == "" {
		return "", fmt.Errorf("no `X-Medialog-Token` set in request header")
	}

	apiToken, err := database.FindToken(token)
	if err != nil {
		return "", fmt.Errorf("could not find supplied token: %s", token)
	}

	if !apiToken.IsValid {
		return "", fmt.Errorf("invalid token - please reauthenticate")
	}

	return token, nil
}
