package database

import "github.com/nyudlts/go-medialog/models"

func FindAccessions() []models.Accession {
	accessions := []models.Accession{}
	db.Find(&accessions)
	return accessions
}

func FindAccession(id int) (models.Accession, error) {
	accession := models.Accession{}

	if err := db.Where("id = ?", id).First(&accession).Error; err != nil {
		return accession, err
	}
	return accession, nil
}
