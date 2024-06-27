package database

import (
	"github.com/nyudlts/go-medialog/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var pgdb *gorm.DB

func ConnectPGSQL() error {
	var err error
	pgdb, err = gorm.Open(postgres.New(postgres.Config{
		DSN:                  "host=localhost user=medialog password=medialog dbname=medialog port=5432 sslmode=disable",
		PreferSimpleProtocol: true,
	}), &gorm.Config{})
	if err != nil {
		return err
	}
	return nil
}

func GetAccessionsPG() ([]models.AccessionPG, error) {
	var accessions = []models.AccessionPG{}
	if err := pgdb.Find(&accessions).Error; err != nil {
		return []models.AccessionPG{}, err
	}
	return accessions, nil
}

func GetEntriesPG() ([]models.Mlog_EntryPG, error) {
	var entries = []models.Mlog_EntryPG{}
	if err := pgdb.Find(&entries).Error; err != nil {
		return []models.Mlog_EntryPG{}, err
	}
	return entries, nil
}

func GetUsersPG() ([]models.UserPG, error) {
	var users = []models.UserPG{}
	if err := pgdb.Find(&users).Error; err != nil {
		return []models.UserPG{}, err
	}
	return users, nil
}

func GetCollectionsPG() ([]models.CollectionPG, error) {
	var collections = []models.CollectionPG{}
	if err := pgdb.Find(&collections).Error; err != nil {
		return []models.CollectionPG{}, err
	}
	return collections, nil
}

func CountAccessionsPG() int64 {
	var count int64
	pgdb.Model(models.AccessionPG{}).Count(&count)
	return count
}

func CountEntriesPG() int64 {
	var count int64
	pgdb.Model(models.Mlog_EntryPG{}).Count(&count)
	return count
}

func CountResourcesPG() int64 {
	var count int64
	pgdb.Model(models.CollectionPG{}).Count(&count)
	return count
}

func CountUsersPG() int64 {
	var count int64
	pgdb.Model(models.UserPG{}).Count(&count)
	return count
}
