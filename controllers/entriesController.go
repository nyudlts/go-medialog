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
	"github.com/nyudlts/go-medialog/utils"
)

func GetEntry(c *gin.Context) {
	if err := checkSession(c); err != nil {
		c.Redirect(302, "/")
		return
	}
	isAdmin := getCookie("is-admin", c)

	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	entry, err := database.FindEntry(id.String())
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	accession, err := database.FindAccession(entry.AccessionID)
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	resource, err := database.FindResource(uint(accession.CollectionID))
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	repository, err := database.FindRepository(uint(resource.RepositoryID))
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	c.HTML(http.StatusOK, "entries-show.html", gin.H{
		"entry":           entry,
		"accession":       accession,
		"resource":        resource,
		"repository":      repository,
		"isAuthenticated": true,
		"isAdmin":         isAdmin,
	})
}

func GetPreviousEntry(c *gin.Context) {
	if err := checkSession(c); err != nil {
		c.Redirect(302, "/")
		return
	}

	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	entry, err := database.FindEntry(id.String())
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	mediaID, err := strconv.Atoi(entry.MediaID)
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	prevEntryID, err := database.FindEntryByMediaIDAndCollectionID(mediaID-1, entry.CollectionID)
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	c.Redirect(http.StatusPermanentRedirect, fmt.Sprintf("/entries/%s/show", prevEntryID))
}

func GetNextEntry(c *gin.Context) {
	if err := checkSession(c); err != nil {
		c.Redirect(302, "/")
		return
	}

	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	entry, err := database.FindEntry(id.String())
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	mediaID, err := strconv.Atoi(entry.MediaID)
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	prevEntryID, err := database.FindEntryByMediaIDAndCollectionID(mediaID+1, entry.CollectionID)
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	c.Redirect(http.StatusPermanentRedirect, fmt.Sprintf("/entries/%s/show", prevEntryID))
}

func GetEntries(c *gin.Context) {
	if err := checkSession(c); err != nil {
		c.Redirect(302, "/")
		return
	}
	isAdmin := getCookie("is-admin", c)

	//pagination
	var p = 0
	var err error
	page := c.Request.URL.Query()["page"]

	if len(page) > 0 {
		p, err = strconv.Atoi(page[0])
		if err != nil {
			c.JSON(http.StatusBadRequest, err.Error())
			return
		}

	}

	pagination := utils.Pagination{Limit: 10, Offset: (p * 10), Sort: "updated_at desc"}

	entries, err := database.FindPaginatedEntries(pagination)
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	c.HTML(http.StatusOK, "entries-index.html", gin.H{
		"entries":         entries,
		"isAuthenticated": true,
		"isAdmin":         isAdmin,
		"page":            p,
	})
}

func NewEntry(c *gin.Context) {
	if err := checkSession(c); err != nil {
		c.Redirect(302, "/")
		return
	}
	isAdmin := getCookie("is-admin", c)

	aID := c.Request.URL.Query().Get("accession_id")
	accessionID, err := strconv.Atoi(aID)
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	accession, err := database.FindAccession(accessionID)
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	resource, err := database.FindResource(accession.Collection.ID)
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	repository, err := database.FindRepository(resource.Repository.ID)
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	c.HTML(http.StatusOK, "entries-create.html", gin.H{
		"isAdmin":    isAdmin,
		"accession":  accession,
		"resource":   resource,
		"repository": repository,
	})

}

func CreateEntry(c *gin.Context) {
	var createEntry = models.Entry{}
	if err := c.Bind(&createEntry); err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	createEntry.ID, _ = uuid.NewUUID()
	createEntry.CreatedAt = time.Now()
	createEntry.UpdatedAt = time.Now()

	if err := database.InsertEntry(createEntry); err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	c.Redirect(http.StatusPermanentRedirect, fmt.Sprintf("entries/%s/show", createEntry.ID.String()))
}

func DeleteEntry(c *gin.Context) {
	if err := checkSession(c); err != nil {
		c.Redirect(302, "/")
		return
	}

	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	entry, err := database.FindEntry(id.String())
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	if err := database.DeleteEntry(id); err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	c.Redirect(http.StatusPermanentRedirect, fmt.Sprintf("/accessions/%d/show", entry.AccessionID))

}
