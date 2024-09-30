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

	if err := isLoggedIn(c); err != nil {
		ThrowError(http.StatusUnauthorized, err.Error(), c, false)
		return
	}

	isLoggedIn := true

	sessionCookies, err := getSessionCookies(c)
	if err != nil {
		ThrowError(http.StatusInternalServerError, err.Error(), c, true)
		return
	}

	user, err := database.GetRedactedUser(sessionCookies.UserID)
	if err != nil {
		ThrowError(http.StatusBadRequest, err.Error(), c, true)
		return
	}

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
		"isLoggedIn":    isLoggedIn,
		"user":          user,
	})
}

func GetAccession(c *gin.Context) {
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

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		ThrowError(http.StatusBadRequest, err.Error(), c, isLoggedIn)
		return
	}

	accession, err := database.FindAccession(uint(id))
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

	pagination := database.Pagination{Limit: 10, Offset: (p * 10), Sort: "media_id"}

	entries, err := database.FindEntriesByAccessionIDPaginated(accession.ID, pagination)
	if err != nil {
		ThrowError(http.StatusBadRequest, err.Error(), c, isLoggedIn)
		return
	}

	entryCount := database.GetCountOfEntriesInAccession(accession.ID)

	repository, err := database.FindRepository(uint(accession.Resource.RepositoryID))
	if err != nil {
		ThrowError(http.StatusBadRequest, err.Error(), c, isLoggedIn)
		return
	}

	summary, err := database.GetSummaryByAccession(accession.ID)
	if err != nil {
		ThrowError(http.StatusBadRequest, err.Error(), c, isLoggedIn)
		return
	}

	users, err := getUserEmailMap([]int{accession.CreatedBy, accession.UpdatedBy})
	if err != nil {
		ThrowError(http.StatusInternalServerError, err.Error(), c, isLoggedIn)
		return
	}

	c.HTML(http.StatusOK, "accessions-show.html", gin.H{
		"accession":       accession,
		"repository":      repository,
		"entries":         entries,
		"isAuthenticated": true,
		"isAdmin":         sessionCookies.IsAdmin,
		"page":            p,
		"summary":         summary,
		"totals":          summary.GetTotals(),
		"users":           users,
		"entryCount":      entryCount,
		"isLoggedIn":      isLoggedIn,
		"user":            user,
	})
}

func NewAccession(c *gin.Context) {
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

	resourceID, err := strconv.Atoi(c.Query("resource_id"))
	if err != nil {
		ThrowError(http.StatusBadRequest, err.Error(), c, isLoggedIn)
		return
	}

	resource, err := database.FindResource(uint(resourceID))
	if err != nil {
		ThrowError(http.StatusInternalServerError, err.Error(), c, isLoggedIn)
		return
	}

	repository, err := database.FindRepository(uint(resource.RepositoryID))
	if err != nil {
		ThrowError(http.StatusInternalServerError, err.Error(), c, isLoggedIn)
		return
	}

	c.HTML(200, "accessions-new.html", gin.H{
		"resource":   resource,
		"repository": repository,
		"isLoggedIn": isLoggedIn,
		"isAdmin":    sessionCookies.IsAdmin,
		"user":       user,
	})
}

func CreateAccession(c *gin.Context) {
	//check the user is logged in
	if err := isLoggedIn(c); err != nil {
		ThrowError(http.StatusUnauthorized, err.Error(), c, false)
		return
	}

	isLoggedIn := true

	//bind the form to an accession
	accession := models.Accession{}
	if err := c.Bind(&accession); err != nil {
		ThrowError(http.StatusBadRequest, err.Error(), c, isLoggedIn)
		return
	}

	//get the parent resource from the database
	resource, err := database.FindResource(uint(accession.ResourceID))
	if err != nil {
		ThrowError(http.StatusBadRequest, err.Error(), c, isLoggedIn)
		return
	}

	if err := c.Bind(&accession); err != nil {
		ThrowError(http.StatusBadRequest, err.Error(), c, isLoggedIn)
		return
	}
	accession.Resource = resource

	//get the current user's id
	userID, err := getUserkey(c)
	if err != nil {
		ThrowError(http.StatusBadRequest, err.Error(), c, isLoggedIn)
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
		ThrowError(http.StatusBadRequest, err.Error(), c, isLoggedIn)
		return
	}

	//redirect to show
	c.Redirect(http.StatusFound, fmt.Sprintf(AccessionsShow, accessionID))

}

func EditAccession(c *gin.Context) {
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

	accessionID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	accession, err := database.FindAccession(uint(accessionID))
	if err != nil {
		ThrowError(http.StatusBadRequest, err.Error(), c, isLoggedIn)
		return
	}

	repository, err := database.FindRepository(accession.Resource.RepositoryID)
	if err != nil {
		ThrowError(http.StatusBadRequest, err.Error(), c, isLoggedIn)
		return
	}

	c.HTML(200, "accessions-edit.html", gin.H{
		"isAdmin":    sessionCookies.IsAdmin,
		"accession":  accession,
		"repository": repository,
		"isLoggedIn": isLoggedIn,
		"user":       user,
	})

}

func UpdateAccession(c *gin.Context) {
	if err := isLoggedIn(c); err != nil {
		ThrowError(http.StatusUnauthorized, err.Error(), c, false)
		return
	}

	isLoggedIn := true

	accessionID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		ThrowError(http.StatusBadRequest, err.Error(), c, isLoggedIn)
		return
	}

	accession, err := database.FindAccession(uint(accessionID))
	if err != nil {
		ThrowError(http.StatusBadRequest, err.Error(), c, isLoggedIn)
		return
	}

	updatedAccession := models.Accession{}
	if err := c.Bind(&updatedAccession); err != nil {
		ThrowError(http.StatusBadRequest, err.Error(), c, isLoggedIn)
		return
	}

	userId, err := getUserkey(c)
	if err != nil {
		ThrowError(http.StatusBadRequest, err.Error(), c, isLoggedIn)
		return
	}

	accession.UpdatedBy = userId
	accession.UpdatedAt = time.Now()
	accession.AccessionNum = updatedAccession.AccessionNum

	if err := database.UpdateAccession(&accession); err != nil {
		ThrowError(http.StatusBadRequest, err.Error(), c, isLoggedIn)
		return
	}

	c.Redirect(http.StatusFound, fmt.Sprintf(AccessionsShow, accession.ID))
}

func DeleteAccession(c *gin.Context) {
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

	accession, err := database.FindAccession(uint(id))
	if err != nil {
		ThrowError(http.StatusBadRequest, err.Error(), c, isLoggedIn)
		return
	}

	if err := database.DeleteAccession(uint(id)); err != nil {
		ThrowError(http.StatusBadRequest, err.Error(), c, isLoggedIn)
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

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		ThrowError(http.StatusBadRequest, err.Error(), c, isLoggedIn)
		return
	}

	accession, err := database.FindAccession(uint(id))
	if err != nil {
		ThrowError(http.StatusInternalServerError, err.Error(), c, isLoggedIn)
		return
	}

	repository, err := database.FindRepository(accession.Resource.RepositoryID)
	if err != nil {
		ThrowError(http.StatusInternalServerError, err.Error(), c, isLoggedIn)
		return
	}

	pagination := database.Pagination{Limit: 10, Offset: 0, Sort: "media_id"}

	entries, err := database.FindEntriesByAccessionIDPaginated(accession.ID, pagination)
	if err != nil {
		ThrowError(http.StatusInternalServerError, err.Error(), c, isLoggedIn)
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
		"isLoggedIn":  isLoggedIn,
		"user":        user,
	})
}

func CreateAccessionSlew(c *gin.Context) {
	if err := isLoggedIn(c); err != nil {
		ThrowError(http.StatusUnauthorized, err.Error(), c, false)
		return
	}
	isLoggedIn := true

	var slew = Slew{}

	if err := c.Bind(&slew); err != nil {
		ThrowError(http.StatusInternalServerError, err.Error(), c, isLoggedIn)
		return
	}

	accession, err := database.FindAccession(uint(slew.AccessionID))
	if err != nil {
		ThrowError(http.StatusInternalServerError, err.Error(), c, isLoggedIn)
		return
	}

	userId, err := getUserkey(c)
	if err != nil {
		ThrowError(http.StatusInternalServerError, err.Error(), c, isLoggedIn)
		return
	}

	slew.userID = userId

	if err := createSlewEntry(slew, accession); err != nil {
		ThrowError(http.StatusInternalServerError, err.Error(), c, isLoggedIn)
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

	accession, err := database.FindAccession(uint(id))
	if err != nil {
		ThrowError(http.StatusInternalServerError, err.Error(), c, isLoggedIn)
		return
	}

	entries, err := database.FindEntriesByAccessionID(accession.ID)
	if err != nil {
		ThrowError(http.StatusBadRequest, err.Error(), c, isLoggedIn)
		return
	}

	repository, err := database.FindRepository(accession.Resource.RepositoryID)
	if err != nil {
		ThrowError(http.StatusBadRequest, err.Error(), c, isLoggedIn)
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
	c.Writer.Write([]byte(csvBuffer.String()))
}
