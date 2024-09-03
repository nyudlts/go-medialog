package controllers

import (
	"crypto/sha512"
	"encoding/hex"
	"fmt"
	"log"
	"net/http"
	"runtime"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
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
const apiVersion = "v0.1.3"
const medialogVersion = "v1.0.6"

var ACCESS_DENIED = map[string]string{"error": "access denied"}

type APIError struct {
	Message map[string][]string `json:"error"`
}

func APILogin(c *gin.Context) {
	expireTokens()
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

	token := GenerateStringRunes(24)
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

func GetV0Index(c *gin.Context) {
	medialogInfo := models.MedialogInfo{
		Version:       medialogVersion,
		GolangVersion: runtime.Version(),
		GinVersion:    gin.Version,
		APIVersion:    apiVersion,
	}

	c.JSON(http.StatusOK, medialogInfo)
}

/* Repository Functions */

func GetRepositoriesV0(c *gin.Context) {

	_, err := checkToken(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, ACCESS_DENIED)
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

	_, err := checkToken(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, ACCESS_DENIED)
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

func CreateRepositoryV0(c *gin.Context) {
	token, err := checkToken(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, ACCESS_DENIED)
		return
	}

	repo := models.Repository{}
	if err := c.Bind(&repo); err != nil {
		c.JSON(http.StatusBadRequest, err)
		return
	}

	userID, err := database.FindUserIDByToken(token)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	repo.CreatedBy = int(userID)
	repo.UpdatedBy = int(userID)
	repo.CreatedAt = time.Now()
	repo.UpdatedAt = time.Now()

	_, err = database.CreateRepository(&repo)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, repo)
}

func DeleteRepositoryV0(c *gin.Context) {
	_, err := checkToken(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, ACCESS_DENIED)
		return
	}

	repositoryIDParam := c.Param("id")
	repositoryID, err := strconv.Atoi(repositoryIDParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	if err := database.DeleteRepository(uint(repositoryID)); err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, fmt.Sprintf("Repository %d deleted", repositoryID))

}

//repository functions

func GetRepositoryEntriesV0(c *gin.Context) {
	_, err := checkToken(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, ACCESS_DENIED)
		return
	}

	repositoryIDParam := c.Param("id")
	repositoryID, err := strconv.Atoi(repositoryIDParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
	}

	allIDsParam := c.Query("all_ids")
	pageParam := c.Query("page")
	pageSizeParam := c.Query("page_size")
	log.Println(allIDsParam)
	var allIds bool
	if allIDsParam != "" {
		var err error
		allIds, err = strconv.ParseBool(allIDsParam)
		if err != nil {
			c.JSON(http.StatusBadRequest, err.Error())
			return
		}
	}

	if allIds {
		entries, err := database.FindEntryIDsByRepositoryID(uint(repositoryID))
		if err != nil {
			c.JSON(http.StatusInternalServerError, err.Error())
			return
		}
		c.JSON(http.StatusOK, entries)
		return
	} else {
		page, err := strconv.Atoi(pageParam)
		if err != nil {
			c.JSON(http.StatusBadRequest, err.Error())
			return
		}
		pagination := database.Pagination{}
		pagination.Offset = page

		pageSize, err := strconv.Atoi(pageSizeParam)
		if err != nil {
			pagination.Limit = 25
		} else {
			pagination.Limit = pageSize
		}

		entries, err := database.FindEntriesByRepositoryIDPaginated(uint(repositoryID), pagination)
		if err != nil {
			c.JSON(http.StatusInternalServerError, err.Error())
			return
		}

		e := EntryResultSet{}
		e.Total = database.GetCountOfEntriesInRepository(uint(repositoryID))
		e.FirstPage = 1
		e.ThisPage = page
		e.Results = entries
		r := int(e.Total / int64(pagination.Limit))
		m := int(e.Total % int64(pagination.Limit))
		var t int
		if m > 0 {
			t = r + 1
		}
		e.LastPage = t
		c.JSON(http.StatusOK, e)
	}
}

func GetRepositorySummaryV0(c *gin.Context) {
	_, err := checkToken(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, ACCESS_DENIED)
		return
	}

	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
	}

	repository, err := database.FindRepository(uint(id))
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	summaryMap, err := database.GetSummaryByRepository(uint(id))
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	summaryTotals := SummaryTotalsRepo{
		Repository: repository.Title,
		Totals:     summaryMap.GetTotals(),
		Summaries:  summaryMap.GetSlice(),
	}

	c.JSON(http.StatusOK, summaryTotals)
}

/* Resource Functions */
func CreateResourceV0(c *gin.Context) {
	token, err := checkToken(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, ACCESS_DENIED)
		return
	}

	resource := models.Resource{}
	if err := c.Bind(&resource); err != nil {
		c.JSON(http.StatusBadRequest, err)
		return
	}

	userID, err := database.FindUserIDByToken(token)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	repository, err := database.FindRepository(resource.RepositoryID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	resource.CreatedBy = int(userID)
	resource.UpdatedBy = int(userID)
	resource.CreatedAt = time.Now()
	resource.UpdatedAt = time.Now()
	resource.Repository = repository

	_, err = database.InsertResource(&resource)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, resource)
}

func DeleteResourceV0(c *gin.Context) {
	_, err := checkToken(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, ACCESS_DENIED)
		return
	}

	resourceIDParam := c.Param("id")
	resourceID, err := strconv.Atoi(resourceIDParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, err)
		return
	}

	if err := database.DeleteResource(uint(resourceID)); err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, fmt.Sprintf("Resource %d deleted", resourceID))

}

func GetResourcesV0(c *gin.Context) {

	_, err := checkToken(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, ACCESS_DENIED)
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

	_, err := checkToken(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, ACCESS_DENIED)
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

func GetResourceEntriesV0(c *gin.Context) {

	_, err := checkToken(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, ACCESS_DENIED)
		return
	}

	resourceIDParam := c.Param("id")
	resourceID, err := strconv.Atoi(resourceIDParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
	}

	allIDsParam := c.Query("all_ids")
	pageParam := c.Query("page")
	pageSizeParam := c.Query("page_size")

	var allIds bool
	if allIDsParam != "" {
		var err error
		allIds, err = strconv.ParseBool(allIDsParam)
		if err != nil {
			c.JSON(http.StatusBadRequest, err.Error())
			return
		}
	}

	if allIds {
		entries, err := database.FindEntryIDsByResourceID(uint(resourceID))
		if err != nil {
			c.JSON(http.StatusBadRequest, err.Error())
			return
		}
		c.JSON(http.StatusOK, entries)
	} else {
		page, err := strconv.Atoi(pageParam)
		if err != nil {
			c.JSON(http.StatusBadRequest, err.Error())
			return
		}
		pagination := database.Pagination{}
		pagination.Offset = page

		pageSize, err := strconv.Atoi(pageSizeParam)
		if err != nil {
			pagination.Limit = 25
		} else {
			pagination.Limit = pageSize
		}

		entries, err := database.FindEntriesByResourceIDPaginated(uint(resourceID), pagination)
		if err != nil {
			c.JSON(http.StatusInternalServerError, err.Error())
			return
		}

		e := EntryResultSet{}
		e.Total = database.GetCountOfEntriesInResource(uint(resourceID))
		e.FirstPage = 1
		e.ThisPage = page
		e.Results = entries
		r := int(e.Total / int64(pagination.Limit))
		m := int(e.Total % int64(pagination.Limit))
		var t int
		if m > 0 {
			t = r + 1
		}
		e.LastPage = t
		c.JSON(http.StatusOK, e)
	}
}

func GetResourceSummaryV0(c *gin.Context) {

	_, err := checkToken(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, ACCESS_DENIED)
		return
	}

	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	resource, err := database.FindResource(uint(id))
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	summaries, err := database.GetSummaryByResource(uint(id))
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	resourceSummary := SummaryTotalsResource{}
	resourceSummary.ResourceIdentifier = resource.CollectionCode
	resourceSummary.ResourceTitle = resource.Title
	resourceSummary.Totals = summaries.GetTotals()
	resourceSummary.Summaries = summaries.GetSlice()

	c.JSON(http.StatusOK, resourceSummary)
}

/* Accession Functions */

func CreateAccessionV0(c *gin.Context) {
	token, err := checkToken(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, ACCESS_DENIED)
		return
	}

	accession := models.Accession{}
	if err := c.Bind(&accession); err != nil {
		c.JSON(http.StatusBadRequest, err)
		return
	}

	userId, err := database.FindUserIDByToken(token)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	resource, err := database.FindResource(accession.ResourceID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	accession.CreatedBy = int(userId)
	accession.UpdatedBy = int(userId)
	accession.CreatedAt = time.Now()
	accession.UpdatedAt = time.Now()
	accession.Resource = resource

	_, err = database.InsertAccession(&accession)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, accession)

}

func GetAccessionsV0(c *gin.Context) {

	_, err := checkToken(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, ACCESS_DENIED)
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

	_, err := checkToken(c)
	if err != nil {
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

	repository, err := database.FindRepository(accession.Resource.RepositoryID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error)
		return
	}

	accession.Resource.Repository = repository

	c.JSON(http.StatusOK, accession)
}

func DeleteAccessionV0(c *gin.Context) {
	_, err := checkToken(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, err.Error())
		return
	}

	accessionIDParam := c.Param("id")
	accessionID, err := strconv.Atoi(accessionIDParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, err)
		return
	}

	if err := database.DeleteAccession(uint(accessionID)); err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, fmt.Sprintf("Resource %d deleted", accessionID))

}

func GetAccessionEntriesV0(c *gin.Context) {
	_, err := checkToken(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, ACCESS_DENIED)
		return
	}

	accessionIDParam := c.Param("id")
	accessionID, err := strconv.Atoi(accessionIDParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
	}

	allIDsParam := c.Query("all_ids")
	pageParam := c.Query("page")
	pageSizeParam := c.Query("page_size")
	var allIds bool

	if allIDsParam != "" {
		var err error
		allIds, err = strconv.ParseBool(allIDsParam)
		if err != nil {
			c.JSON(http.StatusBadRequest, err.Error())
			return
		}
	}

	if allIds {
		entries, err := database.FindEntryIDsByAccessionID(uint(accessionID))
		if err != nil {
			c.JSON(http.StatusBadRequest, err.Error())
			return
		}
		c.JSON(http.StatusOK, entries)
	} else {
		page, err := strconv.Atoi(pageParam)
		if err != nil {
			c.JSON(http.StatusBadRequest, err.Error())
			return
		}
		pagination := database.Pagination{}
		pagination.Offset = page

		pageSize, err := strconv.Atoi(pageSizeParam)
		if err != nil {
			pagination.Limit = 25
		} else {
			pagination.Limit = pageSize
		}

		entries, err := database.FindEntriesByAccessionID(uint(accessionID), pagination)
		if err != nil {
			c.JSON(http.StatusInternalServerError, err.Error())
			return
		}

		e := EntryResultSet{}
		e.Total = database.GetCountOfEntriesInAccession(uint(accessionID))
		e.FirstPage = 1
		e.ThisPage = page
		e.Results = entries
		r := int(e.Total / int64(pagination.Limit))
		m := int(e.Total % int64(pagination.Limit))
		var t int
		if m > 0 {
			t = r + 1
		}
		e.LastPage = t

		c.JSON(http.StatusOK, e)
	}

}

func GetAccessionSummaryV0(c *gin.Context) {

	_, err := checkToken(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, ACCESS_DENIED)
		return
	}

	idParam := c.Param("id")
	accessionID, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	accession, err := database.FindAccession(uint(accessionID))
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	summaries, err := database.GetSummaryByAccession(uint(accessionID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	summaryAccession := SummaryTotalsAccession{}
	summaryAccession.AccessionIdentifier = accession.AccessionNum
	summaryAccession.Totals = summaries.GetTotals()
	summaryAccession.Summaries = summaries.GetSlice()

	c.JSON(http.StatusOK, summaryAccession)

}

/* Entry Functions */

func GetEntryV0(c *gin.Context) {

	_, err := checkToken(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, err.Error())
		return
	}

	id := c.Param("id")

	uId, err := uuid.Parse(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	entry, err := database.FindEntry(uId)
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	c.JSON(http.StatusOK, entry)
}

func GetEntriesV0(c *gin.Context) {

	_, err := checkToken(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, ACCESS_DENIED)
		return
	}

	allIDsParam := c.Query("all_ids")

	var allIds bool
	if allIDsParam != "" {
		var err error
		allIds, err = strconv.ParseBool(allIDsParam)
		if err != nil {
			c.JSON(http.StatusBadRequest, err.Error())
			return
		}
	}

	if allIds {
		ids, err := database.GetEntryIDs()
		if err != nil {
			c.JSON(http.StatusBadRequest, err.Error())
			return
		} else {
			c.JSON(http.StatusOK, ids)
			return
		}
	} else {

		pageSizeParam := c.Query("page_size")
		var pageSize int
		if pageSizeParam != "" {
			var err error
			pageSize, err = strconv.Atoi(pageSizeParam)
			if err != nil {
				c.JSON(http.StatusBadRequest, err.Error())
				return
			}
		} else {
			pageSize = 25
		}

		pageParam := c.Query("page")
		var entries []models.Entry
		var page int
		if pageParam != "" {
			var err error
			page, err = strconv.Atoi(pageParam)
			if err != nil {
				c.JSON(http.StatusBadRequest, err.Error())
				return
			}

			pagination := database.Pagination{Offset: page, Limit: pageSize}
			fmt.Println(pagination)
			entries, err = database.FindEntriesPaginated(pagination)
			if err != nil {
				c.JSON(http.StatusBadRequest, err.Error())
				return
			}
		}

		results := EntryResultSet{}
		results.Total = database.GetCountOfEntriesInDB()
		r := int(results.Total / int64(pageSize))
		m := int(results.Total % int64(pageSize))
		var t int
		if m > 0 {
			t = r + 1
		}
		results.Results = entries
		results.FirstPage = 1
		results.ThisPage = page
		results.LastPage = t

		c.JSON(http.StatusOK, results)
		return
	}
}

func UpdateEntryLocationV0(c *gin.Context) {
	token, err := checkToken(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, err.Error())
		return
	}

	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, "no id provided")
		return
	}

	uid, err := uuid.Parse(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, "provided id is not a valid uuid")
		return
	}

	location := c.Query("location")
	if location == "" {
		c.JSON(http.StatusBadRequest, "no location provided")
		return
	}

	storageLocation := GetStorageLocation(location)

	if storageLocation == "No Match" {
		c.JSON(http.StatusBadRequest, fmt.Sprintf("`%s` is not a valid location", location))
		return
	}

	entry, err := database.FindEntry(uid)
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	userID, err := database.FindUserIDByToken(token)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	entry.Location = location
	entry.UpdatedAt = time.Now()
	entry.UpdatedBy = int(userID)

	if err := database.UpdateEntry(&entry); err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, fmt.Sprintf("id: %s, location: %s, storage location: %s", id, location, storageLocation))

}

func checkToken(c *gin.Context) (string, error) {
	expireTokens()
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
