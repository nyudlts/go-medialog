package database

import (
	"fmt"
	"time"
)

type Session struct {
	ID        string    `json:"id" gorm:"primaryKey"`
	Data      string    `json:"data"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	ExpiresAt time.Time `json:"expires_at"`
}

func DeleteSession(id string) error {
	fmt.Println("deleting session", id)
	if err := db.Delete(&Session{}, id).Error; err != nil {
		return err
	}
	return nil
}

func GetSession(id string) (Session, error) {
	session := Session{}
	if err := db.Where("id = ?", id).First(session).Error; err != nil {
		return session, err
	}
	return session, nil
}
