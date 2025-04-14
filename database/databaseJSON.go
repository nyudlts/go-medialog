package database

import (
	"encoding/json"
	"log"

	"github.com/google/uuid"
	"github.com/nyudlts/go-medialog/models"
)

func GetJSONs() ([]models.EntryJSON, error) {
	entryJSONs := []models.EntryJSON{}
	if err := db.Find(&entryJSONs).Error; err != nil {
		return entryJSONs, err
	}
	return entryJSONs, nil
}

func GetJSON(id uint) (models.EntryJSON, error) {
	j := models.EntryJSON{}
	if err := db.Find(&j, id).Error; err != nil {
		return j, err
	}
	return j, nil
}

func CreateJSON() error {
	entryIDs, err := getEntryIDs()
	if err != nil {
		return err
	}

	for _, entryID := range entryIDs {
		log.Printf("Creating JSON for: %s", entryID.String())
		entry, err := FindEntry(entryID)
		if err != nil {
			return err
		}

		if err := InsertEntryJSON(entry); err != nil {
			return err
		}
	}
	log.Println("Success")
	return nil
}

func InsertEntryJSON(entry models.Entry) error {
	entryJson := models.EntryJSON{}
	entryJson.EntryID = entry.ID
	em := entry.Minimal()
	ebBytes, err := json.Marshal(em)
	if err != nil {
		return err
	}
	entryJson.JSON = string(ebBytes)
	if err := db.Create(&entryJson).Error; err != nil {
		return err
	}
	return nil
}

func UpdateEntryJSON(ej models.EntryJSON) error {

	if err := db.Save(&ej).Error; err != nil {
		return err
	}
	return nil
}

func DeleteEntryJSON(id uint) error {
	if err := db.Delete(&models.EntryJSON{}, id).Error; err != nil {
		return err
	}
	return nil
}

func FindEntryJSONByEntryID(u uuid.UUID) (models.EntryJSON, error) {
	var ej models.EntryJSON
	if err := db.First(&ej).Where("entry_id = ?", u).Error; err != nil {
		return ej, err
	}
	return ej, nil
}
