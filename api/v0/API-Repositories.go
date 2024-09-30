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
