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

	entry, err := database.FindEntry(id)
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

	entry, err := database.FindEntry(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	prevEntryID, err := database.FindEntryByMediaIDAndCollectionID(entry.MediaID-1, entry.CollectionID)
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

	entry, err := database.FindEntry(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	prevEntryID, err := database.FindEntryByMediaIDAndCollectionID(entry.MediaID+1, entry.CollectionID)
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

	if p < 0 {
		p = 0
	}

	pagination := utils.Pagination{Limit: 10, Offset: (p * 10), Sort: "updated_at desc"}

	entries, err := database.FindPaginatedEntries(pagination)
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	repositoryMap, err := database.GetRepositoryMap()
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	c.HTML(http.StatusOK, "entries-index.html", gin.H{
		"entries":         entries,
		"isAuthenticated": true,
		"isAdmin":         isAdmin,
		"page":            p,
		"repositoryMap":   repositoryMap,
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

	mediaID, err := database.FindNextMediaCollectionInResource(resource.ID)
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	c.HTML(http.StatusOK, "entries-create.html", gin.H{
		"isAdmin":                isAdmin,
		"accession":              accession,
		"resource":               resource,
		"repository":             repository,
		"mediatypes":             getMediatypes(),
		"interfaces":             getInterfaces(),
		"stock_units":            getStockUnits(),
		"optical_content_types":  getOpticalContentTypes(),
		"hdd_interfaces":         getHDDInterfaces(),
		"imaging_success":        getImageSuccess(),
		"interpretation_success": getInterpretSuccess(),
		"imaging_software":       getImagingSoftware(),
		"image_formats":          getImageFormats(),
		"media_id":               mediaID,
		"is_refreshed":           is_refreshed,
	})

}

func CreateEntry(c *gin.Context) {
	var createEntry = models.Entry{}
	if err := c.Bind(&createEntry); err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	b, err := database.IsMediaIDUniqueInResource(createEntry.MediaID, createEntry.Collection.ID)
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	if b != true {
		c.JSON(http.StatusBadRequest, fmt.Errorf("%d is not a unique ID in resource %d", createEntry.MediaID, createEntry.CollectionID))
		return
	}

	createEntry.ID, _ = uuid.NewUUID()
	createEntry.CreatedAt = time.Now()
	createEntry.UpdatedAt = time.Now()

	if err := database.InsertEntry(createEntry); err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	c.Redirect(301, fmt.Sprintf("entries/%s/show", createEntry.ID.String()))
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

	entry, err := database.FindEntry(id)
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

func EditEntry(c *gin.Context) {
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

	entry, err := database.FindEntry(id)
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

	c.HTML(http.StatusOK, "entries-edit.html", gin.H{
		"isAdmin":                isAdmin,
		"entry":                  entry,
		"accession":              entry.Accession,
		"resource":               resource,
		"repository":             repository,
		"mediatypes":             getMediatypes(),
		"interfaces":             getInterfaces(),
		"stock_units":            getStockUnits(),
		"optical_content_types":  getOpticalContentTypes(),
		"hdd_interfaces":         getHDDInterfaces(),
		"imaging_success":        getImageSuccess(),
		"interpretation_success": getInterpretSuccess(),
		"imaging_software":       getImagingSoftware(),
		"image_formats":          getImageFormats(),
		"is_refreshed":           is_refreshed,
	})
}

func UpdateEntry(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	var editedEntry = models.Entry{}

	if err := c.Bind(&editedEntry); err != nil {
		c.JSON(http.StatusBadRequest, fmt.Sprintf("%s, %s", "bind", err.Error()))
		return
	}

	entry, err := database.FindEntry(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, fmt.Sprintf("%s, %s", "find entry", err.Error()))
		return
	}

	entry.UpdateEntry(editedEntry)

	if err := database.UpdateEntry(&entry); err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	c.Redirect(301, fmt.Sprintf("/entries/%s/show", entry.ID.String()))
}
