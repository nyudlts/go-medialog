package models

import (
	"fmt"
	"time"

	"github.com/google/uuid"
)

type Repository struct {
	ID        uint      `json:"id" gorm:"primaryKey" form:"id"`
	CreatedAt time.Time `json:"created_at"`
	CreatedBy int       `json:"created_by"`
	UpdatedAt time.Time `json:"updated_at"`
	UpdatedBy int       `json:"updated_by"`
	Slug      string    `json:"slug" form:"slug"`
	Title     string    `json:"title" form:"title"`
}

type Resource struct {
	ID             uint       `json:"id" gorm:"primaryKey" form:"id"`
	Title          string     `json:"title" form:"title"`
	CollectionCode string     `json:"collection_code" form:"collection_code"`
	PartnerCode    string     `json:"partner_code" form:"partner_code"`
	CreatedAt      time.Time  `json:"created_at"`
	UpdatedAt      time.Time  `json:"updated_at"`
	CreatedBy      int        `json:"created_by"`
	UpdatedBy      int        `json:"updated_by"`
	RepositoryID   uint       `json:"repository_id" form:"repository_id"`
	Repository     Repository `json:"repository"`
}

type Accession struct {
	ID             uint      `json:"id" gorm:"primaryKey" form:"id"`
	AccessionNum   string    `json:"accession_num" form:"accession_num"`
	AccessionNote  string    `json:"accession_note"`  //deprecated
	AccessionState string    `json:"accession_state"` //deprecated
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
	CreatedBy      int       `json:"created_by"`
	UpdatedBy      int       `json:"updated_by"`
	ResourceID     uint      `json:"resource_id" form:"resource_id"`
	Resource       Resource  `json:"resource"`
}

type Entry struct {
	ID                    uuid.UUID  `json:"id" gorm:"primaryKey" form:"id"`
	CreatedAt             time.Time  `json:"created_at"`
	UpdatedAt             time.Time  `json:"updated_at"`
	CreatedBy             int        `json:"created_by"`
	UpdatedBy             int        `json:"updated_by"`
	MediaID               uint       `json:"media_id" form:"media_id"`
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
	BoxNumber             string     `json:"box_number" form:"box_number"`
	OriginalID            string     `json:"original_id" form:"original_id"`
	DispositionNote       string     `json:"disposition_note" form:"disposition_note"`
	StockUnit             string     `json:"stock_unit" form:"stock_unit"`
	StockSizeNum          float32    `json:"stock_size_num" form:"stock_size_num"`
	RepositoryID          uint       `json:"repository_id" form:"repository_id"`
	Repository            Repository `json:"repository"`
	ResourceID            uint       `json:"resource_id" form:"resource_id"`
	Resource              Resource   `json:"resource"`
	AccessionID           uint       `json:"accession_id" form:"accession_id"`
	Accession             Accession  `json:"accession"`
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

func (e *Entry) ValidateEntry() error {

	if _, err := uuid.Parse(e.ID.String()); err != nil {
		return fmt.Errorf("ID: `%v` is not valid", e.ID)
	}

	if e.Mediatype == "" {
		return fmt.Errorf("mediatype: `%s` is not valid", e.Mediatype)
	}

	if e.MediaID < 1 {
		return fmt.Errorf("mediaID: `%d` is not valid", e.MediaID)
	}

	if e.StockSizeNum < 1 {
		return fmt.Errorf("stock size number: `%f` is not valid", e.StockSizeNum)
	}

	if e.StockUnit == "" {
		return fmt.Errorf("stock Unit: `%s` is not valid", e.StockUnit)
	}

	return nil
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
	UpdatedBy         int       `json:"updated_by"`
	IsActive          bool      `json:"is_active"`
	IsAdmin           bool      `json:"admin"`
	isLoggedIn        bool      `json:isLoggedIn`
	CurrentIPAddress  string    `json:current_ip`
	PreviousIPAddress string    `json:previous_ip`
}
