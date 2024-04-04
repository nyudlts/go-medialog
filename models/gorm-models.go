package models

import (
	"time"

	"github.com/google/uuid"
)

type Repository struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Slug      string    `json:"slug"`
	Title     string    `json:"title"`
}

type Collection struct {
	ID             uint       `json:"id" gorm:"primaryKey"`
	Title          string     `json:"title"`
	CollectionCode string     `json:"collection_code"`
	PartnerCode    string     `json:"partner_code"`
	CreatedAt      time.Time  `json:"created_at"`
	UpdatedAt      time.Time  `json:"updated_at"`
	CreatedBy      int        `json:"created_by"`
	UpdatedBy      int        `json:"modified_by"`
	RepositoryID   int        `json:"repository_id"`
	Repository     Repository `json:"repository"`
}

type Accession struct {
	ID             uint       `json:"id" gorm:"primaryKey"`
	AccessionNum   string     `json:"accession_num"`
	AccessionNote  string     `json:"accession_note"`
	AccessionState string     `json:"accession_state"`
	CreatedAt      time.Time  `json:"created_at"`
	UpdatedAt      time.Time  `json:"updated_at"`
	CreatedBy      int        `json:"created_by"`
	UpdatedBy      int        `json:"updated_by"`
	CollectionID   int        `json:"collection_id"`
	Collection     Collection `json:"collection"`
}

type Entry struct {
	ID                    uuid.UUID  `json:"id" gorm:"primaryKey"`
	CreatedAt             time.Time  `json:"created_at"`
	UpdatedAt             time.Time  `json:"updated_at"`
	CreatedBy             int        `json:"created_by"`
	UpdatedBy             int        `json:"updated_by"`
	MediaID               string     `json:"media_id"`
	Mediatype             string     `json:"mediatype"`
	Manufacturer          string     `json:"manufacturer"`
	ManufacturerSerial    string     `json:"manufacturer_serial"`
	LabelText             string     `json:"label_text"`
	MediaNote             string     `json:"media_note"`
	HDDInterface          string     `json:"hdd_interface"`
	ImagingSuccess        string     `json:"imaging_success"`
	ImageFilename         string     `json:"image_filename"`
	Interface             string     `json:"interface"`
	ImagingSoftware       string     `json:"imaging_software"`
	InterpretationSuccess string     `json:"interpretation_success"`
	ImagedBy              string     `json:"imaged_by"`
	ImagingNote           string     `json:"imaging_note"`
	ImageFormat           string     `json:"image_format"`
	BoxNumber             int        `json:"box_number"`
	OriginalID            string     `json:"original_id"`
	DispositionNote       string     `json:"disposition_note"`
	StockUnit             string     `json:"stock_unit"`
	StockSizeNum          float32    `json:"stock_size_num"`
	CollectionID          int        `json:"collection_id"`
	Collection            Collection `json:"collection"`
	AccessionID           int        `json:"accession_id"`
	Accession             Accession  `json:"accession"`
	RepositoryID          int        `json:"repository_id"`
	IsRefreshed           bool       `json:"is_refreshed"`
	IsTransferred         bool       `json:"is_transferred"`
	ContentType           string     `json:"content_type"`
	Structure             string     `json:"structure"`
}

type User struct {
	ID                uint      `json:"id" gorm:"primaryKey"`
	Email             string    `json:"email"`
	Salt              string    `json:"salt"`
	EncryptedPassword string    `json:"encrypted_password"`
	SignInCount       int       `json:"sign_in_count"`
	CreatedAt         time.Time `json:"created_at"`
	UpdatedAt         time.Time `json:"updated_at"`
	CreatedBy         int       `json:"created_by"`
	UpdatedBy         int       `json:"modified_by"`
	IsActive          bool      `json:"is_active"`
	IsAdmin           bool      `json:"admin"`
}
