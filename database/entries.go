package database

import (
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
