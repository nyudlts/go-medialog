package database

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/google/uuid"
	"github.com/nyudlts/bytemath"
	"github.com/nyudlts/go-medialog/models"
	"gorm.io/gorm/clause"
)

func InsertEntry(entry *models.Entry) error {
	if err := db.Create(&entry).Error; err != nil {
		return err
	}

	if err := InsertEntryJSON(*entry); err != nil {
		return err
	}

	return nil
}

func DeleteEntry(id uuid.UUID) error {
	ej, err := FindEntryJSONByEntryID(id)
	if err != nil {
		return err
	}

	if err := DeleteEntryJSON(ej.ID); err != nil {
		return err
	}

	if err := db.Delete(models.Entry{}, id).Error; err != nil {
		return err
	}

	return nil
}

func UpdateEntry(entry *models.Entry) error {
	if err := db.Save(entry).Error; err != nil {
		return err
	}

	ej, err := FindEntryJSONByEntryID(entry.ID)
	if err != nil {
		return err
	}

	em := entry.Minimal()
	emBytes, err := json.Marshal(em)
	if err != nil {
		return err
	}
	ej.JSON = string(emBytes)
	ej.EntryID = entry.ID

	if err := UpdateEntryJSON(ej); err != nil {
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

func FindEntriesFiltered(filter string) ([]models.Entry, error) {
	entries := []models.Entry{}
	if filter == "" {
		return FindEntries()
	} else {
		if err := db.Table("entries").Where("mediatype = ?", filter).Find(&entries).Error; err != nil {
			return entries, err
		}
		return entries, nil
	}
}

func FindEntryIDsByResourceID(id uint) ([]string, error) {
	entries := []string{}
	if err := db.Table("entries").Where("resource_id = ?", id).Select("id").Find(&entries).Error; err != nil {
		return entries, err
	}
	return entries, nil
}

func FindEntriesByResourceID(id uint) ([]models.Entry, error) {
	entries := []models.Entry{}
	if err := db.Preload(clause.Associations).Where("resource_id = ?", id).Find(&entries).Error; err != nil {
		return []models.Entry{}, err
	}
	return entries, nil
}

func FindEntriesByResourceIDFiltered(id uint, filter string) ([]models.Entry, error) {
	entries := []models.Entry{}
	if filter == "" {
		return FindEntriesByResourceID(id)
	} else {
		if err := db.Preload(clause.Associations).Where("resource_id = ? AND mediatype = ?", id, filter).Find(&entries).Error; err != nil {
			return []models.Entry{}, err
		}
		return entries, nil
	}
}

func FindEntriesByResourceIDPaginated(id uint, pagination Pagination) ([]models.Entry, error) {
	entries := []models.Entry{}
	if pagination.Filter == "" {
		if err := db.Preload(clause.Associations).Where("resource_id = ?", id).Limit(pagination.Limit).Offset(pagination.Offset).Order(pagination.Sort).Find(&entries).Error; err != nil {
			return entries, err
		}
	} else {
		if err := db.Preload(clause.Associations).Where("resource_id = ? AND mediatype = ?", id, pagination.Filter).Limit(pagination.Limit).Offset(pagination.Offset).Order(pagination.Sort).Find(&entries).Error; err != nil {
			return entries, err
		}
	}
	return entries, nil
}

func FindEntryIDsByAccessionID(id uint) ([]string, error) {
	ids := []string{}
	if err := db.Table("entries").Where("accession_id = ?", id).Select("id").Find(&ids).Error; err != nil {
		return []string{}, err
	}
	return ids, nil
}

func FindEntriesByAccessionIDPaginated(id uint, pagination Pagination) ([]models.Entry, error) {
	entries := []models.Entry{}
	if pagination.Filter == "" {
		if err := db.Where("accession_id = ?", id).Limit(pagination.Limit).Offset(pagination.Offset).Order(pagination.Sort).Find(&entries).Error; err != nil {
			return entries, err
		}
	} else {
		if err := db.Where("accession_id = ? AND mediatype = ?", id, pagination.Filter).Limit(pagination.Limit).Offset(pagination.Offset).Order(pagination.Sort).Find(&entries).Error; err != nil {
			return entries, err
		}
	}

	return entries, nil
}

func FindEntriesByAccessionID(id uint) ([]models.Entry, error) {
	entries := []models.Entry{}
	if err := db.Preload(clause.Associations).Where("accession_id = ?", id).Find(&entries).Error; err != nil {
		return entries, err
	}
	return entries, nil
}

func FindEntriesByAccessionIDFiltered(id uint, filter string) ([]models.Entry, error) {

	if filter == "" {
		return FindEntriesByAccessionID(id)
	} else {
		entries := []models.Entry{}
		if err := db.Preload(clause.Associations).Where("accession_id = ? AND mediatype = ?", id).Find(&entries).Error; err != nil {
			return entries, err
		}
		return entries, nil
	}

}

func FindEntriesByRepositoryID(repositoryID uint) ([]models.Entry, error) {
	entries := []models.Entry{}
	if err := db.Preload(clause.Associations).Where("repository_id = ?", repositoryID).Find(&entries).Error; err != nil {
		return []models.Entry{}, err
	}
	return entries, nil
}

func FindEntriesByRepositoryIDPaginated(repositoryID uint, pagination Pagination) ([]models.Entry, error) {
	entries := []models.Entry{}
	if pagination.Filter == "" {
		if err := db.Where("repository_id = ?", repositoryID).Limit(pagination.Limit).Offset(pagination.Offset).Order(pagination.Sort).Find(&entries).Error; err != nil {
			return entries, err
		}
	} else {
		if err := db.Where("repository_id = ? AND mediatype = ?", repositoryID, pagination.Filter).Limit(pagination.Limit).Offset(pagination.Offset).Order(pagination.Sort).Find(&entries).Error; err != nil {
			return entries, err
		}
	}
	return entries, nil
}

func FindEntryIDsByRepositoryID(repositoryID uint) ([]string, error) {
	ids := []string{}
	if err := db.Table("entries").Where("repository_id = ?", repositoryID).Select("id").Find(&ids).Error; err != nil {
		return []string{}, err
	}
	return ids, nil
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

	if pagination.Filter == "" {
		if err := db.Preload(clause.Associations).Limit(pagination.Limit).Offset(pagination.Offset).Order(pagination.Sort).Find(&entries).Error; err != nil {
			return entries, err
		}
	} else {
		log.Println("FILTER", pagination.Filter)
		if err := db.Preload(clause.Associations).Where("mediatype = ?", pagination.Filter).Limit(pagination.Limit).Offset(pagination.Offset).Order(pagination.Sort).Find(&entries).Error; err != nil {
			return entries, err
		}
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

func GetCountOfEntriesInDBPaginated(pagination *Pagination) int64 {
	var count int64
	if pagination.Filter != "" {
		db.Table("entries").Where("mediatype = ?", pagination.Filter).Count(&count)
		return count
	} else {
		return GetCountOfEntriesInDB()
	}
}

func GetCountOfEntriesInAccession(accessionID uint) int64 {
	var count int64
	db.Model(&models.Entry{}).Where("accession_id = ?", accessionID).Count(&count)
	return count
}

func GetCountOfEntriesInAccessionPaginated(accessionID uint, pagination *Pagination) int64 {
	var count int64
	if pagination.Filter != "" {
		db.Table("entries").Where("mediatype = ? AND accession_id = ?", pagination.Filter, accessionID).Count(&count)
		return count
	} else {
		return GetCountOfEntriesInAccession(accessionID)
	}
}

func GetCountOfEntriesInResource(resourceID uint) int64 {
	var count int64
	db.Model(&models.Entry{}).Where("resource_id = ?", resourceID).Count(&count)
	return count
}

func GetCountOfEntriesInResourcePaginated(resourceID uint, pagination Pagination) int64 {
	if pagination.Filter != "" {
		var count int64
		db.Table("entries").Where("mediatype = ? AND resource_id = ?", pagination.Filter, resourceID).Count(&count)
		return count
	} else {
		return GetCountOfEntriesInResource(resourceID)
	}
}

func GetCountOfEntriesInRepository(repositoryID uint) int64 {
	var count int64
	db.Model(&models.Entry{}).Where("repository_id = ?", repositoryID).Count(&count)
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

func (s Summaries) GetSlice() []Summary {
	summarySlice := []Summary{}
	for _, v := range s {
		summarySlice = append(summarySlice, v)
	}
	return summarySlice
}

func (s Summaries) GetTotals() Totals {
	totals := Totals{}
	for _, summary := range s {
		totals.Count += summary.Count
		totals.Size += summary.Size
	}
	totals.HumanSize = bytemath.ConvertBytesToHumanReadable(int64(totals.Size))
	return totals
}

func GetSummaryByRepository(repositoryID uint) (Summaries, error) {
	entries := []models.Entry{}
	if err := db.Where("repository_id = ?", repositoryID).Find(&entries).Error; err != nil {
		return Summaries{}, err
	}
	return getSummary(entries), nil
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

func (dr DateRange) String() string {
	return fmt.Sprintf("%d-%d-%d to %d-%d-%d", dr.StartYear, dr.StartMonth, dr.StartDay, dr.EndYear, dr.EndMonth, dr.EndDay)
}

func GetEntriesByDateRange(dr DateRange) ([]models.Entry, error) {
	startDate := fmt.Sprintf("%d-%d-%dT00:00:00Z", dr.StartYear, dr.StartMonth, dr.StartDay)
	endDate := fmt.Sprintf("%d-%d-%dT23:59:59Z", dr.EndYear, dr.EndMonth, dr.EndDay)
	entries := []models.Entry{}
	if dr.RepositoryID == 0 {
		if err := db.Preload(clause.Associations).Where("created_at BETWEEN ? AND ?", startDate, endDate).Find(&entries).Error; err != nil {
			return entries, err
		}
		return entries, nil
	} else {
		if err := db.Preload(clause.Associations).Where("repository_id = ?", dr.RepositoryID).Where("created_at BETWEEN ? AND ?", startDate, endDate).Find(&entries).Error; err != nil {
			return entries, err
		}
		return entries, nil
	}

}

func GetSummaryByDateRange(dr DateRange) (Summaries, error) {

	startDate := fmt.Sprintf("%d-%d-%dT00:00:00Z", dr.StartYear, dr.StartMonth, dr.StartDay)
	endDate := fmt.Sprintf("%d-%d-%dT23:59:59Z", dr.EndYear, dr.EndMonth, dr.EndDay)

	entries := []models.Entry{}
	//this needs to be simplified
	if dr.IsRefreshed {
		if dr.RepositoryID == 0 {
			if err := db.Where("created_at BETWEEN ? AND ?", startDate, endDate).Where("is_refreshed = true").Find(&entries).Error; err != nil {
				return Summaries{}, err
			}
		} else {
			if err := db.Where("repository_id = ?", dr.RepositoryID).Where("created_at BETWEEN ? AND ?", startDate, endDate).Where("is_refreshed = true").Find(&entries).Error; err != nil {
				return Summaries{}, err
			}
		}
	} else {
		if dr.RepositoryID == 0 {
			if err := db.Where("created_at BETWEEN ? AND ?", startDate, endDate).Find(&entries).Error; err != nil {
				return Summaries{}, err
			}
		} else {
			if err := db.Where("repository_id = ?", dr.RepositoryID).Where("created_at BETWEEN ? AND ?", startDate, endDate).Find(&entries).Error; err != nil {
				return Summaries{}, err
			}
		}
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

func FindEntryByMediaIDAndCollectionID(mediaID uint, ResourceID uint) (uuid.UUID, error) {
	entry := models.Entry{}
	if err := db.Where("media_id = ? AND resource_id = ?", mediaID, ResourceID).First(&entry).Error; err != nil {
		return uuid.New(), err
	}
	return entry.ID, nil
}

func FindNextMediaCollectionInResource(resourceID uint) (uint, error) {

	var entries = []models.Entry{}

	if err := db.Where("resource_id = ?", resourceID).Order("media_id desc").Find(&entries).Error; err != nil {
		return 0, err
	}

	if (len(entries)) == 0 {
		return 1, nil
	}

	return entries[0].MediaID + 1, nil
}

func IsMediaIDUniqueInResource(mediaID uint, resourceID uint) (bool, error) {

	entries := []models.Entry{}

	if err := db.Where("resource_id = ?", int(resourceID)).Find(&entries).Error; err != nil {
		return false, err
	}

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

func getEntryIDs() ([]uuid.UUID, error) {
	entryIDs := []uuid.UUID{}
	if err := db.Table("entries").Select("id").Scan(&entryIDs).Error; err != nil {
		return entryIDs, err
	}
	return entryIDs, nil
}
