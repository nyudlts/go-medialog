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
	"github.com/nyudlts/go-medialog/version"
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

var ACCESS_DENIED = map[string]string{"error": "access denied"}

type APIError struct {
	Message map[string][]string `json:"error"`
}

// APILogin authenticates a user and returns an API token.
// @Summary      Login
// @Description  Authenticates a user by email and password, returning a session token valid for 3 hours.
// @Tags         auth
// @Produce      json
// @Param        user      path   string  true  "User email address"
// @Param        password  query  string  true  "User password"
// @Success      200  {object}  models.Token
// @Failure      400  {object}  APIError
// @Failure      401  {string}  string
// @Failure      500  {string}  string
// @Router       /users/{user}/login [post]
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

// APILogout invalidates the current API token.
// @Summary      Logout
// @Description  Invalidates the current API token.
// @Tags         auth
// @Produce      json
// @Security     ApiKeyAuth
// @Success      200  {string}  string
// @Failure      401  {object}  map[string]string
// @Failure      500  {string}  string
// @Router       /logout [delete]
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

	c.JSON(http.StatusOK, "Logged Out")
}

// GetV0Root returns application version information.
// @Summary      API root
// @Description  Returns version information about the application and API.
// @Tags         root
// @Produce      json
// @Success      200  {object}  models.MedialogInfo
// @Router       / [get]
func GetV0Root(c *gin.Context) {
	medialogInfo := models.MedialogInfo{
		Version:       version.AppVersion,
		GolangVersion: runtime.Version(),
		GinVersion:    gin.Version,
		APIVersion:    version.APIVersion,
	}

	c.JSON(http.StatusOK, medialogInfo)
}

// DeleteSessionsV0 deletes all active web sessions.
// @Summary      Delete all sessions
// @Description  Deletes all active web sessions from the database.
// @Tags         sessions
// @Produce      json
// @Security     ApiKeyAuth
// @Success      200  {string}  string
// @Failure      401  {string}  string
// @Failure      500  {string}  string
// @Router       /delete_sessions [delete]
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
