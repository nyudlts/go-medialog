package controllers

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/nyudlts/go-medialog/database"
)

func ReportsIndex(c *gin.Context) {
	c.HTML(http.StatusOK, "reports-index.html", gin.H{
		"months": months,
		"days":   days,
		"years":  years,
	})
}

func ReportYear(c *gin.Context) {
	q := c.Request.URL.Query()
	year, err := strconv.Atoi(q.Get("selected-year"))
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	summary, err := database.GetSummaryByYear(year)
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	c.HTML(http.StatusOK, "reports-year.html", gin.H{
		"summary": summary,
		"totals":  summary.GetTotals(),
		"year":    year,
		"years":   years,
		"months":  months,
		"days":    days,
	})
}

func ReportRange(c *gin.Context) {
	var dateRange = database.DateRange{}
	if err := c.Bind(&dateRange); err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	summary, err := database.GetSummaryByDateRange(dateRange)
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	c.HTML(http.StatusOK, "reports-range.html", gin.H{
		"summary":   summary,
		"totals":    summary.GetTotals(),
		"dateRange": fmt.Sprintf("%v", dateRange),
		"years":     years,
		"months":    months,
		"days":      days,
	})

}

var months = []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12}
var years = []int{2014, 2015, 2016, 2017, 2018, 2019, 2020, 2021, 2022, 2023, 2024}
var days = []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23, 24, 25, 26, 27, 28, 29, 30, 31}
