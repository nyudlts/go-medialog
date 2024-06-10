package controllers

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/nyudlts/go-medialog/database"
	"github.com/nyudlts/go-medialog/models"
)

func GetAccessions(c *gin.Context) {

	if !isLoggedIn(c) {
		c.Redirect(302, "/error")
		return
	}

	isAdmin := getCookie("is-admin", c)

	accessions, err := database.FindAccessions()
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	repositoryMap, err := database.GetRepositoryMap()
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	repositoryMap2 := map[uint]string{}
	for k, v := range repositoryMap {
		repositoryMap2[uint(k)] = v
	}

	c.HTML(200, "accessions-index.html", gin.H{
		"accessions":      accessions,
		"isAuthenticated": true,
		"isAdmin":         isAdmin,
		"repositoryMap":   repositoryMap2,
	})
}

func GetAccession(c *gin.Context) {
	if !isLoggedIn(c) {
		c.Redirect(302, "/error")
		return
	}

	isAdmin := getCookie("is-admin", c)

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	accession, err := database.FindAccession(uint(id))
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
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

	pagination := database.Pagination{Limit: 10, Offset: (p * 10), Sort: "media_id"}

	entries, err := database.FindEntriesByAccessionID(accession.ID, pagination)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	repository, err := database.FindRepository(uint(accession.Resource.RepositoryID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	summary, err := database.GetSummaryByAccession(accession.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	users, err := getUserEmailMap([]int{accession.CreatedBy, accession.UpdatedBy})
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	c.HTML(http.StatusOK, "accessions-show.html", gin.H{
		"accession":       accession,
		"repository":      repository,
		"entries":         entries,
		"isAuthenticated": true,
		"isAdmin":         isAdmin,
		"page":            p,
		"summary":         summary,
		"totals":          summary.GetTotals(),
		"users":           users,
	})
}

func NewAccession(c *gin.Context) {
	if !isLoggedIn(c) {
		c.Redirect(302, "/error")
		return
	}

	resourceID, err := strconv.Atoi(c.Query("resource_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	resource, err := database.FindResource(uint(resourceID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	repository, err := database.FindRepository(uint(resource.RepositoryID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	c.HTML(200, "accessions-new.html", gin.H{
		"resource":   resource,
		"repository": repository,
	})
}

func CreateAccession(c *gin.Context) {
	//check the user is logged in
	if !isLoggedIn(c) {
		c.Redirect(302, "/error")
		return
	}

	//bind the form to an accession
	accession := models.Accession{}
	if err := c.Bind(&accession); err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	//get the parent resource from the database
	resource, err := database.FindResource(uint(accession.ResourceID))
	if err := c.Bind(&accession); err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}
	accession.Resource = resource

	//get the current user's id
	userID, err := getUserkey(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
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
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	//redirect to show
	c.Redirect(302, fmt.Sprintf("/accessions/%d/show", accessionID))

}

func EditAccession(c *gin.Context) {
	if !isLoggedIn(c) {
		c.Redirect(302, "/error")
		return
	}

	isAdmin := getCookie("is-admin", c)

	accessionID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	accession, err := database.FindAccession(uint(accessionID))
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	repository, err := database.FindRepository(accession.Resource.RepositoryID)
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	c.HTML(200, "accessions-edit.html", gin.H{
		"isAdmin":    isAdmin,
		"accession":  accession,
		"repository": repository,
	})

}

func UpdateAccession(c *gin.Context) {
	if !isLoggedIn(c) {
		c.Redirect(302, "/error")
		return
	}

	accessionID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	accession, err := database.FindAccession(uint(accessionID))
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	updatedAccession := models.Accession{}
	if err := c.Bind(&updatedAccession); err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	userId, err := getUserkey(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	accession.UpdatedBy = userId
	accession.UpdatedAt = time.Now()
	accession.AccessionNum = updatedAccession.AccessionNum

	if err := database.UpdateAccession(&accession); err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	c.Redirect(302, fmt.Sprintf("/accessions/%d/show", accession.ID))
}

func DeleteAccession(c *gin.Context) {
	if !isLoggedIn(c) {
		c.Redirect(302, "/error")
		return
	}

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	accession, err := database.FindAccession(uint(id))
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	if err := database.DeleteAccession(uint(id)); err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	c.Redirect(302, fmt.Sprintf("/resources/%d/show", accession.ResourceID))
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
	if !isLoggedIn(c) {
		c.Redirect(302, "/error")
		return
	}

	isAdmin := getCookie("is-admin", c)

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	accession, err := database.FindAccession(uint(id))
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	repository, err := database.FindRepository(uint(accession.Resource.RepositoryID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	pagination := database.Pagination{Limit: 10, Offset: 0, Sort: "media_id"}

	entries, err := database.FindEntriesByAccessionID(accession.ID, pagination)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	c.HTML(200, "accessions-slew.html", gin.H{
		"is_admin":    isAdmin,
		"accession":   accession,
		"repository":  repository,
		"mediatypes":  getMediatypes(),
		"stock_units": getStockUnits(),
		"pagination":  pagination,
		"page":        0,
		"entries":     entries,
	})
}

func CreateAccessionSlew(c *gin.Context) {
	if !isLoggedIn(c) {
		c.Redirect(302, "/error")
		return
	}

	var slew = Slew{}

	if err := c.Bind(&slew); err != nil {
		c.JSON(http.StatusBadRequest, fmt.Sprintf("%s, %s", "bind", err.Error()))
		return
	}

	accession, err := database.FindAccession(uint(slew.AccessionID))
	if err != nil {
		c.JSON(http.StatusBadRequest, fmt.Sprintf("%s, %s", "bind", err.Error()))
		return
	}

	userId, err := getUserkey(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	slew.userID = userId

	if err := createSlewEntry(slew, accession); err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	c.Redirect(302, fmt.Sprintf("/accessions/%d/show", accession.ID))
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

		if _, err := database.InsertEntry(&entry); err != nil {
			return err
		}
	}
	return nil
}
