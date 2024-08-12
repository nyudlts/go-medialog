package database

import (
	"fmt"

	"github.com/google/uuid"
	"github.com/nyudlts/bytemath"
	"github.com/nyudlts/go-medialog/models"
	"gorm.io/gorm/clause"
)

func InsertEntry(entry *models.Entry) (uuid.UUID, error) {
	if err := db.Create(&entry).Error; err != nil {
		fakeUUID, _ := uuid.NewUUID()
		return fakeUUID, err
	}
	return entry.ID, nil
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

func FindEntriesByResourceID(id uint, pagination Pagination) ([]models.Entry, error) {
	entries := []models.Entry{}
	if err := db.Preload(clause.Associations).Where("resource_id = ?", id).Limit(pagination.Limit).Offset(pagination.Offset).Order(pagination.Sort).Find(&entries).Error; err != nil {
		return entries, err
	}
	return entries, nil
}

func FindEntriesByAccessionID(id uint, pagination Pagination) ([]models.Entry, error) {
	entries := []models.Entry{}
	if err := db.Where("accession_id = ?", id).Limit(pagination.Limit).Offset(pagination.Offset).Order(pagination.Sort).Find(&entries).Error; err != nil {
		return entries, err
	}
	return entries, nil
}

func FindEntry(id uuid.UUID) (models.Entry, error) {
	entry := models.Entry{}
	if err := db.Preload(clause.Associations).Where("id = ?", id).First(&entry).Error; err != nil {
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

func FindMaxMediaIDInResource(resourceID uint) int {
	var maxMediaID int
	db.Table("entries").Where("resource_id = ?", resourceID).Order("media_id desc").Select("media_id").Limit(1).Find(&maxMediaID)
	return maxMediaID
}

func FindPaginatedEntries(pagination Pagination) ([]models.Entry, error) {
	entries := []models.Entry{}
	if err := db.Preload(clause.Associations).Limit(pagination.Limit).Offset(pagination.Offset).Order(pagination.Sort).Find(&entries).Error; err != nil {
		return entries, err
	}
	return entries, nil
}

func GetNumberPagesInResource(resourceID uint) (int, error) {
	entryIDs := []uuid.UUID{}
	if err := db.Table("entries").Where("resource_id = ?", resourceID).Select("id").Find(&entryIDs).Error; err != nil {
		return 0, err
	}

	l := len(entryIDs)
	p := l / 10

	return p, nil
}

func GetCountOfEntriesInDB() int64 {
	var count int64
	db.Model(&models.Entry{}).Count(&count)
	return count
}

func GetCountOfEntriesInAccession(accessionID uint) int64 {
	var count int64
	db.Model(&models.Entry{}).Where("accession_id = ?", accessionID).Count(&count)
	return count
}

func GetCountOfEntriesInResource(resourceID uint) int64 {
	var count int64
	db.Model(&models.Entry{}).Where("resource_id = ?", resourceID).Count(&count)
	return count
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
	if err := db.Where("resource_id = ?", id).Find(&entries).Error; err != nil {
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
	StartYear    int `form:"start-year"`
	StartMonth   int `form:"start-month"`
	StartDay     int `form:"start-day"`
	EndYear      int `form:"end-year"`
	EndMonth     int `form:"end-month"`
	EndDay       int `form:"end-day"`
	RepositoryID int `form:"partner_code"`
}

func (dr DateRange) String() string {
	return fmt.Sprintf("%d-%d-%d to %d-%d-%d", dr.StartYear, dr.StartMonth, dr.StartDay, dr.EndYear, dr.EndMonth, dr.EndDay)
}

func GetSummaryByDateRange(dr DateRange) (Summaries, error) {
	startDate := fmt.Sprintf("%d-%d-%dT00:00:00Z", dr.StartYear, dr.StartMonth, dr.StartDay)
	endDate := fmt.Sprintf("%d-%d-%dT23:59:59Z", dr.EndYear, dr.EndMonth, dr.EndDay)

	entries := []models.Entry{}
	if dr.RepositoryID == 0 {
		if err := db.Where("created_at BETWEEN ? AND ?", startDate, endDate).Find(&entries).Error; err != nil {
			return Summaries{}, err
		}
		return getSummary(entries), nil
	} else {
		if err := db.Where("repository_id = ?", dr.RepositoryID).Where("created_at BETWEEN ? AND ?", startDate, endDate).Find(&entries).Error; err != nil {
			return Summaries{}, err
		}
		return getSummary(entries), nil
	}
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

func FindEntryByMediaIDAndCollectionID(mediaID uint, ResourceID uint) (uuid.UUID, error) {
	entry := models.Entry{}
	if err := db.Where("media_id = ? AND resource_id = ?", mediaID, ResourceID).First(&entry).Error; err != nil {
		return uuid.New(), err
	}
	return entry.ID, nil
}

func FindNextMediaCollectionInResource(resourceID uint) (uint, error) {
	var entry models.Entry
	if err := db.Where("resource_id = ?", resourceID).Order("media_id desc").First(&entry).Error; err != nil {
		if (entry == models.Entry{}) {
			return 1, nil
		} else {
			return 0, err
		}
	}

	return entry.MediaID + 1, nil
}

func IsMediaIDUniqueInResource(mediaID uint, resourceID uint) (bool, error) {
	fmt.Println("TEST", mediaID, resourceID)
	entries := []models.Entry{}

	if err := db.Where("resource_id = ?", int(resourceID)).Find(&entries).Error; err != nil {
		return false, err
	}

	fmt.Println("TEST", entries)

	for _, entry := range entries {
		if entry.MediaID == mediaID {
			return false, nil
		}
	}

	return true, nil
}

func FindEntryInResource(resourceID int, mediaID int) (string, error) {
	entry := models.Entry{}
	if err := db.Where("resource_id = ? AND media_id = ?", resourceID, mediaID).First(&entry).Error; err != nil {
		return "", err
	}
	return entry.ID.String(), nil
}

func GetEntryIDs() ([]string, error) {
	ids := []string{}
	if err := db.Table("entries").Select("id").Find(&ids).Error; err != nil {
		return []string{}, err
	}
	return ids, nil
}

func GetEntryIDsPaginated(pagination Pagination) ([]string, error) {
	ids := []string{}
	if err := db.Table("entries").Select("id").Limit(pagination.Limit).Offset(pagination.Offset).Find(&ids).Error; err != nil {
		return []string{}, err
	}
	return ids, nil
}

func FindEntriesPaginated(pagination Pagination) ([]models.Entry, error) {
	entries := []models.Entry{}
	if err := db.Preload(clause.Associations).Limit(pagination.Limit).Offset(pagination.Offset).Find(&entries).Error; err != nil {
		return entries, err
	}
	return entries, nil
}
