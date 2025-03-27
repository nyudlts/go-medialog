package database

import (
	"github.com/nyudlts/go-medialog/models"
	"gorm.io/gorm/clause"
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
	if err := db.Preload(clause.Associations).Table("resources").Where("title LIKE ? OR collection_code LIKE ?", "%"+query+"%", "%"+query+"%").Scan(&resources).Error; err != nil {
		return resources, err
	}
	return resources, nil
}

func SearchAccessions(query string) ([]models.Accession, error) {
	accessions := []models.Accession{}
	if err := db.Preload(clause.Associations).Where("accession_num LIKE ?", "%"+query+"%").Find(&accessions).Error; err != nil {
		return accessions, err
	}
	return accessions, nil
}
