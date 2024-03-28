package database

import "github.com/nyudlts/go-medialog/models"

func FindUsers() ([]models.User, error) {
	users := []models.User{}
	if err := db.Find(&users).Error; err != nil {
		return users, err
	}
	return users, nil
}
