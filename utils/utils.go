package utils

import (
	"text/template"
	"time"

	"github.com/gin-gonic/gin"
)

func Add(a int, b int) int { return a + b }

func Subtract(a int, b int) int { return a - b }

func FormatAsDate(t time.Time) string { return t.Format("2006-01-02") }

//func getMediatypes() map[string]string { return mediatypes }

func GetMediatype(s string) string {

	for k, y := range mediatypes {
		if k == s {
			return y
		}
	}
	return "No Match"
}

var mediatypes = map[string]string{
	"":                          "",
	"mediatype_transfer":        "Network Transfer",
	"mediatype_floppy_3_5":      "3.5 in. Floppy Disk",
	"mediatype_floppy_5_25":     "5.25 in. Floppy Disk",
	"mediatype_floppy_8":        "8 in. Floppy Disk",
	"mediatype_hard_disk_drive": "Hard Disk Drive",
	"mediatype_flash_drive":     "Flash Drive",
	"mediatype_dvd":             "DVD commercial",
	"mediatype_dvdr":            "DVD-R",
	"mediatype_dvdrw":           "DVD-RW",
	"mediatype_cd":              "CD commercial",
	"mediatype_cdr":             "CD-R",
	"mediatype_cdrw":            "CD-RW",
	"mediatype_cdda":            "Audio CD",
	"mediatype_jaz":             "Jaz Drive",
	"mediatype_zip":             "Zip Disk",
	"mediatype_sd":              "SD Card",
	"mediatype_minidisc":        "MiniDisc",
	"mediatype_data_cartridge":  "Data Cartridge",
	"mediatype_laserdisc":       "Laserdisc",
	"mediatype_orb":             "Orb Disk",
}

func SetGlobalFuncs(router *gin.Engine) {
	router.SetFuncMap(template.FuncMap{
		"formatAsDate": FormatAsDate,
		"add":          Add,
		"subtract":     Subtract,
		"getMediatype": GetMediatype,
	})
}
