package database

import (
	"github.com/nyudlts/go-medialog/models"
	"github.com/nyudlts/go-medialog/utils"
	"gorm.io/gorm/clause"
)

func FindAccessions() ([]models.Accession, error) {
	accessions := []models.Accession{}
	if err := db.Preload(clause.Associations).Find(&accessions).Error; err != nil {
		return accessions, err
	}
	return accessions, nil
}

func FindAccessionsByResourceID(id uint) ([]models.Accession, error) {
	accessions := []models.Accession{}
	if err := db.Where("collection_id = ?", id).Find(&accessions).Error; err != nil {
		return accessions, err
	}
	return accessions, nil
}

func FindAccession(id int) (models.Accession, error) {
	accession := models.Accession{}

	if err := db.Preload(clause.Associations).Where("id = ?", id).First(&accession).Error; err != nil {
		return accession, err
	}
	return accession, nil
}

func FindPaginatedAccessions(pagination utils.Pagination) ([]models.Accession, error) {
	accessions := []models.Accession{}
	if err := db.Limit(pagination.Limit).Offset(pagination.Offset).Order(pagination.Sort).Find(&accessions).Error; err != nil {
		return accessions, err
	}
	return accessions, nil
}

func UpdateAccession(accession *models.Accession) error {
	if err := db.Save(accession).Error; err != nil {
		return err
	}
	return nil
}

func DeleteAccession(id int) error {
	if err := db.Delete(models.Accession{}, id).Error; err != nil {
		return err
	}
	return nil
}
