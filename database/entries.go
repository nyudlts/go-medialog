package database

import "github.com/nyudlts/go-medialog/models"

func FindEntriesByResourceID(id uint) ([]models.Entry, error) {
	entries := []models.Entry{}
	if err := db.Where("collection_id = ?", id).Find(&entries).Error; err != nil {
		return entries, err
	}
	return entries, nil
}

func FindEntriesByAccessionID(id uint) ([]models.Entry, error) {
	entries := []models.Entry{}
	if err := db.Where("accession_id = ?", id).Find(&entries).Error; err != nil {
		return entries, err
	}
	return entries, nil
}
