package api

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/nyudlts/go-medialog/database"
	"github.com/nyudlts/go-medialog/models"
)

// GetRepositoriesV0 returns all repositories.
// @Summary      List repositories
// @Description  Returns a list of all repositories.
// @Tags         repositories
// @Produce      json
// @Security     ApiKeyAuth
// @Success      200  {array}   models.Repository
// @Failure      401  {object}  map[string]string
// @Failure      500  {string}  string
// @Router       /repositories [get]
func GetRepositoriesV0(c *gin.Context) {

	_, err := checkToken(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, ACCESS_DENIED)
		return
	}

	repositories, err := database.FindRepositories()
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error)
		return
	}

	c.JSON(http.StatusOK, repositories)
}

// GetRepositoryV0 returns a repository by ID.
// @Summary      Get repository
// @Description  Returns a single repository by its ID.
// @Tags         repositories
// @Produce      json
// @Security     ApiKeyAuth
// @Param        id   path      int  true  "Repository ID"
// @Success      200  {object}  models.Repository
// @Failure      400  {string}  string
// @Failure      401  {object}  map[string]string
// @Failure      500  {string}  string
// @Router       /repositories/{id} [get]
func GetRepositoryV0(c *gin.Context) {

	_, err := checkToken(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, ACCESS_DENIED)
		return
	}

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error)
	}

	repository, err := database.FindRepository(uint(id))
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error)
		return
	}

	c.JSON(http.StatusOK, repository)
}

// CreateRepositoryV0 creates a new repository.
// @Summary      Create repository
// @Description  Creates a new repository record.
// @Tags         repositories
// @Accept       json
// @Produce      json
// @Security     ApiKeyAuth
// @Param        repository  body      models.Repository  true  "Repository data"
// @Success      200         {object}  models.Repository
// @Failure      400         {string}  string
// @Failure      401         {object}  map[string]string
// @Failure      500         {string}  string
// @Router       /repositories [post]
func CreateRepositoryV0(c *gin.Context) {
	token, err := checkToken(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, ACCESS_DENIED)
		return
	}

	repo := models.Repository{}
	if err := c.Bind(&repo); err != nil {
		c.JSON(http.StatusBadRequest, err)
		return
	}

	userID, err := database.FindUserIDByToken(token)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	repo.CreatedBy = int(userID)
	repo.UpdatedBy = int(userID)
	repo.CreatedAt = time.Now()
	repo.UpdatedAt = time.Now()

	_, err = database.CreateRepository(&repo)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, repo)
}

// DeleteRepositoryV0 deletes a repository by ID.
// @Summary      Delete repository
// @Description  Deletes a repository by its ID.
// @Tags         repositories
// @Produce      json
// @Security     ApiKeyAuth
// @Param        id   path      int  true  "Repository ID"
// @Success      200  {string}  string
// @Failure      400  {string}  string
// @Failure      401  {object}  map[string]string
// @Failure      500  {string}  string
// @Router       /repositories/{id} [delete]
func DeleteRepositoryV0(c *gin.Context) {
	_, err := checkToken(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, ACCESS_DENIED)
		return
	}

	repositoryIDParam := c.Param("id")
	repositoryID, err := strconv.Atoi(repositoryIDParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	if err := database.DeleteRepository(uint(repositoryID)); err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, fmt.Sprintf("Repository %d deleted", repositoryID))

}

// GetRepositoryEntriesV0 returns entries for a repository.
// @Summary      Get repository entries
// @Description  Returns paginated entries for a given repository. Use all_ids=true to return only entry UUIDs.
// @Tags         repositories
// @Produce      json
// @Security     ApiKeyAuth
// @Param        id         path   int   true   "Repository ID"
// @Param        all_ids    query  bool  false  "Return all entry IDs (no pagination)"
// @Param        page       query  int   false  "Page number"
// @Param        page_size  query  int   false  "Results per page (default 25)"
// @Success      200  {object}  EntryResultSet
// @Failure      400  {string}  string
// @Failure      401  {object}  map[string]string
// @Failure      500  {string}  string
// @Router       /repositories/{id}/entries [get]
func GetRepositoryEntriesV0(c *gin.Context) {
	_, err := checkToken(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, ACCESS_DENIED)
		return
	}

	repositoryIDParam := c.Param("id")
	repositoryID, err := strconv.Atoi(repositoryIDParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
	}

	allIDsParam := c.Query("all_ids")
	pageParam := c.Query("page")
	pageSizeParam := c.Query("page_size")
	log.Println(allIDsParam)
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
		entries, err := database.FindEntryIDsByRepositoryID(uint(repositoryID))
		if err != nil {
			c.JSON(http.StatusInternalServerError, err.Error())
			return
		}
		c.JSON(http.StatusOK, entries)
		return
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

		entries, err := database.FindEntriesByRepositoryIDPaginated(uint(repositoryID), pagination)
		if err != nil {
			c.JSON(http.StatusInternalServerError, err.Error())
			return
		}

		e := EntryResultSet{}
		e.Total = database.GetCountOfEntriesInRepository(uint(repositoryID))
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

// GetRepositorySummaryV0 returns a media type summary for a repository.
// @Summary      Get repository summary
// @Description  Returns media type totals and per-type summaries for a given repository.
// @Tags         repositories
// @Produce      json
// @Security     ApiKeyAuth
// @Param        id   path      int  true  "Repository ID"
// @Success      200  {object}  SummaryTotalsRepo
// @Failure      400  {string}  string
// @Failure      401  {object}  map[string]string
// @Failure      500  {string}  string
// @Router       /repositories/{id}/summary [get]
func GetRepositorySummaryV0(c *gin.Context) {
	_, err := checkToken(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, ACCESS_DENIED)
		return
	}

	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
	}

	repository, err := database.FindRepository(uint(id))
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	summaryMap, err := database.GetSummaryByRepository(uint(id))
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	summaryTotals := SummaryTotalsRepo{
		Repository: repository.Title,
		Totals:     summaryMap.GetTotals(),
		Summaries:  summaryMap.GetSlice(),
	}

	c.JSON(http.StatusOK, summaryTotals)
}
