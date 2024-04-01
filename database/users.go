package database

import "github.com/nyudlts/go-medialog/models"

func FindUserByID(id int) (models.User, error) {
	user := models.User{}
	if err := db.Where("id = ?", id).First(&user).Error; err != nil {
		return user, err
	}
	return user, nil
}

func UpdateUser(user models.User) error {
	if err := db.Save(&user).Error; err != nil {
		return err
	}
	return nil
}

func FindUsers() ([]models.User, error) {
	users := []models.User{}
	if err := db.Order("email").Find(&users).Error; err != nil {
		return users, err
	}
	return users, nil
}

func FindUserByEmail(email string) (models.User, error) {
	user := models.User{}
	if err := db.Where("email = ?", email).First(&user).Error; err != nil {
		return user, err
	}
	return user, nil
}

func InsertUser(user *models.User) error {
	if err := db.Create(&user).Error; err != nil {
		return err
	}
	return nil
}
