package database

import "github.com/nyudlts/go-medialog/models"

func CreateRepository(repository models.Repository) error {
	if err := db.Create(&repository).Error; err != nil {
		return err
	}
	return nil
}

func FindRepositories() ([]models.Repository, error) {
	repositories := []models.Repository{}
	if err := db.Find(&repositories).Error; err != nil {
		return repositories, err
	}
	return repositories, nil
}

func FindRepository(id int) (models.Repository, error) {
	repository := models.Repository{}
	if err := db.Where("id = ?", id).First(&repository).Error; err != nil {
		return repository, err
	}
	return repository, nil
}
