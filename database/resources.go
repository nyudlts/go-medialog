package database

import (
	"github.com/nyudlts/go-medialog/models"
	"gorm.io/gorm/clause"
)

func FindResources() ([]models.Resource, error) {
	resources := []models.Resource{}
	if err := db.Find(&resources).Error; err != nil {
		return resources, err
	}
	return resources, nil
}

func FindResource(id uint) (models.Resource, error) {
	resource := models.Resource{}
	if err := db.Preload(clause.Associations).Where("id = ?", id).First(&resource).Error; err != nil {
		return resource, err
	}
	return resource, nil
}

func FindResourcesByRepositoryID(repositoryID uint) ([]models.Resource, error) {
	resources := []models.Resource{}
	if err := db.Where("repository_id = ?", repositoryID).Find(&resources).Error; err != nil {
		return resources, err
	}
	return resources, nil
}

func FindPaginatedResources(pagination Pagination) ([]models.Resource, error) {
	resources := []models.Resource{}
	if err := db.Limit(pagination.Limit).Offset(pagination.Offset).Order(pagination.Sort).Find(&resources).Error; err != nil {
		return resources, err
	}
	return resources, nil
}

func InsertResource(resource *models.Resource) (uint, error) {
	if err := db.Create(resource).Error; err != nil {
		return 0, err
	}
	return resource.ID, nil
}

func DeleteResource(id uint) error {
	if err := db.Delete(models.Resource{}, id).Error; err != nil {
		return err
	}
	return nil
}

func UpdateResource(resource *models.Resource) error {
	if err := db.Save(resource).Error; err != nil {
		return err
	}
	return nil
}
