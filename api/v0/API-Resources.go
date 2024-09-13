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

func CreateResourceV0(c *gin.Context) {
	token, err := checkToken(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, ACCESS_DENIED)
		return
	}

	resource := models.Resource{}
	if err := c.Bind(&resource); err != nil {
		c.JSON(http.StatusBadRequest, err)
		return
	}

	userID, err := database.FindUserIDByToken(token)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	repository, err := database.FindRepository(resource.RepositoryID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	resource.CreatedBy = int(userID)
	resource.UpdatedBy = int(userID)
	resource.CreatedAt = time.Now()
	resource.UpdatedAt = time.Now()
	resource.Repository = repository

	_, err = database.InsertResource(&resource)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, resource)
}

func DeleteResourceV0(c *gin.Context) {
	_, err := checkToken(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, ACCESS_DENIED)
		return
	}

	resourceIDParam := c.Param("id")
	resourceID, err := strconv.Atoi(resourceIDParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, err)
		return
	}

	if err := database.DeleteResource(uint(resourceID)); err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, fmt.Sprintf("Resource %d deleted", resourceID))

}

func GetResourcesV0(c *gin.Context) {

	_, err := checkToken(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, ACCESS_DENIED)
		return
	}

	resources, err := database.FindResources()
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error)
		return
	}

	c.JSON(http.StatusOK, resources)
}

func GetResourceV0(c *gin.Context) {

	_, err := checkToken(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, ACCESS_DENIED)
		return
	}

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error)
	}

	resource, err := database.FindResource(uint(id))
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error)
		return
	}

	c.JSON(http.StatusOK, resource)
}

func GetResourceEntriesV0(c *gin.Context) {

	_, err := checkToken(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, ACCESS_DENIED)
		return
	}

	resourceIDParam := c.Param("id")
	resourceID, err := strconv.Atoi(resourceIDParam)
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
		entries, err := database.FindEntryIDsByResourceID(uint(resourceID))
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

		entries, err := database.FindEntriesByResourceIDPaginated(uint(resourceID), pagination)
		if err != nil {
			c.JSON(http.StatusInternalServerError, err.Error())
			return
		}

		e := EntryResultSet{}
		e.Total = database.GetCountOfEntriesInResource(uint(resourceID))
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

func GetResourceSummaryV0(c *gin.Context) {

	_, err := checkToken(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, ACCESS_DENIED)
		return
	}

	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	resource, err := database.FindResource(uint(id))
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	summaries, err := database.GetSummaryByResource(uint(id))
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	resourceSummary := SummaryTotalsResource{}
	resourceSummary.ResourceIdentifier = resource.CollectionCode
	resourceSummary.ResourceTitle = resource.Title
	resourceSummary.Totals = summaries.GetTotals()
	resourceSummary.Summaries = summaries.GetSlice()

	c.JSON(http.StatusOK, resourceSummary)
}
