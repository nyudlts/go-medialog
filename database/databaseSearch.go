package database

import (
	"github.com/nyudlts/go-medialog/models"
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
	if err := db.Where("title LIKE ?", "%"+query+"%").Find(&resources).Error; err != nil {
		return resources, err
	}
	return resources, nil
}
