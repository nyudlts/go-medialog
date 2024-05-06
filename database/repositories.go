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

func FindRepository(id uint) (models.Repository, error) {
	repository := models.Repository{}
	if err := db.Where("id = ?", id).First(&repository).Error; err != nil {
		return repository, err
	}
	return repository, nil
}

func GetRepositoryMap() (map[int]string, error) {
	repositories, err := FindRepositories()
	if err != nil {
		return map[int]string{}, err
	}
	repositoryMap := map[int]string{}
	for _, repo := range repositories {
		repositoryMap[int(repo.ID)] = repo.Slug
	}
	return repositoryMap, nil
}
