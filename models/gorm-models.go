package models

import (
	"gorm.io/gorm"
)

type Accession struct {
	gorm.Model
	AccessionNum   string `json:"accession_num"`
	CollectionID   int    `json:"collection_id"`
	AccessionNote  string `json:"accession_note"`
	AccessionState string `json:"accession_state"`
	CreatedBy      int    `json:"created_by"`
	UpdatedBy      int    `json:"updated_by"`
}

type Collection struct {
	gorm.Model
	Title          string `json:"title"`
	CollectionCode string `json:"collection_code"`
	PartnerCode    string `json:"partner_code"`
	CreatedBy      int    `json:"created_by"`
	UpdatedBy      int    `json:"modified_by"`
}

type User struct {
	gorm.Model
	Email             string `json:"email"`
	Salt              string `json:"salt"`
	EncryptedPassword string `json:"encrypted_password"`
	SignInCount       int    `json:"sign_in_count"`
	CreatedBy         int    `json:"created_by"`
	UpdatedBy         int    `json:"modified_by"`
	IsActive          bool   `json:"is_active"`
	IsAdmin           bool   `json:"admin"`
}
