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
	"github.com/nyudlts/go-medialog/utils"
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

	accession, err := database.FindAccession(id)
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

	pagination := utils.Pagination{Limit: 10, Offset: (p * 10), Sort: "media_id"}

	entries, err := database.FindEntriesByAccessionID(accession.ID, pagination)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	repository, err := database.FindRepository(uint(accession.Collection.RepositoryID))
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

	resourceID, err := strconv.Atoi(c.Param("id"))
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
	if !isLoggedIn(c) {
		c.Redirect(302, "/error")
		return
	}

	accession := models.Accession{}
	if err := c.Bind(&accession); err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	log.Println("Controller", accession.CollectionID)

	userId, err := getUserkey(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	accession.CreatedAt = time.Now()
	accession.CreatedBy = userId
	accession.UpdatedAt = time.Now()
	accession.UpdatedBy = userId

	if err := database.InsertAccession(&accession); err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	c.Redirect(302, fmt.Sprintf("/accessions/%d/show", accession.ID))

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

	accession, err := database.FindAccession(accessionID)
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	repository, err := database.FindRepository(uint(accession.Collection.RepositoryID))
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

	accession, err := database.FindAccession(accessionID)
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

	accession, err := database.FindAccession(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	if err := database.DeleteAccession(id); err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	c.Redirect(302, fmt.Sprintf("/resources/accession.ResourceID/show", accession.CollectionID))
}

type Slew struct {
	AccessionID    uint    `form:"accession_id"`
	NumObjects     int     `form:"num_objects"`
	Mediatype      string  `form:"mediatype"`
	MediaStockSize float32 `form:"media_stock_size"`
	MediaStockUnit string  `form:"media_stock_unit"`
	BoxNum         int     `form:"box_num"`
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

	accession, err := database.FindAccession(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	repository, err := database.FindRepository(uint(accession.Collection.RepositoryID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	pagination := utils.Pagination{Limit: 10, Offset: 0, Sort: "media_id"}

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

	accession, err := database.FindAccession(int(slew.AccessionID))
	if err != nil {
		c.JSON(http.StatusBadRequest, fmt.Sprintf("%s, %s", "bind", err.Error()))
		return
	}

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
		mediaID, err := database.FindNextMediaCollectionInResource(uint(accession.CollectionID))
		if err != nil {
			return err
		}
		entry.MediaID = mediaID
		entry.AccessionID = int(accession.ID)
		entry.RepositoryID = int(accession.Collection.RepositoryID)
		entry.CollectionID = int(accession.Collection.ID)
		entry.Mediatype = slew.Mediatype
		entry.StockSizeNum = slew.MediaStockSize
		entry.StockUnit = slew.MediaStockUnit
		entry.CreatedAt = time.Now()
		entry.UpdatedAt = time.Now()

		if err := database.InsertEntry(entry); err != nil {
			return err
		}
	}
	return nil
}
