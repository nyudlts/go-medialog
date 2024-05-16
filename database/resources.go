package database

import (
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

func InsertResource(resource *models.Collection) (uint, error) {
	if err := db.Debug().Create(resource).Error; err != nil {
		return 0, err
	}
	return resource.ID, nil
}

func DeleteResource(id uint) error {
	if err := db.Delete(models.Collection{}, id).Error; err != nil {
		return err
	}
	return nil
}

func UpdateResource(resource *models.Collection) error {
	if err := db.Save(resource).Error; err != nil {
		return err
	}
	return nil
}
