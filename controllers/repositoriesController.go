package controllers

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/nyudlts/go-medialog/database"
	"github.com/nyudlts/go-medialog/models"
)

func NewRepository(c *gin.Context) {
	sessionCookies := c.MustGet(ContextKeySessionCookies).(SessionCookies)
	user := c.MustGet(ContextKeyUser).(models.User)

	c.HTML(http.StatusOK, "repositories-new.html", gin.H{
		"isAdmin":         sessionCookies.IsAdmin,
		"isAuthenticated": true,
		"isLoggedIn": true,
		"user":            user,
	})
}

func CreateRepository(c *gin.Context) {
	var repo = models.Repository{}
	if err := c.Bind(&repo); err != nil {
		ThrowError(http.StatusBadRequest, err.Error(), c, false)
	}

	userID, err := getUserkey(c)
	if err != nil {
		ThrowError(http.StatusBadRequest, err.Error(), c, true)
	}

	repo.CreatedAt = time.Now()
	repo.CreatedBy = userID
	repo.UpdatedAt = time.Now()
	repo.UpdatedBy = userID

	repository_id, err := database.CreateRepository(&repo)
	if err != nil {
		ThrowError(http.StatusBadRequest, err.Error(), c, true)
	}

	c.Redirect(http.StatusFound, fmt.Sprintf("/repositories/%d/show", repository_id))
}

func EditRepository(c *gin.Context) {
	sessionCookies := c.MustGet(ContextKeySessionCookies).(SessionCookies)
	user := c.MustGet(ContextKeyUser).(models.User)

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		ThrowError(http.StatusBadRequest, err.Error(), c, true)
		return
	}

	repository, err := database.FindRepository(uint(id))
	if err != nil {
		ThrowError(http.StatusBadRequest, err.Error(), c, true)
		return
	}

	c.HTML(200, "repositories-edit.html", gin.H{
		"isAdmin":    sessionCookies.IsAdmin,
		"repository": repository,
		"isLoggedIn": true,
		"user":       user,
	})

}

func UpdateRepository(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		ThrowError(http.StatusBadRequest, err.Error(), c, true)
		return
	}

	repository, err := database.FindRepository(uint(id))
	if err != nil {
		ThrowError(http.StatusBadRequest, err.Error(), c, true)
		return
	}

	var updatedRepository = models.Repository{}
	if err := c.Bind(&updatedRepository); err != nil {
		ThrowError(http.StatusBadRequest, err.Error(), c, true)
	}

	userID, err := getUserkey(c)
	if err != nil {
		ThrowError(http.StatusBadRequest, err.Error(), c, true)
		return
	}

	repository.ID = updatedRepository.ID
	repository.Title = updatedRepository.Title
	repository.Slug = updatedRepository.Slug
	repository.UpdatedAt = time.Now()
	repository.UpdatedBy = userID

	if err := database.UpdateRepository(&repository); err != nil {
		ThrowError(http.StatusBadRequest, err.Error(), c, true)
		return
	}

	c.Redirect(http.StatusFound, fmt.Sprintf("/repositories/%d/show", id))
}

func DeleteRepository(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		ThrowError(http.StatusBadRequest, err.Error(), c, true)
		return
	}

	if err := database.DeleteRepository(uint(id)); err != nil {
		ThrowError(http.StatusBadRequest, err.Error(), c, true)
		return
	}

	c.Redirect(http.StatusFound, "/repositories")
}

func GetRepositories(c *gin.Context) {
	sessionCookies := c.MustGet(ContextKeySessionCookies).(SessionCookies)
	user := c.MustGet(ContextKeyUser).(models.User)

	repositories, err := database.FindRepositories()
	if err != nil {
		ThrowError(http.StatusBadRequest, err.Error(), c, true)
		return
	}

	c.HTML(http.StatusOK, "repositories-index.html", gin.H{
		"repositories": repositories,
		"isAdmin":      sessionCookies.IsAdmin,
		"isLoggedIn": true,
		"user":         user,
	})

}

func GetRepository(c *gin.Context) {
	sessionCookies := c.MustGet(ContextKeySessionCookies).(SessionCookies)
	user := c.MustGet(ContextKeyUser).(models.User)

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	repository, err := database.FindRepository(uint(id))
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	resources, err := database.FindResourcesByRepositoryID(repository.ID)
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	c.HTML(http.StatusOK, "repositories-show.html", gin.H{
		"repository": repository,
		"resources":  resources,
		"isAdmin":    sessionCookies.IsAdmin,
		"isLoggedIn": true,
		"user":       user,
	})
}
