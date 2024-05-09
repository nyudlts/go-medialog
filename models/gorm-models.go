package models

import (
	"time"

	"github.com/google/uuid"
)

type Repository struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	CreatedAt time.Time `json:"created_at"`
	CreatedBy int       `json:"created_by"`
	UpdatedAt time.Time `json:"updated_at"`
	UpdatedBy int       `json:"udpate_by"`
	Slug      string    `json:"slug"`
	Title     string    `json:"title"`
}

type Collection struct {
	ID             uint       `json:"id" gorm:"primaryKey" form:"id"`
	Title          string     `json:"title" form:"title"`
	CollectionCode string     `json:"collection_code" form:"collection_code"`
	PartnerCode    string     `json:"partner_code" form:"partner_code"`
	CreatedAt      time.Time  `json:"created_at"`
	UpdatedAt      time.Time  `json:"updated_at"`
	CreatedBy      int        `json:"created_by"`
	UpdatedBy      int        `json:"modified_by"`
	RepositoryID   int        `json:"repository_id" form:"repository_id"`
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
	ID                    uuid.UUID  `json:"id" gorm:"primaryKey" form:"id"`
	CreatedAt             time.Time  `json:"created_at"`
	UpdatedAt             time.Time  `json:"updated_at"`
	CreatedBy             int        `json:"created_by"`
	UpdatedBy             int        `json:"updated_by"`
	MediaID               int        `json:"media_id" form:"media_id"`
	Mediatype             string     `json:"mediatype" form:"mediatype"`
	Manufacturer          string     `json:"manufacturer" form:"manufacturer"`
	ManufacturerSerial    string     `json:"manufacturer_serial" form:"manufacturer_serial"`
	LabelText             string     `json:"label_text" form:"label_text"`
	MediaNote             string     `json:"media_note" form:"media_note"`
	HDDInterface          string     `json:"hdd_interface" form:"hdd_interface"`
	ImagingSuccess        string     `json:"imaging_success" form:"imaging_success"`
	ImageFilename         string     `json:"image_filename" form:"image_filename"`
	Interface             string     `json:"interface" form:"interface"`
	ImagingSoftware       string     `json:"imaging_software" form:"imaging_software"`
	InterpretationSuccess string     `json:"interpretation_success" form:"interpretation_success"`
	ImagedBy              string     `json:"imaged_by" form:"imaged_by"`
	ImagingNote           string     `json:"imaging_note" form:"imaging_note"`
	ImageFormat           string     `json:"image_format" form:"image_format"`
	BoxNumber             int        `json:"box_number" form:"box_number"`
	OriginalID            string     `json:"original_id" form:"original_id"`
	DispositionNote       string     `json:"disposition_note" form:"disposition_note"`
	StockUnit             string     `json:"stock_unit" form:"stock_unit"`
	StockSizeNum          float32    `json:"stock_size_num" form:"stock_size_num"`
	CollectionID          int        `json:"collection_id" form:"collection_id"`
	Collection            Collection `json:"collection"`
	AccessionID           int        `json:"accession_id" form:"accession_id"`
	Accession             Accession  `json:"accession"`
	RepositoryID          int        `json:"repository_id" form:"repository_id"`
	Repository            Repository `json:"repository"`
	IsRefreshed           bool       `json:"is_refreshed" form:"is_refreshed"`
	IsTransferred         bool       `json:"is_transferred"`
	ContentType           string     `json:"content_type" form:"content_type"`
	Structure             string     `json:"structure"`
}

func (e *Entry) UpdateEntry(updatedEntry Entry) {
	e.Mediatype = updatedEntry.Mediatype
	e.DispositionNote = updatedEntry.DispositionNote
	e.BoxNumber = updatedEntry.BoxNumber
	e.StockSizeNum = updatedEntry.StockSizeNum
	e.StockUnit = updatedEntry.StockUnit
	e.ContentType = updatedEntry.ContentType
	e.LabelText = updatedEntry.LabelText
	e.OriginalID = updatedEntry.OriginalID
	e.Manufacturer = updatedEntry.Manufacturer
	e.ManufacturerSerial = updatedEntry.ManufacturerSerial
	e.MediaNote = updatedEntry.MediaNote
	e.DispositionNote = updatedEntry.DispositionNote
	e.IsRefreshed = updatedEntry.IsRefreshed
	e.UpdatedAt = time.Now()
	e.ImageFilename = updatedEntry.ImageFilename
	e.Interface = updatedEntry.Interface
	e.HDDInterface = updatedEntry.HDDInterface
	e.ImagingSuccess = updatedEntry.ImagingSuccess
	e.InterpretationSuccess = updatedEntry.InterpretationSuccess
	e.ImagedBy = updatedEntry.ImagedBy
	e.ImagingNote = updatedEntry.ImagingNote
	e.ImagingSoftware = updatedEntry.ImagingSoftware
	e.ImageFormat = updatedEntry.ImageFormat

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
