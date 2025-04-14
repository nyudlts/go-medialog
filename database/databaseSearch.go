package database

import (
	"github.com/google/uuid"
	"github.com/nyudlts/go-medialog/models"
	"gorm.io/gorm/clause"
)

func SearchRepositories(query string) ([]models.Repository, error) {
	repositories := []models.Repository{}
	if err := db.Where("title LIKE ?", "%"+query+"%").Find(&repositories).Error; err != nil {
		return repositories, err
	}
	return repositories, nil
}

func SearchResources(query string) ([]models.Resource, error) {
	resources := []models.Resource{}
	if err := db.Preload(clause.Associations).Table("resources").Where("title LIKE ? OR collection_code LIKE ?", "%"+query+"%", "%"+query+"%").Scan(&resources).Error; err != nil {
		return resources, err
	}
	return resources, nil
}

func SearchAccessions(query string) ([]models.Accession, error) {
	accessions := []models.Accession{}
	if err := db.Preload(clause.Associations).Where("accession_num LIKE ?", "%"+query+"%").Find(&accessions).Error; err != nil {
		return accessions, err
	}
	return accessions, nil
}

func SearchEntries(query string) ([]models.Entry, error) {
	hitIds := []uuid.UUID{}
	entries := []models.Entry{}
	if err := db.Table("entry_jsons").Select("entry_id").Where("json LIKE ? AND deleted_at IS NULL", "%"+query+"%").Scan(&hitIds).Error; err != nil {
		return entries, err
	}

	for _, hitID := range hitIds {
		entry, err := FindEntry(hitID)
		if err != nil {
			return entries, err
		}
		entries = append(entries, entry)
	}
	return entries, nil
}
