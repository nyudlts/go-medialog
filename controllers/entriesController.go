package controllers

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/nyudlts/go-medialog/database"
	"github.com/nyudlts/go-medialog/models"
)

func GetEntry(c *gin.Context) {
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

	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	entry, err := database.FindEntry(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	maxMediaID := database.FindMaxMediaIDInResource(entry.ResourceID)

	accession, err := database.FindAccession(uint(entry.AccessionID))
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	resource, err := database.FindResource(accession.ResourceID)
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	repository, err := database.FindRepository(resource.RepositoryID)
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	entryUsers, err := database.FindEntryUsers(entry.CreatedBy, entry.UpdatedBy)
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	c.HTML(http.StatusOK, "entries-show.html", gin.H{
		"entry":      entry,
		"accession":  accession,
		"resource":   resource,
		"repository": repository,
		"isAdmin":    sessionCookies.IsAdmin,
		"entryUsers": entryUsers,
		"isLoggedIn": isLoggedIn,
		"maxMediaID": maxMediaID,
		"user":       user,
	})
}

func GetPreviousEntry(c *gin.Context) {
	if !isLoggedIn(c) {
		throwError(http.StatusUnauthorized, UNAUTHORIZED, c)
		return
	}

	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	entry, err := database.FindEntry(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	prevEntryID, err := database.FindEntryByMediaIDAndCollectionID(entry.MediaID-1, entry.ResourceID)
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	c.Redirect(http.StatusFound, fmt.Sprintf("/entries/%s/show", prevEntryID))
}

func GetNextEntry(c *gin.Context) {
	if !isLoggedIn(c) {
		throwError(http.StatusUnauthorized, UNAUTHORIZED, c)
		return
	}

	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	entry, err := database.FindEntry(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	prevEntryID, err := database.FindEntryByMediaIDAndCollectionID(entry.MediaID+1, entry.ResourceID)
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	c.Redirect(http.StatusFound, fmt.Sprintf("/entries/%s/show", prevEntryID))
}

func GetEntries(c *gin.Context) {
	if !isLoggedIn(c) {
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

	pagination := database.Pagination{Limit: 10, Offset: (p * 10), Sort: "updated_at desc"}

	entries, err := database.FindPaginatedEntries(pagination)
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	entryCount := database.GetCountOfEntriesInDB()

	repositoryMap, err := database.GetRepositoryMap()
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	c.HTML(http.StatusOK, "entries-index.html", gin.H{
		"entries":       entries,
		"isAdmin":       sessionCookies.IsAdmin,
		"page":          p,
		"repositoryMap": repositoryMap,
		"entryCount":    entryCount,
		"isLoggedIn":    isLoggedIn,
		"user":          user,
	})
}

func NewEntry(c *gin.Context) {
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

	aID := c.Query("accession_id")
	accessionID, err := strconv.Atoi(aID)
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	log.Println("Accession:", accessionID)

	accession, err := database.FindAccession(uint(accessionID))
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
	}

	resource, err := database.FindResource(accession.ResourceID)
	if err != nil {
		fmt.Println("RESOURCE")
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	repository, err := database.FindRepository(resource.RepositoryID)
	if err != nil {
		fmt.Println("REPOSITORY")
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	mediaID, err := database.FindNextMediaCollectionInResource(resource.ID)
	if err != nil {
		fmt.Println("MEDIA")
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	c.HTML(http.StatusOK, "entries-create.html", gin.H{
		"isAdmin":                sessionCookies.IsAdmin,
		"accession":              accession,
		"resource":               resource,
		"repository":             repository,
		"interfaces":             getInterfaces(),
		"stock_units":            getStockUnits(),
		"optical_content_types":  getOpticalContentTypes(),
		"hdd_interfaces":         getHDDInterfaces(),
		"imaging_success":        getImageSuccess(),
		"interpretation_success": getInterpretSuccess(),
		"imaging_software":       getImagingSoftware(),
		"image_formats":          getImageFormats(),
		"media_id":               mediaID,
		"is_refreshed":           is_refreshed,
		"isLoggedIn":             isLoggedIn,
		"user":                   user,
	})

}

func CreateEntry(c *gin.Context) {
	//check user is logged in
	if !isLoggedIn(c) {
		throwError(http.StatusUnauthorized, UNAUTHORIZED, c)
		return
	}

	//bind form to entry
	var createEntry = models.Entry{}
	if err := c.Bind(&createEntry); err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	//validate the form
	if err := createEntry.ValidateEntry(); err != nil {
		throwError(http.StatusBadRequest, err.Error(), c)
		return
	}

	//check if media id is unique
	b, err := database.IsMediaIDUniqueInResource(createEntry.MediaID, createEntry.ResourceID)
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	if !b {
		c.JSON(http.StatusBadRequest, fmt.Errorf("%d is not a unique ID in resource %d", createEntry.MediaID, createEntry.ResourceID))
		return
	}

	//get the user's id
	userID, err := getUserkey(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	createEntry.ID, _ = uuid.NewUUID()
	createEntry.CreatedAt = time.Now()
	createEntry.CreatedBy = userID
	createEntry.UpdatedAt = time.Now()
	createEntry.UpdatedBy = userID

	//get the accession
	accession, err := database.FindAccession(uint(createEntry.AccessionID))
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	createEntry.Accession = accession

	//get the resource
	resource, err := database.FindResource(accession.ResourceID)
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	createEntry.Resource = resource

	//get the repository
	repository, err := database.FindRepository(uint(accession.Resource.RepositoryID))
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	createEntry.Repository = repository

	//insert the entry
	if _, err := database.InsertEntry(&createEntry); err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	//redirect
	c.Redirect(http.StatusFound, fmt.Sprintf("entries/%s/show", createEntry.ID.String()))
}

func DeleteEntry(c *gin.Context) {
	if !isLoggedIn(c) {
		throwError(http.StatusUnauthorized, UNAUTHORIZED, c)
		return
	}

	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	entry, err := database.FindEntry(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	if err := database.DeleteEntry(id); err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	c.Redirect(http.StatusFound, fmt.Sprintf("/accessions/%d/show", entry.AccessionID))

}

func EditEntry(c *gin.Context) {
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

	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	entry, err := database.FindEntry(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	accession, err := database.FindAccession(entry.AccessionID)
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	resource, err := database.FindResource(uint(accession.ResourceID))
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	repository, err := database.FindRepository(uint(resource.RepositoryID))
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	c.HTML(http.StatusOK, "entries-edit.html", gin.H{
		"isAdmin":                sessionCookies.IsAdmin,
		"entry":                  entry,
		"accession":              accession,
		"resource":               resource,
		"repository":             repository,
		"interfaces":             getInterfaces(),
		"stock_units":            getStockUnits(),
		"optical_content_types":  getOpticalContentTypes(),
		"hdd_interfaces":         getHDDInterfaces(),
		"imaging_success":        getImageSuccess(),
		"interpretation_success": getInterpretSuccess(),
		"imaging_software":       getImagingSoftware(),
		"image_formats":          getImageFormats(),
		"is_refreshed":           is_refreshed,
		"isLoggedIn":             isLoggedIn,
		"user":                   user,
	})
}

func UpdateEntry(c *gin.Context) {
	//check for login
	if !isLoggedIn(c) {
		throwError(http.StatusUnauthorized, UNAUTHORIZED, c)
		return
	}

	//parse the id
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	var editedEntry = models.Entry{}

	//bind form to entry
	if err := c.Bind(&editedEntry); err != nil {
		c.JSON(http.StatusBadRequest, fmt.Sprintf("%s, %s", "bind", err.Error()))
		return
	}

	//validate the form
	if err := editedEntry.ValidateEntry(); err != nil {
		throwError(http.StatusBadRequest, err.Error(), c)
		return
	}

	//find the original entry
	entry, err := database.FindEntry(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, fmt.Sprintf("%s, %s", "find entry", err.Error()))
		return
	}

	// get the user id
	userID, err := getUserkey(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, err)
	}

	//updated user and timestamp
	entry.UpdatedBy = userID
	entry.UpdatedAt = time.Now()
	entry.UpdateEntry(editedEntry)

	//update the entry
	if err := database.UpdateEntry(&entry); err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	c.Redirect(http.StatusFound, fmt.Sprintf("/entries/%s/show", entry.ID.String()))
}

func CloneEntry(c *gin.Context) {

	//check login
	if !isLoggedIn(c) {
		throwError(http.StatusUnauthorized, UNAUTHORIZED, c)
		return
	}

	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	entry, err := database.FindEntry(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	nextID, err := database.FindNextMediaCollectionInResource(uint(entry.ResourceID))
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	//generate a new uuid
	newUUID, err := uuid.NewUUID()
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	userID, err := getUserkey(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	accession, err := database.FindAccession(entry.AccessionID)
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	resource, err := database.FindResource(accession.ResourceID)
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	repository, err := database.FindRepository(resource.RepositoryID)
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	entry.ID = newUUID
	entry.LabelText = ""
	entry.MediaID = nextID
	entry.CreatedAt = time.Now()
	entry.CreatedBy = userID
	entry.UpdatedAt = time.Now()
	entry.UpdatedBy = userID
	entry.AccessionID = accession.ID
	entry.Accession = accession
	entry.ResourceID = resource.ID
	entry.Resource = resource
	entry.RepositoryID = repository.ID
	entry.Repository = repository

	if _, err := database.InsertEntry(&entry); err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	c.Redirect(http.StatusFound, fmt.Sprintf("/entries/%s/show", newUUID.String()))

}

type FindEntryInResource struct {
	MediaID    int `form:"media_id" json:"media_id"`
	ResourceID int `form:"resource_id" json:"resource_id"`
}

func FindEntry(c *gin.Context) {
	if !isLoggedIn(c) {
		throwError(http.StatusUnauthorized, UNAUTHORIZED, c)
		return
	}

	findEntry := FindEntryInResource{}
	if err := c.Bind(&findEntry); err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	id, err := database.FindEntryInResource(findEntry.ResourceID, findEntry.MediaID)
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	c.Redirect(http.StatusFound, fmt.Sprintf("/entries/%s/show", id))

}
