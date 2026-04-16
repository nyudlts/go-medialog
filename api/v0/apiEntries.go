package api

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/nyudlts/go-medialog/controllers"
	"github.com/nyudlts/go-medialog/database"
	"github.com/nyudlts/go-medialog/models"
)

// CreateEntryV0 creates a new media entry.
// @Summary      Create entry
// @Description  Creates a new media entry within an accession.
// @Tags         entries
// @Accept       json
// @Produce      json
// @Security     ApiKeyAuth
// @Param        entry  body      models.Entry  true  "Entry data"
// @Success      200    {object}  models.Entry
// @Failure      400    {string}  string
// @Failure      401    {string}  string
// @Failure      500    {string}  string
// @Router       /entries [post]
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

// DeleteEntryV0 deletes an entry by UUID.
// @Summary      Delete entry
// @Description  Deletes a media entry by its UUID.
// @Tags         entries
// @Produce      json
// @Security     ApiKeyAuth
// @Param        id   path      string  true  "Entry UUID"
// @Success      200  {string}  string
// @Failure      400  {string}  string
// @Failure      401  {string}  string
// @Failure      500  {string}  string
// @Router       /entries/{id} [delete]
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

// GetEntryV0 returns an entry by UUID.
// @Summary      Get entry
// @Description  Returns a single media entry by its UUID.
// @Tags         entries
// @Produce      json
// @Security     ApiKeyAuth
// @Param        id   path      string  true  "Entry UUID"
// @Success      200  {object}  models.Entry
// @Failure      400  {string}  string
// @Failure      401  {string}  string
// @Router       /entries/{id} [get]
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

// GetEntriesV0 returns all entries.
// @Summary      List entries
// @Description  Returns paginated entries across all accessions. Use all_ids=true to return only UUIDs.
// @Tags         entries
// @Produce      json
// @Security     ApiKeyAuth
// @Param        all_ids    query  bool  false  "Return all entry IDs (no pagination)"
// @Param        page       query  int   false  "Page number"
// @Param        page_size  query  int   false  "Results per page (default 25)"
// @Success      200  {object}  EntryResultSet
// @Failure      400  {string}  string
// @Failure      401  {object}  map[string]string
// @Router       /entries [get]
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

// UpdateEntryLocationV0 updates the storage location of an entry.
// @Summary      Update entry location
// @Description  Updates the storage location for a given entry by UUID.
// @Tags         entries
// @Produce      json
// @Security     ApiKeyAuth
// @Param        id        path   string  true  "Entry UUID"
// @Param        location  query  string  true  "Storage location code"
// @Success      200  {string}  string
// @Failure      400  {string}  string
// @Failure      401  {string}  string
// @Failure      500  {string}  string
// @Router       /entries/{id}/update_location [patch]
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

// UpdateEntryV0 updates all fields of an entry.
// @Summary      Update entry
// @Description  Replaces all updatable fields of an entry by UUID.
// @Tags         entries
// @Accept       json
// @Produce      json
// @Security     ApiKeyAuth
// @Param        id     path      string        true  "Entry UUID"
// @Param        entry  body      models.Entry  true  "Updated entry data"
// @Success      200  {string}  string
// @Failure      400  {string}  string
// @Failure      401  {string}  string
// @Failure      500  {string}  string
// @Router       /entries/{id}/update [post]
func UpdateEntryV0(c *gin.Context) {
	id := c.Param("id")

	tkn, err := checkToken(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, err.Error())
		return
	}

	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	entry := models.Entry{}
	if err := json.Unmarshal(body, &entry); err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
	}

	userID, err := database.FindUserIDByToken(tkn)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	entry.UpdatedBy = int(userID)
	entry.UpdatedAt = time.Now()

	if err := database.UpdateEntry(&entry); err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, fmt.Sprintf("entry %s updated", id))
}
