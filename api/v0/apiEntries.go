package api

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/nyudlts/go-medialog/controllers"
	"github.com/nyudlts/go-medialog/database"
	"github.com/nyudlts/go-medialog/models"
)

func CreateEntryV0(c *gin.Context) {
	token, err := checkToken(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, err.Error())
		return
	}

	userID, err := database.FindUserIDByToken(token)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	entry := models.Entry{}
	if err := c.Bind(&entry); err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	entry.CreatedBy = int(userID)
	entry.UpdatedBy = int(userID)
	entry.CreatedAt = time.Now()
	entry.UpdatedAt = time.Now()
	entry.ID, _ = uuid.NewUUID()

	accession, err := database.FindAccession(entry.AccessionID)
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}
	entry.Accession = accession

	resource, err := database.FindResource(accession.ResourceID)
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}
	entry.Resource = resource

	repository, err := database.FindRepository(resource.RepositoryID)
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}
	entry.Repository = repository

	if err = database.InsertEntry(&entry); err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, entry)
}

func DeleteEntryV0(c *gin.Context) {
	_, err := checkToken(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, err.Error())
		return
	}

	entryID := c.Param("id")
	entryUUID, err := uuid.Parse(entryID)
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	if err := database.DeleteEntry(entryUUID); err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, fmt.Sprintf("Entry %s deleted", entryUUID))
}

func GetEntryV0(c *gin.Context) {
	_, err := checkToken(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, err.Error())
		return
	}

	id := c.Param("id")

	uId, err := uuid.Parse(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	entry, err := database.FindEntry(uId)
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	c.JSON(http.StatusOK, entry)
}

func GetEntriesV0(c *gin.Context) {

	_, err := checkToken(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, ACCESS_DENIED)
		return
	}

	allIDsParam := c.Query("all_ids")

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
		ids, err := database.GetEntryIDs()
		if err != nil {
			c.JSON(http.StatusBadRequest, err.Error())
			return
		} else {
			c.JSON(http.StatusOK, ids)
			return
		}
	} else {

		pageSizeParam := c.Query("page_size")
		var pageSize int
		if pageSizeParam != "" {
			var err error
			pageSize, err = strconv.Atoi(pageSizeParam)
			if err != nil {
				c.JSON(http.StatusBadRequest, err.Error())
				return
			}
		} else {
			pageSize = 25
		}

		pageParam := c.Query("page")
		var entries []models.Entry
		var page int
		if pageParam != "" {
			var err error
			page, err = strconv.Atoi(pageParam)
			if err != nil {
				c.JSON(http.StatusBadRequest, err.Error())
				return
			}

			pagination := database.Pagination{Offset: page, Limit: pageSize}
			fmt.Println(pagination)
			entries, err = database.FindEntriesPaginated(pagination)
			if err != nil {
				c.JSON(http.StatusBadRequest, err.Error())
				return
			}
		} else {
			page = 1
		}

		results := EntryResultSet{}
		results.Total = database.GetCountOfEntriesInDB()
		r := int(results.Total / int64(pageSize))
		m := int(results.Total % int64(pageSize))
		var t int
		if m > 0 {
			t = r + 1
		}
		results.Results = entries
		results.FirstPage = 1
		results.ThisPage = page
		results.LastPage = t

		c.JSON(http.StatusOK, results)
		return
	}
}

func UpdateEntryLocationV0(c *gin.Context) {
	token, err := checkToken(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, err.Error())
		return
	}

	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, "no id provided")
		return
	}

	uid, err := uuid.Parse(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, "provided id is not a valid uuid")
		return
	}

	location := c.Query("location")
	if location == "" {
		c.JSON(http.StatusBadRequest, "no location provided")
		return
	}

	storageLocation := controllers.GetStorageLocation(location)

	if storageLocation == "No Match" {
		c.JSON(http.StatusBadRequest, fmt.Sprintf("`%s` is not a valid location", location))
		return
	}

	entry, err := database.FindEntry(uid)
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	userID, err := database.FindUserIDByToken(token)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	entry.Location = location
	entry.UpdatedAt = time.Now()
	entry.UpdatedBy = int(userID)

	if err := database.UpdateEntry(&entry); err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, fmt.Sprintf("id: %s, location: %s, storage location: %s", id, location, storageLocation))

}
