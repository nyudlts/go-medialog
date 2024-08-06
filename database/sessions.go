package database

import "time"

type session struct {
	ID        string    `json:"id"`
	Data      string    `json:"data"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	ExpiredAt time.Time `json:"expires_at"`
}

func GetSessions() ([]session, error) {
	seshes := []session{}
	if err := db.Find(&seshes).Error; err != nil {
		return seshes, err
	}
	return seshes, nil
}
