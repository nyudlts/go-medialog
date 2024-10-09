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
