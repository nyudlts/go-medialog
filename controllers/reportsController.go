package controllers

import (
	"encoding/csv"
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/nyudlts/go-medialog/database"
	"github.com/nyudlts/go-medialog/models"
)

func ReportsIndex(c *gin.Context) {

	if err := isLoggedIn(c); err != nil {
		ThrowError(http.StatusUnauthorized, err.Error(), c, false)
		return
	}

	isLoggedIn := true

	partnerCodes, err := database.GetRepositoryMap()
	if err != nil {
		ThrowError(http.StatusInternalServerError, err.Error(), c, isLoggedIn)
		return
	}

	sessionCookies, err := getSessionCookies(c)
	if err != nil {
		ThrowError(http.StatusInternalServerError, err.Error(), c, isLoggedIn)
		return
	}

	user, err := database.GetRedactedUser(sessionCookies.UserID)
	if err != nil {
		ThrowError(http.StatusBadRequest, err.Error(), c, isLoggedIn)
		return
	}

	c.HTML(http.StatusOK, "reports-index.html", gin.H{
		"months":        months,
		"days":          days,
		"years":         years,
		"partner_codes": partnerCodes,
		"isLoggedIn":    isLoggedIn,
		"isAdmin":       sessionCookies.IsAdmin,
		"user":          user,
	})
}

func ReportsRange(c *gin.Context) {

	if err := isLoggedIn(c); err != nil {
		ThrowError(http.StatusUnauthorized, err.Error(), c, false)
		return
	}

	isLoggedIn := true

	sessionCookies, err := getSessionCookies(c)
	if err != nil {
		ThrowError(http.StatusInternalServerError, err.Error(), c, isLoggedIn)
		return
	}

	user, err := database.GetRedactedUser(sessionCookies.UserID)
	if err != nil {
		ThrowError(http.StatusBadRequest, err.Error(), c, isLoggedIn)
		return
	}

	var dateRange = database.DateRange{}
	if err := c.Bind(&dateRange); err != nil {
		ThrowError(http.StatusBadRequest, err.Error(), c, isLoggedIn)
		return
	}

	summary, err := database.GetSummaryByDateRange(dateRange)
	if err != nil {
		ThrowError(http.StatusBadRequest, err.Error(), c, isLoggedIn)
		return
	}

	partnerCodes, err := database.GetRepositoryMap()
	if err != nil {
		ThrowError(http.StatusInternalServerError, err.Error(), c, isLoggedIn)
		return
	}

	c.HTML(http.StatusOK, "reports-range.html", gin.H{
		"summary":       summary,
		"totals":        summary.GetTotals(),
		"dateRange":     dateRange,
		"years":         years,
		"months":        months,
		"days":          days,
		"repository":    partnerCodes[dateRange.RepositoryID],
		"partner_codes": partnerCodes,
		"isLoggedIn":    isLoggedIn,
		"isAdmin":       sessionCookies.IsAdmin,
		"user":          user,
	})

}

func ReportsCSV(c *gin.Context) {
	var dateRange = database.DateRange{}
	if err := c.Bind(&dateRange); err != nil {
		ThrowError(http.StatusBadRequest, err.Error(), c, true)
		return
	}

	entries, err := database.GetEntriesByDateRange(dateRange)
	if err != nil {
		ThrowError(http.StatusInternalServerError, err.Error(), c, true)
		return
	}

	csvBuffer := new(strings.Builder)
	var csvWriter = csv.NewWriter(csvBuffer)
	csvWriter.Write(models.CSVHeader)
	for _, entry := range entries {
		record := entry.ToCSV()
		csvWriter.Write(record)
	}
	csvWriter.Flush()

	csvFileName := fmt.Sprintf("%s.csv", "report") // make astring formatter for dateRage struct
	c.Header("content-type", "text/csv")
	c.Header("Content-Description", "File Transfer")
	c.Header("Content-Disposition", "attachment; filename="+csvFileName)
	c.Writer.Write([]byte(csvBuffer.String()))

	c.JSON(http.StatusOK, entries)
}

var months = []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12}
var years = []int{2014, 2015, 2016, 2017, 2018, 2019, 2020, 2021, 2022, 2023, 2024} //make s range from 2014 to current year
var days = []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23, 24, 25, 26, 27, 28, 29, 30, 31}
