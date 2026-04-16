package controllers

import (
	"encoding/csv"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/nyudlts/go-medialog/database"
	"github.com/nyudlts/go-medialog/models"
)

const EntriesShow = "/entries/%s/show"

func GetEntry(c *gin.Context) {
	sessionCookies := c.MustGet(ContextKeySessionCookies).(SessionCookies)
	user := c.MustGet(ContextKeyUser).(models.User)

	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		ThrowError(http.StatusBadRequest, err.Error(), c, true)
		return
	}

	entry, err := database.FindEntry(id)
	if err != nil {
		ThrowError(http.StatusBadRequest, err.Error(), c, true)
		return
	}

	maxMediaID := database.FindMaxMediaIDInResource(entry.ResourceID)

	accession, err := database.FindAccession(uint(entry.AccessionID))
	if err != nil {
		ThrowError(http.StatusBadRequest, err.Error(), c, true)
		return
	}

	resource, err := database.FindResource(accession.ResourceID)
	if err != nil {
		ThrowError(http.StatusBadRequest, err.Error(), c, true)
		return
	}

	repository, err := database.FindRepository(resource.RepositoryID)
	if err != nil {
		ThrowError(http.StatusBadRequest, err.Error(), c, true)
		return
	}

	entryUsers, err := database.FindEntryUsers(entry.CreatedBy, entry.UpdatedBy)
	if err != nil {
		ThrowError(http.StatusBadRequest, err.Error(), c, true)
		return
	}

	c.HTML(http.StatusOK, "entries-show.html", gin.H{
		"entry":            entry,
		"accession":        accession,
		"resource":         resource,
		"repository":       repository,
		"isAdmin":          sessionCookies.IsAdmin,
		"entryUsers":       entryUsers,
		"isLoggedIn": true,
		"maxMediaID":       maxMediaID,
		"interfaces":       getInterfaces(),
		"hddInterfaces":    getHDDInterfaces(),
		"imageFormats":     getImageFormats(),
		"imagingSuccess":   getImageSuccess(),
		"interpretSuccess": getInterpretSuccess(),
		"user":             user,
	})
}

func GetPreviousEntry(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		ThrowError(http.StatusBadRequest, err.Error(), c, true)
		return
	}

	entry, err := database.FindEntry(id)
	if err != nil {
		ThrowError(http.StatusBadRequest, err.Error(), c, true)
		return
	}

	prevEntryID, err := database.FindEntryByMediaIDAndCollectionID(entry.MediaID-1, entry.ResourceID)
	if err != nil {
		ThrowError(http.StatusBadRequest, err.Error(), c, true)
		return
	}

	c.Redirect(http.StatusFound, fmt.Sprintf(EntriesShow, prevEntryID))
}

func GetNextEntry(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		ThrowError(http.StatusBadRequest, err.Error(), c, true)
		return
	}

	entry, err := database.FindEntry(id)
	if err != nil {
		ThrowError(http.StatusBadRequest, err.Error(), c, true)
		return
	}

	prevEntryID, err := database.FindEntryByMediaIDAndCollectionID(entry.MediaID+1, entry.ResourceID)
	if err != nil {
		ThrowError(http.StatusBadRequest, err.Error(), c, true)
		return
	}

	c.Redirect(http.StatusFound, fmt.Sprintf(EntriesShow, prevEntryID))
}

func GetEntries(c *gin.Context) {
	sessionCookies := c.MustGet(ContextKeySessionCookies).(SessionCookies)
	user := c.MustGet(ContextKeyUser).(models.User)

	//pagination
	var p = 0
	var err error
	page := c.Request.URL.Query()["page"]

	if len(page) > 0 {
		p, err = strconv.Atoi(page[0])
		if err != nil {
			ThrowError(http.StatusBadRequest, err.Error(), c, true)
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
			ThrowError(http.StatusBadRequest, err.Error(), c, true)
			return
		}
	}

	pagination := database.Pagination{Limit: limit, Offset: (p * limit), Sort: "updated_at desc", Page: p}

	//get filter
	filter := c.Request.URL.Query()["filter"]
	if len(filter) > 0 {
		pagination.Filter = filter[0]
	}

	totalEntries := database.GetCountOfEntriesInDBPaginated(&pagination)
	pagination.TotalRecords = totalEntries
	totalPages := totalEntries / int64(pagination.Limit)
	if totalEntries%int64(pagination.Limit) > 0 {
		totalPages++
	}
	pagination.TotalPages = int(totalPages)
	overlimit := ((pagination.Page * pagination.Limit) + pagination.Limit) > int(totalEntries)

	//get entries
	entries, err := database.FindPaginatedEntries(pagination)
	if err != nil {
		ThrowError(http.StatusBadRequest, err.Error(), c, true)
		return
	}

	//get repositoryMap
	repositoryMap, err := database.GetRepositoryMap()
	if err != nil {
		ThrowError(http.StatusBadRequest, err.Error(), c, true)
		return
	}

	//return
	c.HTML(http.StatusOK, "entries-index.html", gin.H{
		"entries":       entries,
		"isAdmin":       sessionCookies.IsAdmin,
		"pagination":    pagination,
		"repositoryMap": repositoryMap,
		"isLoggedIn": true,
		"user":          user,
		"overlimit":     overlimit,
		"limitValues":   LimitValues,
		"mediatypes":    GetMediatypes(),
	})
}

func NewEntry(c *gin.Context) {
	sessionCookies := c.MustGet(ContextKeySessionCookies).(SessionCookies)
	user := c.MustGet(ContextKeyUser).(models.User)

	aID := c.Query("accession_id")
	accessionID, err := strconv.Atoi(aID)
	if err != nil {
		ThrowError(http.StatusBadRequest, err.Error(), c, true)
		return
	}

	log.Println("Accession:", accessionID)

	accession, err := database.FindAccession(uint(accessionID))
	if err != nil {
		ThrowError(http.StatusBadRequest, err.Error(), c, true)
		return
	}

	resource, err := database.FindResource(accession.ResourceID)
	if err != nil {
		fmt.Println("RESOURCE")
		ThrowError(http.StatusBadRequest, err.Error(), c, true)
		return
	}

	repository, err := database.FindRepository(resource.RepositoryID)
	if err != nil {
		fmt.Println("REPOSITORY")
		ThrowError(http.StatusBadRequest, err.Error(), c, true)
		return
	}

	mediaID, err := database.FindNextMediaCollectionInResource(resource.ID)
	if err != nil {
		fmt.Println("MEDIA")
		ThrowError(http.StatusBadRequest, err.Error(), c, true)
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
		"isLoggedIn": true,
		"user":                   user,
	})

}

func CreateEntry(c *gin.Context) {
	//bind form to entry
	var createEntry = models.Entry{}
	if err := c.Bind(&createEntry); err != nil {
		ThrowError(http.StatusBadRequest, err.Error(), c, true)
		return
	}

	//validate the form
	if err := createEntry.ValidateEntry(); err != nil {
		ThrowError(http.StatusBadRequest, err.Error(), c, true)
		return
	}

	//add the default status and storage location
	createEntry.Status = "es_to_be_processed"
	createEntry.Location = "sl_not_imaged"

	//check if media id is unique
	b, err := database.IsMediaIDUniqueInResource(createEntry.MediaID, createEntry.ResourceID)
	if err != nil {
		ThrowError(http.StatusBadRequest, err.Error(), c, true)
		return
	}

	if !b {
		ThrowError(http.StatusBadRequest, fmt.Sprintf("%d is not a unique ID in resource %d", createEntry.MediaID, createEntry.ResourceID), c, true)
		return
	}

	//get the user's id
	userID, err := getUserkey(c)
	if err != nil {
		ThrowError(http.StatusBadRequest, err.Error(), c, true)
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
		ThrowError(http.StatusBadRequest, err.Error(), c, true)
		return
	}

	createEntry.Accession = accession

	//get the resource
	resource, err := database.FindResource(accession.ResourceID)
	if err != nil {
		ThrowError(http.StatusBadRequest, err.Error(), c, true)
		return
	}

	createEntry.Resource = resource

	//get the repository
	repository, err := database.FindRepository(uint(accession.Resource.RepositoryID))
	if err != nil {
		ThrowError(http.StatusBadRequest, err.Error(), c, true)
		return
	}

	createEntry.Repository = repository

	//insert the entry
	if err := database.InsertEntry(&createEntry); err != nil {
		ThrowError(http.StatusBadRequest, err.Error(), c, true)
		return
	}

	//redirect
	c.Redirect(http.StatusFound, fmt.Sprintf("entries/%s/show", createEntry.ID.String()))
}

func DeleteEntry(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		ThrowError(http.StatusBadRequest, err.Error(), c, true)
		return
	}

	entry, err := database.FindEntry(id)
	if err != nil {
		ThrowError(http.StatusBadRequest, err.Error(), c, true)
		return
	}

	if err := database.DeleteEntry(id); err != nil {
		ThrowError(http.StatusBadRequest, err.Error(), c, true)
		return
	}

	c.Redirect(http.StatusFound, fmt.Sprintf("/accessions/%d/show", entry.AccessionID))

}

func EditEntry(c *gin.Context) {
	sessionCookies := c.MustGet(ContextKeySessionCookies).(SessionCookies)
	user := c.MustGet(ContextKeyUser).(models.User)

	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		ThrowError(http.StatusBadRequest, err.Error(), c, true)
		return
	}

	entry, err := database.FindEntry(id)
	if err != nil {
		ThrowError(http.StatusBadRequest, err.Error(), c, true)
		return
	}

	accession, err := database.FindAccession(entry.AccessionID)
	if err != nil {
		ThrowError(http.StatusBadRequest, err.Error(), c, true)
		return
	}

	resource, err := database.FindResource(uint(accession.ResourceID))
	if err != nil {
		ThrowError(http.StatusBadRequest, err.Error(), c, true)
		return
	}

	repository, err := database.FindRepository(uint(resource.RepositoryID))
	if err != nil {
		ThrowError(http.StatusBadRequest, err.Error(), c, true)
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
		"isLoggedIn": true,
		"user":                   user,
	})
}

func UpdateEntry(c *gin.Context) {
	//parse the id
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		ThrowError(http.StatusBadRequest, err.Error(), c, true)
		return
	}

	var editedEntry = models.Entry{}

	//bind form to entry
	if err := c.Bind(&editedEntry); err != nil {
		c.JSON(http.StatusBadRequest, fmt.Sprintf("%s, %s", "bind", err.Error()))
		return
	}

	log.Println("STATUS", editedEntry.Status)

	//validate the form
	if err := editedEntry.ValidateEntry(); err != nil {
		ThrowError(http.StatusBadRequest, err.Error(), c, true)
		return
	}

	//find the original entry
	entry, err := database.FindEntry(id)
	if err != nil {
		ThrowError(http.StatusBadRequest, err.Error(), c, true)
		return
	}

	// get the user id
	userID, err := getUserkey(c)
	if err != nil {
		ThrowError(http.StatusBadRequest, err.Error(), c, true)
		return
	}

	//updated user and timestamp
	entry.UpdatedBy = userID
	entry.UpdatedAt = time.Now()
	entry.UpdateEntry(editedEntry)

	//update the entry
	if err := database.UpdateEntry(&entry); err != nil {
		ThrowError(http.StatusBadRequest, err.Error(), c, true)
		return
	}

	c.Redirect(http.StatusFound, fmt.Sprintf(EntriesShow, entry.ID.String()))
}

func CloneEntry(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		ThrowError(http.StatusBadRequest, err.Error(), c, true)
		return
	}

	entry, err := database.FindEntry(id)
	if err != nil {
		ThrowError(http.StatusBadRequest, err.Error(), c, true)
		return
	}

	nextID, err := database.FindNextMediaCollectionInResource(uint(entry.ResourceID))
	if err != nil {
		ThrowError(http.StatusBadRequest, err.Error(), c, true)
		return
	}

	//generate a new uuid
	newUUID, err := uuid.NewUUID()
	if err != nil {
		ThrowError(http.StatusBadRequest, err.Error(), c, true)
		return
	}

	userID, err := getUserkey(c)
	if err != nil {
		ThrowError(http.StatusBadRequest, err.Error(), c, true)
		return
	}

	accession, err := database.FindAccession(entry.AccessionID)
	if err != nil {
		ThrowError(http.StatusBadRequest, err.Error(), c, true)
		return
	}

	resource, err := database.FindResource(accession.ResourceID)
	if err != nil {
		ThrowError(http.StatusBadRequest, err.Error(), c, true)
		return
	}

	repository, err := database.FindRepository(resource.RepositoryID)
	if err != nil {
		ThrowError(http.StatusBadRequest, err.Error(), c, true)
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

	if err := database.InsertEntry(&entry); err != nil {
		ThrowError(http.StatusBadRequest, err.Error(), c, true)
		return
	}

	c.Redirect(http.StatusFound, fmt.Sprintf(EntriesShow, newUUID.String()))

}

type FindEntryInResource struct {
	MediaID    int `form:"media_id" json:"media_id"`
	ResourceID int `form:"resource_id" json:"resource_id"`
}

func FindEntry(c *gin.Context) {
	findEntry := FindEntryInResource{}
	if err := c.Bind(&findEntry); err != nil {
		ThrowError(http.StatusBadRequest, err.Error(), c, true)
		return
	}

	id, err := database.FindEntryInResource(findEntry.ResourceID, findEntry.MediaID)
	if err != nil {
		ThrowError(http.StatusBadRequest, err.Error(), c, true)
		return
	}

	c.Redirect(http.StatusFound, fmt.Sprintf(EntriesShow, id))

}

func EntriesGenCSV(c *gin.Context) {
	//find the entries
	entries, err := database.FindEntriesFiltered(c.Request.URL.Query()["filter"][0])
	if err != nil {
		ThrowError(http.StatusBadRequest, err.Error(), c, true)
		return
	}

	//get the repository map
	repositoryMap, err := database.GetRepositoryMap()
	if err != nil {
		ThrowError(http.StatusInternalServerError, err.Error(), c, true)
		return
	}

	//get the resource map
	resourceMap, err := database.GetResourceMap()
	if err != nil {
		ThrowError(http.StatusInternalServerError, err.Error(), c, true)
		return
	}

	//get the accessions map
	accessionsMap, err := database.GetAccessionsMap()
	if err != nil {
		ThrowError(http.StatusInternalServerError, err.Error(), c, true)
		return
	}

	//generate the csv
	csvFileName := "medialog_entries.csv"
	csvBuffer := new(strings.Builder)
	var csvWriter = csv.NewWriter(csvBuffer)
	csvWriter.Write(models.CSVHeader)
	for _, entry := range entries {
		record := entry.ToCSVEntryResult().ToCSV()
		record[2] = GetMediaType(record[2])
		repoID, err := strconv.Atoi(record[7])
		if err != nil {
			ThrowError(http.StatusInternalServerError, err.Error(), c, true)
			return
		}
		record[7] = repositoryMap[repoID]

		resID, err := strconv.Atoi(record[8])
		if err != nil {
			ThrowError(http.StatusInternalServerError, err.Error(), c, true)
			return
		}
		record[8] = resourceMap[uint(resID)]

		accID, err := strconv.Atoi(record[9])
		if err != nil {
			ThrowError(http.StatusInternalServerError, err.Error(), c, true)
			return
		}
		record[9] = accessionsMap[uint(accID)]
		csvWriter.Write(record)
	}
	csvWriter.Flush()
	c.Header("content-type", "text/csv")
	c.Header("Content-Description", "File Transfer")
	c.Header("Content-Disposition", "attachment; filename="+csvFileName)
	c.Status(http.StatusOK)
	c.Writer.Write([]byte(csvBuffer.String()))
}
