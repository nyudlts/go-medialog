package database

import (
	"log"

	"github.com/nyudlts/go-medialog/models"
)

func InsertToken(apiToken *models.Token) error {
	if err := db.Create(apiToken).Error; err != nil {
		return err
	}
	return nil
}

func UpdateToken(apiToken *models.Token) error {
	if err := db.Save(apiToken).Error; err != nil {
		return err
	}
	return nil
}

func FindToken(token string) (models.Token, error) {
	apiToken := models.Token{}
	if err := db.Table("tokens").Where("token = ?", token).First(&apiToken).Error; err != nil {
		return apiToken, err
	}
	return apiToken, nil
}

func FindTokenByID(id uint) (models.Token, error) {
	apiToken := models.Token{}
	if err := db.Table("tokens").Where("id = ?", id).First(&apiToken).Error; err != nil {
		return models.Token{}, err
	}
	return apiToken, nil
}

func GetTokens() []models.Token {
	tokens := []models.Token{}
	db.Find(&tokens)
	return tokens
}

func ExpireToken(id uint) error {
	token, err := FindTokenByID(id)
	if err != nil {
		return err
	}

	log.Printf("[INFO] %v", token)

	token.IsValid = false

	if err := UpdateToken(&token); err != nil {
		return err
	}
	return nil

}
