package models

import (
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/pgtype"
)

type AccessionPG struct {
	ID             int       `json:"id"`
	AccessionNum   string    `json:"accession_num"`
	CollectionID   int       `json:"collection_id"`
	AccessionNote  string    `json:"accession_note"`
	AccessionState string    `json:"accession_state"`
	CreatedBy      int       `json:"created_by"`
	CreatedAt      time.Time `json:"created_at"`
	ModifiedBy     int       `json:"modified_by"`
	ModifiedAt     time.Time `json:"modified_at"`
}

func (a *AccessionPG) ToGormModel() Accession {
	accession := Accession{}
	accession.ID = uint(a.ID)
	accession.AccessionNum = a.AccessionNum
	accession.AccessionNote = a.AccessionNote
	accession.AccessionState = a.AccessionState
	accession.ResourceID = uint(a.CollectionID)
	accession.CreatedAt = a.CreatedAt
	accession.CreatedBy = a.CreatedBy
	accession.UpdatedAt = a.ModifiedAt
	accession.UpdatedBy = a.ModifiedBy

	return accession
}

type CollectionPG struct {
	ID             int       `json:"id"`
	Title          string    `json:"title"`
	CollectionCode string    `json:"collection_code"`
	PartnerCode    string    `json:"partner_code"`
	CreatedBy      int       `json:"created_by"`
	CreatedAt      time.Time `json:"created_at"`
	ModifiedBy     int       `json:"modified_by"`
	ModifiedAt     time.Time `json:"modified_at"`
}

func (c *CollectionPG) ToGormModel() Resource {
	resource := Resource{}
	resource.ID = uint(c.ID)
	resource.Title = c.Title
	resource.CollectionCode = c.CollectionCode
	resource.PartnerCode = c.PartnerCode
	resource.CreatedBy = c.CreatedBy
	resource.CreatedAt = c.CreatedAt
	resource.UpdatedBy = c.ModifiedBy
	resource.UpdatedAt = c.ModifiedAt
	return resource
}

type UserPG struct {
	ID                  int         `json:"id"`
	Email               string      `json:"email"`
	EncryptedPassword   string      `json:"encrypted_password"`
	ResetPasswordToken  string      `json:"reset_password_token"`
	ResetPasswordSentAt time.Time   `json:"reset_password_sent_at"`
	RememberCreatedAT   time.Time   `json:"remember_created_at"`
	SignInCount         int         `json:"sign_in_count"`
	CurrentSignInAt     time.Time   `json:"current_sign_in_at"`
	CurrentSingInIP     pgtype.Inet `json:"current_sign_in_ip"`
	LastSignInAt        time.Time   `json:"last_sign_in_at"`
	LastSignInIP        pgtype.Inet `json:"last_sign_in_ip"`
	CreatedBy           int         `json:"created_by"`
	CreatedAt           time.Time   `json:"created_at"`
	ModifiedBy          int         `json:"modified_by"`
	ModifiedAt          time.Time   `json:"modified_at"`
	IsActive            bool        `json:"is_active"`
	Admin               bool        `json:"admin"`
	DeletedAt           time.Time   `json:"deleted_at"`
}

func (u *UserPG) ToGormModel() User {
	user := User{}
	user.ID = uint(u.ID)
	user.Email = u.Email
	user.SignInCount = u.SignInCount
	user.IsActive = u.IsActive
	user.IsAdmin = u.Admin
	user.CreatedAt = u.CreatedAt
	user.CreatedBy = u.CreatedBy
	user.UpdatedAt = u.ModifiedAt
	user.UpdatedBy = u.ModifiedBy
	return user
}

type Mlog_EntryPG struct {
	ID                    uuid.UUID `json:"id"`
	AccessionNum          string    `json:"accession_num"` //remove
	MediaID               int       `json:"media_id"`
	Mediatype             string    `json:"mediatype"`
	Manufacturer          string    `json:"manufacturer"`
	ManufacturerSerial    string    `json:"manufacturer_serial"`
	LabelText             string    `json:"label_text"`
	MediaNote             string    `json:"media_note"`
	PhotoURL              string    `json:"photo_url"`
	HDD_Interface         string    `json:"hdd_interface"`
	Imaging_success       string    `json:"imaging_success"`
	ImageFilename         string    `json:"image_filename"`
	Interface             string    `json:"interface"`
	ImagingSoftware       string    `json:"imaging_software"`
	InterpretationSuccess string    `json:"interpretation_success"`
	ImagedBy              string    `json:"imaged_by"`
	ImagingNote           string    `json:"imaging_note"`
	ImageFormat           string    `json:"image_format"`
	EncodingScheme        string    `json:"encoding_scheme"`
	PartitionTableFormat  string    `json:"partition_table_format"` //remove
	NumberOfPartitions    int       `json:"number of partitions"`   //remove
	FileSystem            string    `json:"filesystem"`
	HasDFXML              bool      `json:"has_dfxml"`        //remove
	Has_FTK_CSV           bool      `json:"has_ftk_csv"`      //remove
	ImageSizeBytes        int64     `json:"image_size_bytes"` //remove
	MD5Checksum           string    `json:"md5_checksum"`     //remove
	SHA1Checksum          string    `json:"sha1_checksum"`    //remove
	DateImaged            time.Time `json:"date_imaged"`
	DateFTKLoaded         time.Time `json:"date_ftk_loaded"`         //remove
	DateMetadataExtracted time.Time `json:"date_metadata_extracted"` //remove
	DatePhotographed      time.Time `json:"date_photographed"`       //remove
	DateQC                time.Time `json:"date_qc"`                 //remove
	DatePackage           time.Time `json:"date_packaged"`           //remove
	DateTransferred       time.Time `json:"date_transferred"`        //remove
	CreatedAt             time.Time `json:"created_at"`
	UpdatedAt             time.Time `json:"updated_at"`
	NumImageSegments      int       `json:"num_of_image_segments"` //remove
	RefID                 string    `json:"ref_id"`                //remove
	HasMactimeCSV         bool      `json:"has_mactime_csv"`       //remove
	BoxNumber             int       `json:"box_number"`
	StockSize             string    `json:"stock_size"`
	SIPID                 int       `json:"sip_id"` //remove
	OriginalID            string    `json:"original_id"`
	DispositionNote       string    `json:"disposition_note"`
	StockUnit             string    `json:"stock_unit"`
	StockSizeNum          float32   `json:"stock_size_num"`
	CreatedBy             int       `json:"created_by"`
	ModifiedBy            int       `json:"modified_by"`
	CollectionID          int       `json:"collection_id"`
	AccessionID           int       `json:"accession_id"`
	IsTransferred         bool      `json:"is_transferred"`
	IsRefreshed           bool      `json:"is_refreshed"`
	SessionCount          int       `json:"session_count"`
	ContentType           string    `json:"content_type"`
	Structure             string    `json:"structure"`
	FileSystems           string    `json:"file_systems"`
	BoxText               string    `json:"box_text"`
}

func (mlog *Mlog_EntryPG) ToGormModel() Entry {
	e := Entry{}
	e.ID = mlog.ID
	e.CreatedAt = mlog.CreatedAt
	e.UpdatedAt = mlog.UpdatedAt
	e.CreatedBy = mlog.CreatedBy
	e.UpdatedBy = mlog.ModifiedBy
	e.MediaID = uint(mlog.MediaID)
	e.Mediatype = mlog.Mediatype
	e.Manufacturer = mlog.Manufacturer
	e.ManufacturerSerial = mlog.ManufacturerSerial
	e.LabelText = mlog.LabelText
	e.MediaNote = mlog.MediaNote
	e.HDDInterface = mlog.HDD_Interface
	e.ImagingSuccess = mlog.Imaging_success
	e.ImageFilename = mlog.ImageFilename
	e.Interface = mlog.Interface
	e.ImagingSoftware = mlog.ImagingSoftware
	e.InterpretationSuccess = mlog.InterpretationSuccess
	e.ImagedBy = mlog.ImagedBy
	e.ImagingNote = mlog.ImagingNote
	e.ImageFormat = mlog.ImageFormat
	e.BoxNumber = mlog.BoxText
	e.OriginalID = mlog.OriginalID
	e.DispositionNote = mlog.DispositionNote
	e.StockUnit = mlog.StockUnit
	e.StockSizeNum = mlog.StockSizeNum
	e.ResourceID = uint(mlog.CollectionID)
	e.AccessionID = uint(mlog.AccessionID)
	e.IsRefreshed = mlog.IsRefreshed
	e.IsTransferred = mlog.IsTransferred
	e.ContentType = mlog.ContentType
	e.Structure = mlog.Structure
	return e
}
