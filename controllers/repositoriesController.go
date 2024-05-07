package controllers

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/nyudlts/go-medialog/database"
	"github.com/nyudlts/go-medialog/models"
)

func NewRepository(c *gin.Context) {
	if !isLoggedIn(c) {
		c.Redirect(302, "/401")
		return
	}

	isAdmin := getCookie("is-admin", c)

	c.HTML(http.StatusOK, "repositories-new.html", gin.H{
		"isAdmin":         isAdmin,
		"isAuthenticated": true,
	})
}

func CreateRepository(c *gin.Context) {
	var input = models.Repository{}
	if err := c.Bind(&input); err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
	}
	input.CreatedAt = time.Now()
	input.UpdatedAt = time.Now()

	if err := database.CreateRepository(input); err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
	}
	c.Redirect(302, "/repositories")
}

func GetRepositories(c *gin.Context) {
	if !isLoggedIn(c) {
		c.Redirect(302, "/401")
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
	})

}

func GetRepository(c *gin.Context) {
	if !isLoggedIn(c) {
		c.Redirect(302, "/401")
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

	collections, err := database.FindResourcesByRepositoryID(repository.ID)
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	c.HTML(http.StatusOK, "repositories-show.html", gin.H{
		"repository":      repository,
		"resources":       collections,
		"isAdmin":         isAdmin,
		"isAuthenticated": true,
	})
}
