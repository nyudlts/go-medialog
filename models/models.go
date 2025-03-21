package models

import (
	"fmt"
	"strings"
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
	CreatedBy             int        `json:"created_by"` //this should be converted to a uint
	UpdatedBy             int        `json:"updated_by"` //this should be converted to a uint
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
	Status                string     `json:"status" form:"status"`
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
	Location              string     `json:"location" form:"location"`
}

var CSVHeader = []string{"id", "media_id", "mediatype", "content_type", "label_text", "is_refreshed", "imaging_success", "repository", "resource", "accession", "storage_location"}

func (e Entry) ToCSV() []string {

	labelText := strings.ReplaceAll(e.LabelText, "\n", " ")

	csv := []string{
		e.ID.String(),
		fmt.Sprintf("%d", e.MediaID),
		e.Mediatype,
		e.ContentType,
		labelText,
		boolToString(e.IsRefreshed),
		e.ImagingSuccess,
		e.Repository.Slug,
		e.Resource.CollectionCode,
		e.Accession.AccessionNum,
		e.Location,
	}
	return csv
}

func (e Entry) ToCSVEntryResult() CSVEntryResult {
	return CSVEntryResult{
		ID:              e.ID,
		MediaID:         e.MediaID,
		MediaType:       e.Mediatype,
		ContentType:     e.ContentType,
		LabelText:       e.LabelText,
		IsRefreshed:     e.IsRefreshed,
		ImagingSuccess:  e.ImagingSuccess,
		RepositoryID:    e.RepositoryID,
		ResourceID:      e.ResourceID,
		AccessionID:     e.AccessionID,
		StorageLocation: e.Location,
	}
}

func boolToString(b bool) string {
	if b {
		return "TRUE"
	}
	return "FALSE"
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
	e.Location = updatedEntry.Location
	e.Status = updatedEntry.Status
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
	ID                uint      `json:"id" gorm:"primaryKey" form:"id"`
	Email             string    `json:"email" form:"email"`
	Salt              string    `json:"salt"`
	EncryptedPassword string    `json:"encrypted_password"`
	SignInCount       int       `json:"sign_in_count"`
	CreatedAt         time.Time `json:"created_at"`
	UpdatedAt         time.Time `json:"updated_at"`
	CreatedBy         int       `json:"created_by"`
	UpdatedBy         int       `json:"updated_by"`
	IsActive          bool      `json:"is_active"`
	IsAdmin           bool      `json:"is_admin"`
	CurrentIPAddress  string    `json:"current_ip_address"`
	PreviousIPAddress string    `json:"previous_ip_address"`
	FirstName         string    `json:"first_name" form:"first_name"`
	LastName          string    `json:"last_name" form:"last_name"`
	CanAccessAPI      bool      `json:"can_access_api" form:"can_access_api"`
}

type Token struct {
	Token   string    `json:"token"`
	ID      uint      `json:"id" gorm:"primaryKey"`
	IsValid bool      `json:"is_valid"`
	Expires time.Time `json:"expires"`
	UserID  uint      `json:"user_id"`
	User    User      `json:"user"`
	Type    string    `json:"type"`
}

// config functions
type Environment struct {
	LogLocation    string         `yaml:"log"`
	DatabaseConfig DatabaseConfig `yaml:"database"`
	TestCreds      TestCreds      `yaml:"test_creds"`
	AdminEmail     string         `yaml:"admin_email"`
}

type DatabaseConfig struct {
	Username     string `yaml:"username"`
	Password     string `yaml:"password"`
	URL          string `yaml:"url"`
	Port         string `yaml:"port"`
	DatabaseName string `yaml:"database_name"`
}

type TestCreds struct {
	Username string `yaml:"username"`
	Password string `yaml:"password"`
}

type MedialogInfo struct {
	Version       string
	GinVersion    string
	GolangVersion string
	APIVersion    string
}

func (mli MedialogInfo) String() string {
	return fmt.Sprintf("Medialog Version: %s, API Version: %s, Golang Version: %s, Gin Version: %s", mli.Version, mli.APIVersion, mli.GolangVersion, mli.GinVersion)
}

type CSVEntryResult struct {
	ID              uuid.UUID
	MediaID         uint
	MediaType       string
	ContentType     string
	LabelText       string
	IsRefreshed     bool
	ImagingSuccess  string
	RepositoryID    uint
	ResourceID      uint
	AccessionID     uint
	StorageLocation string
}

func (er CSVEntryResult) ToCSV() []string {
	return []string{
		er.ID.String(),
		fmt.Sprintf("%d", er.MediaID),
		er.MediaType,
		er.ContentType,
		er.LabelText,
		boolToString(er.IsRefreshed),
		er.ImagingSuccess,
		fmt.Sprintf("%d", er.RepositoryID),
		fmt.Sprintf("%d", er.ResourceID),
		fmt.Sprintf("%d", er.AccessionID),
		er.StorageLocation,
	}
}
