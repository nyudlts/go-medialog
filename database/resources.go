package database

import (
	"log"

	"github.com/nyudlts/go-medialog/models"
	"github.com/nyudlts/go-medialog/utils"
	"gorm.io/gorm/clause"
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
	if err := db.Preload(clause.Associations).Where("id = ?", id).First(&collection).Error; err != nil {
		return collection, err
	}

	log.Println("%v\n", collection)
	return collection, nil
}

func FindResourcesByRepositoryID(repositoryID uint) ([]models.Collection, error) {
	collections := []models.Collection{}
	if err := db.Where("repository_id = ?", repositoryID).Find(&collections).Error; err != nil {
		return collections, err
	}
	return collections, nil
}

func FindPaginatedResources(pagination utils.Pagination) ([]models.Collection, error) {
	resources := []models.Collection{}
	if err := db.Limit(pagination.Limit).Offset(pagination.Offset).Order(pagination.Sort).Find(&resources).Error; err != nil {
		return resources, err
	}
	return resources, nil
}
