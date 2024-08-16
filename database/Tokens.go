package database

import (
	"github.com/nyudlts/go-medialog/models"
	"gorm.io/gorm/clause"
)

func InsertToken(apiToken *models.Token) error {
	if err := db.Preload(clause.Associations).Create(apiToken).Error; err != nil {
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
	if err := db.Preload(clause.Associations).Table("tokens").Where("token = ?", token).First(&apiToken).Error; err != nil {
		return apiToken, err
	}
	return apiToken, nil
}

func FindTokenByID(id uint) (models.Token, error) {
	apiToken := models.Token{}
	if err := db.Table("tokens").Preload(clause.Associations).Where("id = ?", id).First(&apiToken).Error; err != nil {
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

	token.IsValid = false

	if err := UpdateToken(&token); err != nil {
		return err
	}
	return nil

}

func ExpireTokensByUserID(userID uint) error {
	tokens, err := FindTokensByUserID(userID)
	if err != nil {
		return err
	}

	for _, token := range tokens {
		if err := ExpireToken(token.ID); err != nil {
			return err
		}
	}

	return nil
}

func ExpireAPITokensByUserID(userID uint) error {
	tokens, err := FindTokensByUserID(userID)
	if err != nil {
		return err
	}

	for _, token := range tokens {
		if token.Type == "api" {
			if err := ExpireToken(token.ID); err != nil {
				return err
			}
		}
	}

	return nil
}

func ExpireAppTokensByUserID(userID uint) error {
	tokens, err := FindTokensByUserID(userID)
	if err != nil {
		return err
	}

	for _, token := range tokens {
		if token.Type == "application" {
			if err := ExpireToken(token.ID); err != nil {
				return err
			}
		}
	}

	return nil
}

func FindTokensByUserID(id uint) ([]models.Token, error) {
	tokens := []models.Token{}
	if err := db.Where("user_id = ?", id).Find(&tokens).Error; err != nil {
		return []models.Token{}, err
	}
	return tokens, nil
}

func ExpireAllTokens() error {
	tokens := []uint{}
	if err := db.Table("tokens").Select("id").Find(&tokens).Error; err != nil {
		return err
	}
	for _, id := range tokens {
		if err := ExpireToken(id); err != nil {
			return err
		}
	}
	return nil
}
