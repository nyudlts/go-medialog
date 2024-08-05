package database

import (
	"log"

	"github.com/nyudlts/go-medialog/models"
)

func InsertAPIToken(apiToken *models.APIToken) error {
	if err := db.Create(apiToken).Error; err != nil {
		return err
	}
	return nil
}

func UpdateAPIToken(apiToken *models.APIToken) error {
	if err := db.Save(apiToken).Error; err != nil {
		return err
	}
	return nil
}

func FindToken(token string) (models.APIToken, error) {
	apiToken := models.APIToken{}
	if err := db.Table("api_tokens").Where("token = ?", token).First(&apiToken).Error; err != nil {
		return apiToken, err
	}
	return apiToken, nil
}

func FindTokenByID(id uint) (models.APIToken, error) {
	apiToken := models.APIToken{}
	if err := db.Table("api_tokens").Where("id = ?", id).First(&apiToken).Error; err != nil {
		return models.APIToken{}, err
	}
	return apiToken, nil
}

func GetTokens() []models.APIToken {
	tokens := []models.APIToken{}
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

	if err := UpdateAPIToken(&token); err != nil {
		return err
	}
	return nil

}
