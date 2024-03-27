package database

import (
	"github.com/nyudlts/go-medialog/models"
)

func FindResources() ([]models.Collection, error) {
	collections := []models.Collection{}
	if err := db.Find(&collections).Error; err != nil {
		return collections, err
	}
	return collections, nil
}

func FindResource(id uint) (models.Collection, error) {
	collection := models.Collection{}
	if err := db.Where("id = ?", id).First(&collection).Error; err != nil {
		return collection, err
	}
	return collection, nil
}

func FindResourcesByRepositoryID(repositoryID uint) ([]models.Collection, error) {
	collections := []models.Collection{}
	if err := db.Where("repository_id = ?", repositoryID).Find(&collections).Error; err != nil {
		return collections, err
	}
	return collections, nil
}
