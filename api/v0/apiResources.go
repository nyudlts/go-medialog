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

// CreateResourceV0 creates a new resource.
// @Summary      Create resource
// @Description  Creates a new resource (collection) within a repository.
// @Tags         resources
// @Accept       json
// @Produce      json
// @Security     ApiKeyAuth
// @Param        resource  body      models.Resource  true  "Resource data"
// @Success      200       {object}  models.Resource
// @Failure      400       {string}  string
// @Failure      401       {object}  map[string]string
// @Failure      500       {string}  string
// @Router       /resources [post]
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

// DeleteResourceV0 deletes a resource by ID.
// @Summary      Delete resource
// @Description  Deletes a resource by its ID.
// @Tags         resources
// @Produce      json
// @Security     ApiKeyAuth
// @Param        id   path      int  true  "Resource ID"
// @Success      200  {string}  string
// @Failure      400  {string}  string
// @Failure      401  {object}  map[string]string
// @Failure      500  {string}  string
// @Router       /resources/{id} [delete]
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

// GetResourcesV0 returns all resources.
// @Summary      List resources
// @Description  Returns a list of all resources.
// @Tags         resources
// @Produce      json
// @Security     ApiKeyAuth
// @Success      200  {array}   models.Resource
// @Failure      401  {object}  map[string]string
// @Failure      500  {string}  string
// @Router       /resources [get]
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

// GetResourceV0 returns a resource by ID.
// @Summary      Get resource
// @Description  Returns a single resource by its ID.
// @Tags         resources
// @Produce      json
// @Security     ApiKeyAuth
// @Param        id   path      int  true  "Resource ID"
// @Success      200  {object}  models.Resource
// @Failure      400  {string}  string
// @Failure      401  {object}  map[string]string
// @Failure      500  {string}  string
// @Router       /resources/{id} [get]
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

// GetResourceEntriesV0 returns entries for a resource.
// @Summary      Get resource entries
// @Description  Returns paginated entries for a given resource. Use all_ids=true to return only entry UUIDs.
// @Tags         resources
// @Produce      json
// @Security     ApiKeyAuth
// @Param        id         path   int   true   "Resource ID"
// @Param        all_ids    query  bool  false  "Return all entry IDs (no pagination)"
// @Param        page       query  int   false  "Page number"
// @Param        page_size  query  int   false  "Results per page (default 25)"
// @Success      200  {object}  EntryResultSet
// @Failure      400  {string}  string
// @Failure      401  {object}  map[string]string
// @Failure      500  {string}  string
// @Router       /resources/{id}/entries [get]
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

// GetResourceSummaryV0 returns a media type summary for a resource.
// @Summary      Get resource summary
// @Description  Returns media type totals and per-type summaries for a given resource.
// @Tags         resources
// @Produce      json
// @Security     ApiKeyAuth
// @Param        id   path      int  true  "Resource ID"
// @Success      200  {object}  SummaryTotalsResource
// @Failure      400  {string}  string
// @Failure      401  {object}  map[string]string
// @Failure      500  {string}  string
// @Router       /resources/{id}/summary [get]
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
