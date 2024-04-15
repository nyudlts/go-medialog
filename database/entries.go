package database

import (
	"github.com/nyudlts/bytemath"
	"github.com/nyudlts/go-medialog/models"
	"github.com/nyudlts/go-medialog/utils"
)

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

func FindEntry(id string) (models.Entry, error) {
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
