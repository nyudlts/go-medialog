package controllers

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/nyudlts/go-medialog/database"
	"github.com/nyudlts/go-medialog/models"
)

func GetResource(c *gin.Context) {
	isLoggedIn := isLoggedIn(c)
	if !isLoggedIn {
		throwError(http.StatusUnauthorized, UNAUTHORIZED, c)
		return
	}

	isAdmin := getCookie("is-admin", c)

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	resource, err := database.FindResource(uint(id))
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	summary, err := database.GetSummaryByResource(resource.ID)
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	accessions, err := database.FindAccessionsByResourceID(resource.ID)
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	//get page count
	pageCount, err := database.GetNumberPagesInResource(resource.ID)
	if err != nil {
		throwError(http.StatusInternalServerError, err.Error(), c)
	}
	log.Println("PAGE COUNT", pageCount)

	//pagination
	var p = 0
	page := c.Request.URL.Query()["page"]

	if len(page) > 0 {
		p, err = strconv.Atoi(page[0])
		if err != nil {
			c.JSON(http.StatusBadRequest, err.Error())
			return
		}
	}

	if p < 0 {
		p = 0
	}

	pagination := database.Pagination{Limit: 10, Offset: (p * 10), Sort: "media_id"}

	entries, err := database.FindEntriesByResourceID(resource.ID, pagination)
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	entryUsers, err := database.FindEntryUsers(resource.CreatedBy, resource.UpdatedBy)
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	entryCount := database.GetCountOfEntriesInResource(resource.ID)

	c.HTML(http.StatusOK, "resources-show.html", gin.H{
		"resource":        resource,
		"accessions":      accessions,
		"entries":         entries,
		"isAdmin":         isAdmin,
		"isAuthenticated": true,
		"page":            p,
		"summary":         summary,
		"totals":          summary.GetTotals(),
		"entry_users":     entryUsers,
		"page_count":      pageCount,
		"entryCount":      entryCount,
		"isLoggedIn":      isLoggedIn,
	})
}

func GetResources(c *gin.Context) {
	isLoggedIn := isLoggedIn(c)
	if !isLoggedIn {
		throwError(http.StatusUnauthorized, UNAUTHORIZED, c)
		return
	}

	isAdmin := getCookie("is-admin", c)

	resources, err := database.FindResources()
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	repositoryMap, err := database.GetRepositoryMap()
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	c.HTML(http.StatusOK, "resources-index.html", gin.H{
		"resources":       resources,
		"isAuthenticated": true,
		"isAdmin":         isAdmin,
		"repositoryMap":   repositoryMap,
		"isLoggedIn":      isLoggedIn,
	})
}

func NewResource(c *gin.Context) {
	isLoggedIn := isLoggedIn(c)
	if !isLoggedIn {
		throwError(http.StatusUnauthorized, UNAUTHORIZED, c)
		return
	}

	isAdmin := getCookie("is-admin", c)

	repoID, err := strconv.Atoi(c.Query("repository_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	repository, err := database.FindRepository(uint(repoID))
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	c.HTML(200, "resources-new.html", gin.H{
		"isAdmin":    isAdmin,
		"repository": repository,
		"isLoggedIn": isLoggedIn,
	})
}

func CreateResource(c *gin.Context) {
	//ensure user is logged in
	isLoggedIn := isLoggedIn(c)
	if !isLoggedIn {
		throwError(http.StatusUnauthorized, UNAUTHORIZED, c)
		return
	}

	//bind the form to a resource
	var resource = models.Resource{}
	if err := c.Bind(&resource); err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	//get the repository
	repository, err := database.FindRepository(uint(resource.RepositoryID))
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	//set the repository
	resource.Repository = repository

	//get the current user id
	userID, err := getUserkey(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
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
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	//redirect to the new resource
	c.Redirect(302, fmt.Sprintf("/resources/%d/show", resourceID))
}

func EditResource(c *gin.Context) {
	isLoggedIn := isLoggedIn(c)
	if !isLoggedIn {
		throwError(http.StatusUnauthorized, UNAUTHORIZED, c)
		return
	}

	isAdmin := getCookie("is-admin", c)

	resourceID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	resource, err := database.FindResource(uint(resourceID))
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	c.HTML(200, "resources-edit.html", gin.H{
		"isAdmin":    isAdmin,
		"resource":   resource,
		"isLoggedIn": isLoggedIn,
	})
}

func UpdateResource(c *gin.Context) {
	isLoggedIn := isLoggedIn(c)
	if !isLoggedIn {
		throwError(http.StatusUnauthorized, UNAUTHORIZED, c)
		return
	}

	var updateResource = models.Resource{}
	if err := c.Bind(&updateResource); err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	userID, err := getUserkey(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	resource, err := database.FindResource(uint(id))
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	resource.UpdatedBy = userID
	resource.UpdatedAt = time.Now()
	resource.Title = updateResource.Title
	resource.CollectionCode = updateResource.CollectionCode

	if err := database.UpdateResource(&resource); err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	c.Redirect(302, fmt.Sprintf("/resources/%d/show", resource.ID))
}

func DeleteResource(c *gin.Context) {
	isLoggedIn := isLoggedIn(c)
	if !isLoggedIn {
		throwError(http.StatusUnauthorized, UNAUTHORIZED, c)
		return
	}

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	resource, err := database.FindResource(uint(id))
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	if err := database.DeleteResource(uint(id)); err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	c.Redirect(302, fmt.Sprintf("/repositories/%d/show", resource.RepositoryID))
}
