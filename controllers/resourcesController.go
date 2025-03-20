package controllers

import (
	"encoding/csv"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/nyudlts/go-medialog/database"
	"github.com/nyudlts/go-medialog/models"
)

func GetResource(c *gin.Context) {
	//check if user is logged in
	if err := isLoggedIn(c); err != nil {
		ThrowError(http.StatusUnauthorized, err.Error(), c, false)
		return
	}
	isLoggedIn := true

	//get session cookies
	sessionCookies, err := getSessionCookies(c)
	if err != nil {
		ThrowError(http.StatusInternalServerError, err.Error(), c, isLoggedIn)
		return
	}

	//get the user
	user, err := database.GetRedactedUser(sessionCookies.UserID)
	if err != nil {
		ThrowError(http.StatusBadRequest, err.Error(), c, isLoggedIn)
		return
	}

	//get the resource
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		ThrowError(http.StatusBadRequest, err.Error(), c, isLoggedIn)
		return
	}

	resource, err := database.FindResource(uint(id))
	if err != nil {
		ThrowError(http.StatusBadRequest, err.Error(), c, isLoggedIn)
		return
	}

	//get the summary
	summary, err := database.GetSummaryByResource(resource.ID)
	if err != nil {
		ThrowError(http.StatusBadRequest, err.Error(), c, isLoggedIn)
		return
	}

	//get associacted accessions
	accessions, err := database.FindAccessionsByResourceID(resource.ID)
	if err != nil {
		ThrowError(http.StatusBadRequest, err.Error(), c, isLoggedIn)
		return
	}

	//pagination
	var p = 0
	page := c.Request.URL.Query()["page"]

	if len(page) > 0 {
		p, err = strconv.Atoi(page[0])
		if err != nil {
			ThrowError(http.StatusBadRequest, err.Error(), c, isLoggedIn)
			return
		}
	}

	if p < 0 {
		p = 0
	}

	var limit = 10
	l := c.Request.URL.Query()["limit"]
	if len(l) > 0 {
		limit, err = strconv.Atoi(l[0])
		if err != nil {
			ThrowError(http.StatusBadRequest, err.Error(), c, isLoggedIn)
			return
		}
	}

	pagination := database.Pagination{Limit: limit, Offset: (p * limit), Sort: "media_id", Page: p}

	//get filter
	filter := c.Request.URL.Query()["filter"]
	if len(filter) > 0 {
		pagination.Filter = filter[0]
	}

	pagination.TotalRecords = database.GetCountOfEntriesInResourcePaginated(resource.ID, pagination)
	totalPages := pagination.TotalRecords / int64(pagination.Limit)
	if pagination.TotalRecords%int64(pagination.Limit) > 0 {
		totalPages++
	}
	pagination.TotalPages = int(totalPages)
	overlimit := ((pagination.Page * pagination.Limit) + pagination.Limit) > int(pagination.TotalRecords)

	//get entries
	entries, err := database.FindEntriesByResourceIDPaginated(resource.ID, pagination)
	if err != nil {
		ThrowError(http.StatusBadRequest, err.Error(), c, isLoggedIn)
		return
	}

	entryUsers, err := database.FindEntryUsers(resource.CreatedBy, resource.UpdatedBy)
	if err != nil {
		ThrowError(http.StatusBadRequest, err.Error(), c, isLoggedIn)
		return
	}

	c.HTML(http.StatusOK, "resources-show.html", gin.H{
		"resource":        resource,
		"accessions":      accessions,
		"entries":         entries,
		"isAdmin":         sessionCookies.IsAdmin,
		"isAuthenticated": true,
		"pagination":      pagination,
		"summary":         summary,
		"totals":          summary.GetTotals(),
		"entry_users":     entryUsers,
		"isLoggedIn":      isLoggedIn,
		"user":            user,
		"overlimit":       overlimit,
		"limitValues":     LimitValues,
		"mediatypes":      GetMediatypes(),
	})
}

func GetResources(c *gin.Context) {
	if err := isLoggedIn(c); err != nil {
		ThrowError(http.StatusUnauthorized, err.Error(), c, false)
		return
	}

	isLoggedIn := true

	sessionCookies, err := getSessionCookies(c)
	if err != nil {
		ThrowError(http.StatusInternalServerError, err.Error(), c, isLoggedIn)
		return
	}

	user, err := database.GetRedactedUser(sessionCookies.UserID)
	if err != nil {
		ThrowError(http.StatusBadRequest, err.Error(), c, isLoggedIn)
		return
	}

	resources, err := database.FindResources()
	if err != nil {
		ThrowError(http.StatusBadRequest, err.Error(), c, isLoggedIn)
		return
	}

	repositoryMap, err := database.GetRepositoryMap()
	if err != nil {
		ThrowError(http.StatusBadRequest, err.Error(), c, isLoggedIn)
		return
	}

	c.HTML(http.StatusOK, "resources-index.html", gin.H{
		"resources":       resources,
		"isAuthenticated": true,
		"isAdmin":         sessionCookies.IsAdmin,
		"repositoryMap":   repositoryMap,
		"isLoggedIn":      isLoggedIn,
		"user":            user,
	})
}

func NewResource(c *gin.Context) {
	if err := isLoggedIn(c); err != nil {
		ThrowError(http.StatusUnauthorized, err.Error(), c, false)
		return
	}

	isLoggedIn := true

	sessionCookies, err := getSessionCookies(c)
	if err != nil {
		ThrowError(http.StatusInternalServerError, err.Error(), c, isLoggedIn)
		return
	}

	user, err := database.GetRedactedUser(sessionCookies.UserID)
	if err != nil {
		ThrowError(http.StatusBadRequest, err.Error(), c, isLoggedIn)
		return
	}

	repoID, err := strconv.Atoi(c.Query("repository_id"))
	if err != nil {
		ThrowError(http.StatusBadRequest, err.Error(), c, isLoggedIn)
		return
	}

	repository, err := database.FindRepository(uint(repoID))
	if err != nil {
		ThrowError(http.StatusBadRequest, err.Error(), c, isLoggedIn)
		return
	}

	c.HTML(200, "resources-new.html", gin.H{
		"isAdmin":    sessionCookies.IsAdmin,
		"repository": repository,
		"isLoggedIn": isLoggedIn,
		"user":       user,
	})
}

func CreateResource(c *gin.Context) {
	//ensure user is logged in
	if err := isLoggedIn(c); err != nil {
		ThrowError(http.StatusUnauthorized, err.Error(), c, false)
		return
	}

	isLoggedIn := true

	//bind the form to a resource
	var resource = models.Resource{}
	if err := c.Bind(&resource); err != nil {
		ThrowError(http.StatusBadRequest, err.Error(), c, isLoggedIn)
		return
	}

	//get the repository
	repository, err := database.FindRepository(uint(resource.RepositoryID))
	if err != nil {
		ThrowError(http.StatusBadRequest, err.Error(), c, isLoggedIn)
		return
	}

	//set the repository
	resource.Repository = repository

	//get the current user id
	userID, err := getUserkey(c)
	if err != nil {
		ThrowError(http.StatusBadRequest, err.Error(), c, isLoggedIn)
		return
	}

	//update the timestamps
	resource.CreatedAt = time.Now()
	resource.CreatedBy = userID
	resource.UpdatedAt = time.Now()
	resource.UpdatedBy = userID

	//insert the new resource
	resourceID, err := database.InsertResource(&resource)
	if err != nil {
		ThrowError(http.StatusBadRequest, err.Error(), c, isLoggedIn)
		return
	}

	//redirect to the new resource
	c.Redirect(http.StatusFound, fmt.Sprintf("/resources/%d/show", resourceID))
}

func EditResource(c *gin.Context) {
	if err := isLoggedIn(c); err != nil {
		ThrowError(http.StatusUnauthorized, err.Error(), c, false)
		return
	}

	isLoggedIn := true

	sessionCookies, err := getSessionCookies(c)
	if err != nil {
		ThrowError(http.StatusInternalServerError, err.Error(), c, isLoggedIn)
		return
	}

	user, err := database.GetRedactedUser(sessionCookies.UserID)
	if err != nil {
		ThrowError(http.StatusBadRequest, err.Error(), c, isLoggedIn)
		return
	}

	resourceID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		ThrowError(http.StatusBadRequest, err.Error(), c, isLoggedIn)
		return
	}

	resource, err := database.FindResource(uint(resourceID))
	if err != nil {
		ThrowError(http.StatusBadRequest, err.Error(), c, isLoggedIn)
		return
	}

	c.HTML(200, "resources-edit.html", gin.H{
		"isAdmin":    sessionCookies.IsAdmin,
		"resource":   resource,
		"isLoggedIn": isLoggedIn,
		"user":       user,
	})
}

func UpdateResource(c *gin.Context) {
	if err := isLoggedIn(c); err != nil {
		ThrowError(http.StatusUnauthorized, err.Error(), c, false)
		return
	}

	isLoggedIn := true

	var updateResource = models.Resource{}
	if err := c.Bind(&updateResource); err != nil {
		ThrowError(http.StatusBadRequest, err.Error(), c, isLoggedIn)
		return
	}

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		ThrowError(http.StatusBadRequest, err.Error(), c, isLoggedIn)
		return
	}

	userID, err := getUserkey(c)
	if err != nil {
		ThrowError(http.StatusBadRequest, err.Error(), c, isLoggedIn)
		return
	}

	resource, err := database.FindResource(uint(id))
	if err != nil {
		ThrowError(http.StatusBadRequest, err.Error(), c, isLoggedIn)
		return
	}

	resource.UpdatedBy = userID
	resource.UpdatedAt = time.Now()
	resource.Title = updateResource.Title
	resource.CollectionCode = updateResource.CollectionCode

	if err := database.UpdateResource(&resource); err != nil {
		ThrowError(http.StatusBadRequest, err.Error(), c, isLoggedIn)
		return
	}

	c.Redirect(http.StatusFound, fmt.Sprintf("/resources/%d/show", resource.ID))
}

func DeleteResource(c *gin.Context) {
	if err := isLoggedIn(c); err != nil {
		ThrowError(http.StatusUnauthorized, err.Error(), c, false)
		return
	}

	isLoggedIn := true
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		ThrowError(http.StatusBadRequest, err.Error(), c, isLoggedIn)
		return
	}

	resource, err := database.FindResource(uint(id))
	if err != nil {
		ThrowError(http.StatusBadRequest, err.Error(), c, isLoggedIn)
		return
	}

	if err := database.DeleteResource(uint(id)); err != nil {
		ThrowError(http.StatusBadRequest, err.Error(), c, isLoggedIn)
		return
	}

	c.Redirect(http.StatusFound, fmt.Sprintf("/repositories/%d/show", resource.RepositoryID))
}

func ResourceGenCSV(c *gin.Context) {

	if err := isLoggedIn(c); err != nil {
		ThrowError(http.StatusUnauthorized, err.Error(), c, false)
		return
	}
	isLoggedIn := true

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		ThrowError(http.StatusBadRequest, err.Error(), c, isLoggedIn)
		return
	}

	resource, err := database.FindResource(uint(id))
	if err != nil {
		ThrowError(http.StatusInternalServerError, err.Error(), c, isLoggedIn)
		return
	}

	entries, err := database.FindEntriesByResourceID(uint(id))
	if err != nil {
		ThrowError(http.StatusBadRequest, err.Error(), c, isLoggedIn)
		return
	}

	csvBuffer := new(strings.Builder)
	var csvWriter = csv.NewWriter(csvBuffer)
	csvWriter.Write(models.CSVHeader)
	for _, entry := range entries {
		record := entry.ToCSV()
		csvWriter.Write(record)
	}
	csvWriter.Flush()

	csvFileName := fmt.Sprintf("%s_%s.csv", resource.Repository.Slug, resource.CollectionCode)
	c.Header("content-type", "text/csv")
	c.Header("Content-Description", "File Transfer")
	c.Header("Content-Disposition", "attachment; filename="+csvFileName)
	c.Status(http.StatusOK)
	c.Writer.Write([]byte(csvBuffer.String()))
}
