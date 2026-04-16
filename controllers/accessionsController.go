package controllers

import (
	"encoding/csv"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/nyudlts/go-medialog/database"
	"github.com/nyudlts/go-medialog/models"
)

const AccessionsShow = "/accessions/%d/show"

func GetAccessions(c *gin.Context) {
	sessionCookies := c.MustGet(ContextKeySessionCookies).(SessionCookies)
	user := c.MustGet(ContextKeyUser).(models.User)

	accessions, err := database.FindAccessions()
	if err != nil {
		ThrowError(http.StatusBadRequest, err.Error(), c, true)
		return
	}

	repositoryMap, err := database.GetRepositoryMap()
	if err != nil {
		ThrowError(http.StatusBadRequest, err.Error(), c, true)
		return
	}

	repositoryMap2 := map[uint]string{}
	for k, v := range repositoryMap {
		repositoryMap2[uint(k)] = v
	}

	c.HTML(200, "accessions-index.html", gin.H{
		"accessions":    accessions,
		"isAdmin":       sessionCookies.IsAdmin,
		"repositoryMap": repositoryMap2,
		"isLoggedIn":    true,
		"user":          user,
	})
}

func GetAccession(c *gin.Context) {
	sessionCookies := c.MustGet(ContextKeySessionCookies).(SessionCookies)
	user := c.MustGet(ContextKeyUser).(models.User)

	//get accession
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		ThrowError(http.StatusBadRequest, err.Error(), c, true)
		return
	}

	accession, err := database.FindAccession(uint(id))
	if err != nil {
		ThrowError(http.StatusBadRequest, err.Error(), c, true)
		return
	}

	//setup pagination
	var p = 0
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

	pagination := database.Pagination{Limit: limit, Offset: (p * limit), Sort: "media_id", Page: p}

	//get filter
	filter := c.Request.URL.Query()["filter"]
	if len(filter) > 0 {
		pagination.Filter = filter[0]
	}

	totalEntries := database.GetCountOfEntriesInAccessionPaginated(accession.ID, &pagination)
	pagination.TotalRecords = totalEntries
	totalPages := totalEntries / int64(pagination.Limit)
	if totalEntries%int64(pagination.Limit) > 0 {
		totalPages++
	}
	pagination.TotalPages = int(totalPages)

	overlimit := ((pagination.Page * pagination.Limit) + pagination.Limit) > int(totalEntries)

	//get entries
	entries, err := database.FindEntriesByAccessionIDPaginated(accession.ID, pagination)
	if err != nil {
		ThrowError(http.StatusBadRequest, err.Error(), c, true)
		return
	}

	//get repository
	repository, err := database.FindRepository(uint(accession.Resource.RepositoryID))
	if err != nil {
		ThrowError(http.StatusBadRequest, err.Error(), c, true)
		return
	}

	//get summary
	summary, err := database.GetSummaryByAccession(accession.ID)
	if err != nil {
		ThrowError(http.StatusBadRequest, err.Error(), c, true)
		return
	}

	//get users
	users, err := getUserEmailMap([]int{accession.CreatedBy, accession.UpdatedBy})
	if err != nil {
		ThrowError(http.StatusInternalServerError, err.Error(), c, true)
		return
	}

	//return
	c.HTML(http.StatusOK, "accessions-show.html", gin.H{
		"accession":       accession,
		"repository":      repository,
		"entries":         entries,
		"isAuthenticated": true,
		"isAdmin":         sessionCookies.IsAdmin,
		"summary":         summary,
		"totals":          summary.GetTotals(),
		"users":           users,
		"isLoggedIn": true,
		"user":            user,
		"pagination":      pagination,
		"limitValues":     LimitValues,
		"overlimit":       overlimit,
		"mediatypes":      GetMediatypes(),
	})
}

func NewAccession(c *gin.Context) {
	sessionCookies := c.MustGet(ContextKeySessionCookies).(SessionCookies)
	user := c.MustGet(ContextKeyUser).(models.User)

	resourceID, err := strconv.Atoi(c.Query("resource_id"))
	if err != nil {
		ThrowError(http.StatusBadRequest, err.Error(), c, true)
		return
	}

	resource, err := database.FindResource(uint(resourceID))
	if err != nil {
		ThrowError(http.StatusInternalServerError, err.Error(), c, true)
		return
	}

	repository, err := database.FindRepository(uint(resource.RepositoryID))
	if err != nil {
		ThrowError(http.StatusInternalServerError, err.Error(), c, true)
		return
	}

	c.HTML(200, "accessions-new.html", gin.H{
		"resource":   resource,
		"repository": repository,
		"isLoggedIn": true,
		"isAdmin":    sessionCookies.IsAdmin,
		"user":       user,
	})
}

func CreateAccession(c *gin.Context) {
	//bind the form to an accession
	accession := models.Accession{}
	if err := c.Bind(&accession); err != nil {
		ThrowError(http.StatusBadRequest, err.Error(), c, true)
		return
	}

	//get the parent resource from the database
	resource, err := database.FindResource(uint(accession.ResourceID))
	if err != nil {
		ThrowError(http.StatusBadRequest, err.Error(), c, true)
		return
	}

	if err := c.Bind(&accession); err != nil {
		ThrowError(http.StatusBadRequest, err.Error(), c, true)
		return
	}
	accession.Resource = resource

	//get the current user's id
	userID, err := getUserkey(c)
	if err != nil {
		ThrowError(http.StatusBadRequest, err.Error(), c, true)
		return
	}

	//update timestamps and users
	accession.CreatedAt = time.Now()
	accession.CreatedBy = userID
	accession.UpdatedAt = time.Now()
	accession.UpdatedBy = userID

	//insert the accession Record
	accessionID, err := database.InsertAccession(&accession)
	if err != nil {
		ThrowError(http.StatusBadRequest, err.Error(), c, true)
		return
	}

	//redirect to show
	c.Redirect(http.StatusFound, fmt.Sprintf(AccessionsShow, accessionID))

}

func EditAccession(c *gin.Context) {
	sessionCookies := c.MustGet(ContextKeySessionCookies).(SessionCookies)
	user := c.MustGet(ContextKeyUser).(models.User)

	accessionID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	accession, err := database.FindAccession(uint(accessionID))
	if err != nil {
		ThrowError(http.StatusBadRequest, err.Error(), c, true)
		return
	}

	repository, err := database.FindRepository(accession.Resource.RepositoryID)
	if err != nil {
		ThrowError(http.StatusBadRequest, err.Error(), c, true)
		return
	}

	c.HTML(200, "accessions-edit.html", gin.H{
		"isAdmin":    sessionCookies.IsAdmin,
		"accession":  accession,
		"repository": repository,
		"isLoggedIn": true,
		"user":       user,
	})

}

func UpdateAccession(c *gin.Context) {
	accessionID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		ThrowError(http.StatusBadRequest, err.Error(), c, true)
		return
	}

	accession, err := database.FindAccession(uint(accessionID))
	if err != nil {
		ThrowError(http.StatusBadRequest, err.Error(), c, true)
		return
	}

	updatedAccession := models.Accession{}
	if err := c.Bind(&updatedAccession); err != nil {
		ThrowError(http.StatusBadRequest, err.Error(), c, true)
		return
	}

	userId, err := getUserkey(c)
	if err != nil {
		ThrowError(http.StatusBadRequest, err.Error(), c, true)
		return
	}

	accession.UpdatedBy = userId
	accession.UpdatedAt = time.Now()
	accession.AccessionNum = updatedAccession.AccessionNum

	if err := database.UpdateAccession(&accession); err != nil {
		ThrowError(http.StatusBadRequest, err.Error(), c, true)
		return
	}

	c.Redirect(http.StatusFound, fmt.Sprintf(AccessionsShow, accession.ID))
}

func DeleteAccession(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		ThrowError(http.StatusBadRequest, err.Error(), c, true)
		return
	}

	accession, err := database.FindAccession(uint(id))
	if err != nil {
		ThrowError(http.StatusBadRequest, err.Error(), c, true)
		return
	}

	if err := database.DeleteAccession(uint(id)); err != nil {
		ThrowError(http.StatusBadRequest, err.Error(), c, true)
		return
	}

	c.Redirect(http.StatusFound, fmt.Sprintf("/resources/%d/show", accession.ResourceID))
}

type Slew struct {
	AccessionID    uint    `form:"accession_id"`
	NumObjects     int     `form:"num_objects"`
	Mediatype      string  `form:"mediatype"`
	MediaStockSize float32 `form:"media_stock_size"`
	MediaStockUnit string  `form:"media_stock_unit"`
	BoxNum         int     `form:"box_num"`
	userID         int
}

func SlewAccession(c *gin.Context) {
	sessionCookies := c.MustGet(ContextKeySessionCookies).(SessionCookies)
	user := c.MustGet(ContextKeyUser).(models.User)

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		ThrowError(http.StatusBadRequest, err.Error(), c, true)
		return
	}

	accession, err := database.FindAccession(uint(id))
	if err != nil {
		ThrowError(http.StatusInternalServerError, err.Error(), c, true)
		return
	}

	repository, err := database.FindRepository(accession.Resource.RepositoryID)
	if err != nil {
		ThrowError(http.StatusInternalServerError, err.Error(), c, true)
		return
	}

	pagination := database.Pagination{Limit: 10, Offset: 0, Sort: "media_id"}

	entries, err := database.FindEntriesByAccessionIDPaginated(accession.ID, pagination)
	if err != nil {
		ThrowError(http.StatusInternalServerError, err.Error(), c, true)
		return
	}

	c.HTML(200, "accessions-slew.html", gin.H{
		"is_admin":    sessionCookies.IsAdmin,
		"accession":   accession,
		"repository":  repository,
		"stock_units": getStockUnits(),
		"pagination":  pagination,
		"page":        0,
		"entries":     entries,
		"isLoggedIn": true,
		"user":        user,
	})
}

func CreateAccessionSlew(c *gin.Context) {
	var slew = Slew{}

	if err := c.Bind(&slew); err != nil {
		ThrowError(http.StatusInternalServerError, err.Error(), c, true)
		return
	}

	accession, err := database.FindAccession(uint(slew.AccessionID))
	if err != nil {
		ThrowError(http.StatusInternalServerError, err.Error(), c, true)
		return
	}

	userId, err := getUserkey(c)
	if err != nil {
		ThrowError(http.StatusInternalServerError, err.Error(), c, true)
		return
	}

	slew.userID = userId

	if err := createSlewEntry(slew, accession); err != nil {
		ThrowError(http.StatusInternalServerError, err.Error(), c, true)
		return
	}

	c.Redirect(http.StatusFound, fmt.Sprintf(AccessionsShow, accession.ID))
}

func createSlewEntry(slew Slew, accession models.Accession) error {

	for i := 0; i < slew.NumObjects; i++ {
		entry := models.Entry{}
		id, _ := uuid.NewUUID()
		entry.ID = id
		mediaID, err := database.FindNextMediaCollectionInResource(accession.ResourceID)
		userID := slew.userID

		if err != nil {
			return err
		}

		resource, err := database.FindResource(uint(accession.ResourceID))
		if err != nil {
			return err
		}

		repository, err := database.FindRepository(uint(resource.RepositoryID))
		if err != nil {
			return err
		}

		entry.MediaID = mediaID
		entry.AccessionID = accession.ID
		entry.RepositoryID = accession.Resource.RepositoryID
		entry.Repository = repository
		entry.ResourceID = accession.ResourceID
		entry.Resource = resource
		entry.Mediatype = slew.Mediatype
		entry.StockSizeNum = slew.MediaStockSize
		entry.StockUnit = slew.MediaStockUnit
		entry.CreatedBy = userID
		entry.CreatedAt = time.Now()
		entry.UpdatedBy = userID
		entry.UpdatedAt = time.Now()

		if err := database.InsertEntry(&entry); err != nil {
			return err
		}
	}
	return nil
}

func AccessionGenCSV(c *gin.Context) {

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		ThrowError(http.StatusBadRequest, err.Error(), c, true)
		return
	}

	accession, err := database.FindAccession(uint(id))
	if err != nil {
		ThrowError(http.StatusInternalServerError, err.Error(), c, true)
		return
	}

	filter := c.Request.URL.Query()["filter"][0]

	entries, err := database.FindEntriesByAccessionIDFiltered(accession.ID, filter)
	if err != nil {
		ThrowError(http.StatusBadRequest, err.Error(), c, true)
		return
	}

	repository, err := database.FindRepository(accession.Resource.RepositoryID)
	if err != nil {
		ThrowError(http.StatusBadRequest, err.Error(), c, true)
		return
	}

	csvFileName := fmt.Sprintf("%s_%s_%s.csv", repository.Slug, accession.Resource.CollectionCode, accession.AccessionNum)
	csvBuffer := new(strings.Builder)
	var csvWriter = csv.NewWriter(csvBuffer)
	csvWriter.Write(models.CSVHeader)
	for _, entry := range entries {
		record := entry.ToCSV()
		csvWriter.Write(record)
	}
	csvWriter.Flush()

	c.Header("content-type", "text/csv")
	c.Header("Content-Description", "File Transfer")
	c.Header("Content-Disposition", "attachment; filename="+csvFileName)
	c.Status(http.StatusOK)
	c.Writer.Write([]byte(csvBuffer.String()))
}
