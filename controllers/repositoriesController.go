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
	loggedIn := isLoggedIn(c)
	if !loggedIn {
		throwError(http.StatusUnauthorized, UNAUTHORIZED, c)
		return
	}

	isAdmin := getCookie("is-admin", c)

	c.HTML(http.StatusOK, "repositories-new.html", gin.H{
		"isAdmin":         isAdmin,
		"isAuthenticated": true,
		"isLoggedIn":      loggedIn,
	})
}

func CreateRepository(c *gin.Context) {
	var repo = models.Repository{}
	if err := c.Bind(&repo); err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
	}

	userID, err := getUserkey(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
	}

	repo.CreatedAt = time.Now()
	repo.CreatedBy = userID
	repo.UpdatedAt = time.Now()
	repo.UpdatedBy = userID

	repository_id, err := database.CreateRepository(&repo)
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
	}

	c.Redirect(302, fmt.Sprintf("/repositories/%d/show", repository_id))
}

func EditRepository(c *gin.Context) {
	loggedIn := isLoggedIn(c)
	if !loggedIn {
		throwError(http.StatusUnauthorized, UNAUTHORIZED, c)
		return
	}

	isAdmin := getCookie("is-admin", c)

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

	c.HTML(200, "repositories-edit.html", gin.H{
		"isAdmin":    isAdmin,
		"repository": repository,
		"isLoggedIn": loggedIn,
	})

}

func UpdateRepository(c *gin.Context) {
	loggedIn := isLoggedIn(c)
	if !loggedIn {
		throwError(http.StatusUnauthorized, UNAUTHORIZED, c)
		return
	}

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

	var updatedRepository = models.Repository{}
	if err := c.Bind(&updatedRepository); err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
	}

	userID, err := getUserkey(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	repository.ID = updatedRepository.ID
	repository.Title = updatedRepository.Title
	repository.Slug = updatedRepository.Slug
	repository.UpdatedAt = time.Now()
	repository.UpdatedBy = userID

	if err := database.UpdateRepository(&repository); err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	c.Redirect(302, fmt.Sprintf("/repositories/%d/show", id))
}

func DeleteRepository(c *gin.Context) {
	loggedIn := isLoggedIn(c)
	if !loggedIn {
		throwError(http.StatusUnauthorized, UNAUTHORIZED, c)
		return
	}

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	if err := database.DeleteRepository(uint(id)); err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	c.Redirect(302, "/repositories")
}

func GetRepositories(c *gin.Context) {
	loggedIn := isLoggedIn(c)
	if !loggedIn {
		throwError(http.StatusUnauthorized, UNAUTHORIZED, c)
		return
	}

	isAdmin := getCookie("is-admin", c)

	repositories, err := database.FindRepositories()
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	c.HTML(http.StatusOK, "repositories-index.html", gin.H{
		"repositories":    repositories,
		"isAuthenticated": true,
		"isAdmin":         isAdmin,
		"isLoggedIn":      loggedIn,
	})

}

func GetRepository(c *gin.Context) {
	loggedIn := isLoggedIn(c)
	if !loggedIn {
		throwError(http.StatusUnauthorized, UNAUTHORIZED, c)
		return
	}

	isAdmin := getCookie("is-admin", c)

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
		"repository":      repository,
		"resources":       resources,
		"isAdmin":         isAdmin,
		"isAuthenticated": loggedIn,
	})
}
