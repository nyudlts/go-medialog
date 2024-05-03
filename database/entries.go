package database

import (
	"fmt"

	"github.com/google/uuid"
	"github.com/nyudlts/bytemath"
	"github.com/nyudlts/go-medialog/models"
	"github.com/nyudlts/go-medialog/utils"
)

func InsertEntry(entry models.Entry) error {
	if err := db.Create(&entry).Error; err != nil {
		return err
	}
	return nil
}

func DeleteEntry(id uuid.UUID) error {
	if err := db.Delete(models.Entry{}, id).Error; err != nil {
		return err
	}
	return nil
}

func UpdateEntry(entry *models.Entry) error {
	if err := db.Save(entry).Error; err != nil {
		return err
	}
	return nil
}

func FindEntries() ([]models.Entry, error) {
	entries := []models.Entry{}
	if err := db.Find(&entries).Error; err != nil {
		return entries, err
	}
	return entries, nil
}

func FindEntriesByResourceID(id uint, pagination utils.Pagination) ([]models.Entry, error) {
	entries := []models.Entry{}
	if err := db.Where("collection_id = ?", id).Limit(pagination.Limit).Offset(pagination.Offset).Order(pagination.Sort).Find(&entries).Error; err != nil {
		return entries, err
	}
	return entries, nil
}

func FindEntriesByAccessionID(id uint, pagination utils.Pagination) ([]models.Entry, error) {
	entries := []models.Entry{}
	if err := db.Where("accession_id = ?", id).Limit(pagination.Limit).Offset(pagination.Offset).Order(pagination.Sort).Find(&entries).Error; err != nil {
		return entries, err
	}
	return entries, nil
}

func FindEntry(id uuid.UUID) (models.Entry, error) {
	entry := models.Entry{}
	if err := db.Where("id = ?", id).First(&entry).Error; err != nil {
		return entry, err
	}
	return entry, nil
}

func FindEntriesSorted(numRecords int) ([]models.Entry, error) {
	entries := []models.Entry{}
	if err := db.Limit(numRecords).Order("updated_at DESC").Find(&entries).Error; err != nil {
		return entries, err
	}
	return entries, nil
}

func FindPaginatedEntries(pagination utils.Pagination) ([]models.Entry, error) {
	entries := []models.Entry{}
	if err := db.Limit(pagination.Limit).Offset(pagination.Offset).Order(pagination.Sort).Find(&entries).Error; err != nil {
		return entries, err
	}
	return entries, nil
}

type Summary struct {
	Mediatype string
	Count     int
	Size      float64
	HumanSize string
}

type Totals struct {
	Count     int
	Size      float64
	HumanSize string
}

type Summaries map[string]Summary

func (s Summaries) GetTotals() Totals {
	totals := Totals{}
	for _, summary := range s {
		totals.Count += summary.Count
		totals.Size += summary.Size
	}
	totals.HumanSize = bytemath.ConvertBytesToHumanReadable(int64(totals.Size))
	return totals
}

func GetSummaryByResource(id uint) (Summaries, error) {
	entries := []models.Entry{}
	if err := db.Where("collection_id = ?", id).Find(&entries).Error; err != nil {
		return Summaries{}, err
	}
	return getSummary(entries), nil
}

func GetSummaryByAccession(id uint) (Summaries, error) {
	entries := []models.Entry{}
	if err := db.Where("accession_id = ?", id).Find(&entries).Error; err != nil {
		return Summaries{}, err
	}
	return getSummary(entries), nil
}

func GetSummaryByYear(year int) (Summaries, error) {
	startDate := fmt.Sprintf("%d-01-01", year)
	endDate := fmt.Sprintf("%d-01-01", year+1)
	entries := []models.Entry{}
	if err := db.Where("created_at BETWEEN ? AND ?", startDate, endDate).Find(&entries).Error; err != nil {
		return Summaries{}, err
	}
	return getSummary(entries), nil
}

type DateRange struct {
	StartYear  int `form:"start-year"`
	StartMonth int `form:"start-month"`
	StartDay   int `form:"start-day"`
	EndYear    int `form:"end-year"`
	EndMonth   int `form:"end-month"`
	EndDay     int `form:"end-day"`
}

func (dr DateRange) String() string {
	return fmt.Sprintf("%d-%d-%d to %d-%d-%d", dr.StartYear, dr.StartMonth, dr.StartDay, dr.EndYear, dr.EndMonth, dr.EndDay)
}

func GetSummaryByDateRange(dr DateRange) (Summaries, error) {
	startDate := fmt.Sprintf("%d-%d-%d", dr.StartYear, dr.StartMonth, dr.StartDay)
	endDate := fmt.Sprintf("%d-%d-%d", dr.EndYear, dr.EndMonth, dr.EndDay)
	entries := []models.Entry{}
	if err := db.Where("created_at BETWEEN ? AND ?", startDate, endDate).Find(&entries).Error; err != nil {
		return Summaries{}, err
	}
	return getSummary(entries), nil
}

func summaryContains(summaries Summaries, mediatype string) bool {
	for k, _ := range summaries {
		if k == mediatype {
			return true
		}
	}
	return false
}

func getSummary(entries []models.Entry) Summaries {
	summaries := Summaries{}
	for _, entry := range entries {
		if summaryContains(summaries, entry.Mediatype) {
			s := summaries[entry.Mediatype]
			s.Count += 1
			f64 := float64(entry.StockSizeNum)
			sfx := bytemath.GetSuffixByString(entry.StockUnit)
			s.Size = s.Size + bytemath.ConvertToBytes(f64, *sfx)
			s.HumanSize = bytemath.ConvertBytesToHumanReadable(int64(s.Size))
			summaries[entry.Mediatype] = s
		} else {
			s := Summary{Count: 1}
			s.Mediatype = entry.Mediatype
			f64 := float64(entry.StockSizeNum)
			sfx := bytemath.GetSuffixByString(entry.StockUnit)
			s.Size = bytemath.ConvertToBytes(f64, *sfx)
			s.HumanSize = bytemath.ConvertBytesToHumanReadable(int64(s.Size))
			summaries[entry.Mediatype] = s
		}
	}
	return summaries
}

func FindEntryByMediaIDAndCollectionID(mediaID int, collectionID int) (uuid.UUID, error) {
	entry := models.Entry{}
	if err := db.Where("media_id = ? AND collection_id = ?", mediaID, collectionID).First(&entry).Error; err != nil {
		return uuid.New(), err
	}
	return entry.ID, nil
}

func FindNextMediaCollectionInResource(resourceID uint) (int, error) {
	var entry models.Entry
	if err := db.Where("collection_id = ?", resourceID).Order("media_id desc").First(&entry).Error; err != nil {
		return 0, err
	}

	return entry.MediaID + 1, nil
}
