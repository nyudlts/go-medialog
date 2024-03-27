package database

import (
	"log"

	"github.com/nyudlts/go-medialog/models"
)

func FindEntries() ([]models.Entry, error) {
	entries := []models.Entry{}
	if err := db.Find(&entries).Error; err != nil {
		return entries, err
	}
	return entries, nil
}

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

func FindEntry(id string) (models.Entry, error) {
	log.Println(id)
	entry := models.Entry{}
	if err := db.Where("id = ?", id).First(&entry).Error; err != nil {
		return entry, err
	}
	return entry, nil
}
