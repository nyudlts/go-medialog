package api

import (
	"fmt"
	"net/http"
	"regexp"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/nyudlts/bytemath"
	"github.com/nyudlts/go-medialog/database"
)

type SummaryAndTotal struct {
	RepositoryID int                `json:"repository_id"`
	Repository   string             `json:"repository`
	Total        int64              `json:"total_Size"`
	TotalHuman   string             `json:"total_human_size"`
	Summaries    database.Summaries `json:"summaries"`
}

func SummaryDateRange(c *gin.Context) {
	if _, err := checkToken(c); err != nil {
		c.JSON(http.StatusUnauthorized, err.Error())
		return
	}

	//get the start and end dates
	datePattern := regexp.MustCompile("[0-9]{8}")
	startDate := c.Query("start_date")
	endDate := c.Query("end_date")

	//validate start and end dates
	if !datePattern.MatchString(startDate) {
		c.JSON(http.StatusBadRequest, fmt.Errorf("invalid start_date"))
		return
	}

	if !datePattern.MatchString(endDate) {
		c.JSON(http.StatusBadRequest, fmt.Errorf("invalid end_date"))
		return
	}

	var isRefreshed bool
	//parse is refreshed
	if c.Query("is_refreshed") == "true" {
		isRefreshed = true
	} else {
		isRefreshed = false
	}

	//create a date range variable
	var dr = database.DateRange{}
	dr.IsRefreshed = isRefreshed
	//parse the repository id
	var err error
	dr.RepositoryID, err = strconv.Atoi(c.Query("repository_id"))
	if err != nil {
		dr.RepositoryID = 0
	}

	dr.StartYear, _ = strconv.Atoi(startDate[0:4])
	dr.StartMonth, _ = strconv.Atoi(startDate[4:6])
	dr.StartDay, _ = strconv.Atoi(startDate[6:8])
	dr.EndYear, _ = strconv.Atoi(endDate[0:4])
	dr.EndMonth, _ = strconv.Atoi(endDate[4:6])
	dr.EndDay, _ = strconv.Atoi(endDate[6:8])

	summaries, err := database.GetSummaryByDateRange(dr)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	var slug string
	if dr.RepositoryID == 0 {
		slug = "all"
	} else {
		repository, err := database.FindRepository(uint(dr.RepositoryID))
		if err != nil {
			c.JSON(http.StatusInternalServerError, err.Error())
			return
		}
		slug = repository.Slug
	}

	var totalSize int64 = 0
	for _, hit := range summaries.GetSlice() {
		totalSize = totalSize + int64(hit.Size)
	}

	//convert the total size to a human readable format
	f64 := float64(totalSize)
	tsize := bytemath.ConvertToBytes(f64, bytemath.B)
	humanSize := bytemath.ConvertBytesToHumanReadable(int64(tsize))

	c.JSON(http.StatusOK, SummaryAndTotal{Repository: slug, Total: totalSize, Summaries: summaries, TotalHuman: humanSize, RepositoryID: dr.RepositoryID})
}
