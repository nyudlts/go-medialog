package database

import (
	"log"

	"github.com/nyudlts/go-medialog/models"
)

func InsertSesssion(session models.Session) error {
	if err := db.Create(&session).Error; err != nil {
		return err
	}
	return nil
}

func FindSessionByKey(sessionKey string) (models.Session, error) {
	session := models.Session{}
	if err := db.Where("session_key = ?", sessionKey).First(&session).Error; err != nil {
		return session, err
	}
	return session, nil
}

func DropSession(sessionKey string) error {
	session, err := FindSessionByKey(sessionKey)
	if err != nil {
		return err
	}

	log.Println(session)

	if err := db.Delete(&session).Error; err != nil {
		return err
	}

	return nil
}
