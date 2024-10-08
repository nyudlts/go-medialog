package api

import (
	"net/http"
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
	Summary      database.Summaries `json:"summary"`
}

func SummaryDateRange(c *gin.Context) {
	if _, err := checkToken(c); err != nil {
		c.JSON(http.StatusUnauthorized, err.Error())
		return
	}

	//get the start and end dates
	startDate := c.Query("start_date")
	endDate := c.Query("end_date")

	var dr = database.DateRange{}
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

	entries, err := database.GetSummaryByDateRange(dr)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	repository, err := database.FindRepository(uint(dr.RepositoryID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	var totalSize int64 = 0
	for _, hit := range entries.GetSlice() {
		totalSize = totalSize + int64(hit.Size)
	}

	f64 := float64(totalSize)
	tsize := bytemath.ConvertToBytes(f64, bytemath.B)
	humanSize := bytemath.ConvertBytesToHumanReadable(int64(tsize))

	c.JSON(http.StatusOK, SummaryAndTotal{Repository: repository.Slug, Total: totalSize, Summary: entries, TotalHuman: humanSize, RepositoryID: dr.RepositoryID})
}
