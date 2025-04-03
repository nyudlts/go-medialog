package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/nyudlts/go-medialog/database"
)

func GlobalSearch(c *gin.Context) {

	user, err := checkLogin(c)
	if err != nil {
		ThrowError(http.StatusBadRequest, err.Error(), c, false)
		return
	}
	query := c.Query("query")
	//get Repository Matches
	repositories, err := database.SearchRepositories(query)
	if err != nil {
		ThrowError(http.StatusInternalServerError, err.Error(), c, true)
		return
	}

	//get Resource Matches
	resources, err := database.SearchResources(query)
	if err != nil {
		ThrowError(http.StatusInternalServerError, err.Error(), c, true)
		return
	}

	//get Accession Matches
	accessions, err := database.SearchAccessions(query)
	if err != nil {
		ThrowError(http.StatusInternalServerError, err.Error(), c, true)
		return
	}

	//get Entry matches
	entries, err := database.SearchEntries(query)
	if err != nil {
		ThrowError(http.StatusInternalServerError, err.Error(), c, true)
		return
	}

	c.HTML(200, "results.html", gin.H{
		"user":         user,
		"isLoggedIn":   true,
		"isAdmin":      user.IsAdmin,
		"query":        query,
		"repositories": repositories,
		"resources":    resources,
		"accessions":   accessions,
		"entries":      entries,
	})

}
