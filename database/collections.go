package database

import "github.com/nyudlts/go-medialog/models"

func FindCollections() []models.Collection {
	collections := []models.Collection{}
	db.Find(&collections)
	return collections
}

func FindCollection(id int) (models.Collection, error) {
	collection := models.Collection{}
	if err := db.Where("id = ?", id).First(&collection).Error; err != nil {
		return collection, err
	}
	return collection, nil
}
