package database

import "github.com/nyudlts/go-medialog/models"

func InsertSessionToken(token *models.SessionToken) error {
	if err := db.Create(token).Error; err != nil {
		return err
	}
	return nil
}

func FindSessionToken(token string) (models.SessionToken, error) {
	sessionToken := models.SessionToken{}
	if err := db.Table("session_tokens").Where("token = ?", token).First(&sessionToken).Error; err != nil {
		return models.SessionToken{}, err
	}
	return sessionToken, nil
}
