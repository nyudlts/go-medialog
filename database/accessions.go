package database

import (
	"github.com/nyudlts/go-medialog/models"
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

func FindAccession(id uint) (models.Accession, error) {
	accession := models.Accession{}

	if err := db.Preload(clause.Associations).Where("id = ?", id).First(&accession).Error; err != nil {
		return accession, err
	}
	return accession, nil
}

func FindPaginatedAccessions(pagination Pagination) ([]models.Accession, error) {
	accessions := []models.Accession{}
	if err := db.Limit(pagination.Limit).Offset(pagination.Offset).Order(pagination.Sort).Find(&accessions).Error; err != nil {
		return accessions, err
	}
	return accessions, nil
}

func InsertAccession(accession *models.Accession) (uint, error) {
	if err := db.Create(accession).Debug().Error; err != nil {
		return 0, err
	}

	return accession.ID, nil
}

func UpdateAccession(accession *models.Accession) error {
	if err := db.Save(accession).Error; err != nil {
		return err
	}
	return nil
}

func DeleteAccession(id uint) error {
	if err := db.Delete(models.Accession{}, id).Error; err != nil {
		return err
	}
	return nil
}
