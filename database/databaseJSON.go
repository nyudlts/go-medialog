package database

import (
	"encoding/json"
	"log"

	"github.com/google/uuid"
	"github.com/nyudlts/go-medialog/models"
)

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
	if err := db.Delete(models.EntryJSON{}, id).Error; err != nil {
		return err
	}
	return nil
}

func FindEntryJSONByEntryID(u uuid.UUID) (models.EntryJSON, error) {
	var ej models.EntryJSON
	if err := db.Table("entry_jsons").Select("id").Where("entry_id = ?", u).Find(&ej).Error; err != nil {
		return ej, err
	}
	return ej, nil
}
