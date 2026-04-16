package api

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/nyudlts/go-medialog/database"
	"github.com/nyudlts/go-medialog/models"
)

// CreateAccessionV0 creates a new accession.
// @Summary      Create accession
// @Description  Creates a new accession within a resource.
// @Tags         accessions
// @Accept       json
// @Produce      json
// @Security     ApiKeyAuth
// @Param        accession  body      models.Accession  true  "Accession data"
// @Success      200        {object}  models.Accession
// @Failure      400        {string}  string
// @Failure      401        {object}  map[string]string
// @Failure      500        {string}  string
// @Router       /accessions [post]
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

// GetAccessionsV0 returns all accessions.
// @Summary      List accessions
// @Description  Returns a list of all accessions.
// @Tags         accessions
// @Produce      json
// @Security     ApiKeyAuth
// @Success      200  {array}   models.Accession
// @Failure      401  {object}  map[string]string
// @Failure      500  {string}  string
// @Router       /accessions [get]
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

// GetAccessionV0 returns an accession by ID.
// @Summary      Get accession
// @Description  Returns a single accession by its ID.
// @Tags         accessions
// @Produce      json
// @Security     ApiKeyAuth
// @Param        id   path      int  true  "Accession ID"
// @Success      200  {object}  models.Accession
// @Failure      400  {string}  string
// @Failure      401  {string}  string
// @Failure      500  {string}  string
// @Router       /accessions/{id} [get]
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

// DeleteAccessionV0 deletes an accession by ID.
// @Summary      Delete accession
// @Description  Deletes an accession by its ID.
// @Tags         accessions
// @Produce      json
// @Security     ApiKeyAuth
// @Param        id   path      int  true  "Accession ID"
// @Success      200  {string}  string
// @Failure      400  {string}  string
// @Failure      401  {string}  string
// @Failure      500  {string}  string
// @Router       /accessions/{id} [delete]
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

// GetAccessionEntriesV0 returns entries for an accession.
// @Summary      Get accession entries
// @Description  Returns paginated entries for a given accession. Use all_ids=true to return only entry UUIDs.
// @Tags         accessions
// @Produce      json
// @Security     ApiKeyAuth
// @Param        id         path   int   true   "Accession ID"
// @Param        all_ids    query  bool  false  "Return all entry IDs (no pagination)"
// @Param        page       query  int   false  "Page number"
// @Param        page_size  query  int   false  "Results per page (default 25)"
// @Success      200  {object}  EntryResultSet
// @Failure      400  {string}  string
// @Failure      401  {object}  map[string]string
// @Failure      500  {string}  string
// @Router       /accessions/{id}/entries [get]
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

		entries, err := database.FindEntriesByAccessionIDPaginated(uint(accessionID), pagination)
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

// GetAccessionSummaryV0 returns a media type summary for an accession.
// @Summary      Get accession summary
// @Description  Returns media type totals and per-type summaries for a given accession.
// @Tags         accessions
// @Produce      json
// @Security     ApiKeyAuth
// @Param        id   path      int  true  "Accession ID"
// @Success      200  {object}  SummaryTotalsAccession
// @Failure      400  {string}  string
// @Failure      401  {object}  map[string]string
// @Failure      500  {string}  string
// @Router       /accessions/{id}/summary [get]
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
