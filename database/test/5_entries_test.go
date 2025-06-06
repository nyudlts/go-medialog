package test

import (
	"encoding/json"
	"testing"

	"github.com/google/uuid"
	"github.com/nyudlts/go-medialog/database"
	"github.com/nyudlts/go-medialog/models"
)

var entryID uuid.UUID

func TestEntries(t *testing.T) {

	t.Run("Test create an entry", func(t *testing.T) {
		uid, _ := uuid.NewUUID()
		entry := models.Entry{}
		entry.ID = uid
		entry.MediaID = 789
		entry.ResourceID = resourceID
		entry.RepositoryID = repositoryID
		entry.AccessionID = accessionID
		entry.ImagedBy = "Donald Mennerich"
		entry.Mediatype = "stuff"
		entry.CreatedBy = int(userID)
		entry.UpdatedBy = int(userID)
		entry.StockSizeNum = 1.2
		entry.StockUnit = "MB"
		entry.LabelText = "Rusty Buckles"
		var err error
		if err = database.InsertEntry(&entry); err != nil {
			t.Error(err)
		}

		t.Logf("Created entry %s", entry.ID.String())
		entryID = entry.ID
	})

	t.Run("Test invalid ID in entry", func(t *testing.T) {

		entry := models.Entry{}

		if err := entry.ValidateEntry(); err == nil {
			t.Error(err)
		}
	})

	t.Run("Test invalid mediaID in entry", func(t *testing.T) {
		uid, _ := uuid.NewUUID()
		entry := models.Entry{ID: uid}

		if err := entry.ValidateEntry(); err == nil {
			t.Error(err)
		}
	})

	t.Run("Test invalid media type in entry", func(t *testing.T) {
		uid, _ := uuid.NewUUID()
		entry := models.Entry{ID: uid, MediaID: 776}

		if err := entry.ValidateEntry(); err == nil {
			t.Error(err)
		}
	})

	t.Run("Test invalid stock size in entry", func(t *testing.T) {
		uid, _ := uuid.NewUUID()
		entry := models.Entry{ID: uid, MediaID: 765, Mediatype: "thing"}
		if err := entry.ValidateEntry(); err == nil {
			t.Error(err)
		}
	})

	t.Run("Test invalid stock unit in entry", func(t *testing.T) {
		uid, _ := uuid.NewUUID()
		entry := models.Entry{ID: uid, MediaID: 765, Mediatype: "thing", StockSizeNum: 4.7}
		if err := entry.ValidateEntry(); err == nil {
			t.Error(err)
		}
	})

	t.Run("Test valid entry", func(t *testing.T) {
		uid, _ := uuid.NewUUID()
		entry := models.Entry{ID: uid, MediaID: 765, Mediatype: "thing", StockSizeNum: 4.7, StockUnit: "GB"}
		if err := entry.ValidateEntry(); err != nil {
			t.Error(err)
		}
	})

	var entry models.Entry
	t.Run("Test get an entry", func(t *testing.T) {
		var err error
		entry, err = database.FindEntry(entryID)
		if err != nil {
			t.Error(err)
		}

		if _, err = json.Marshal(entry); err != nil {
			t.Error(err)
		}

	})

	t.Run("Test Unique Media ID in Resource", func(t *testing.T) {
		mediaID := uint(5)

		got, err := database.IsMediaIDUniqueInResource(mediaID, resourceID)
		if err != nil {
			t.Error(err)
		}

		want := true
		if want != got {
			t.Errorf("Wanted %v, got %v", want, got)
		}
	})

	t.Run("Test Non-Unique Media ID in Resource", func(t *testing.T) {

		mediaID := uint(789)

		got, err := database.IsMediaIDUniqueInResource(mediaID, resourceID)
		if err != nil {
			t.Error(err)
		}

		want := false
		if want != got {
			t.Errorf("Wanted %v, got %v", want, got)
		}

	})

	t.Run("test update an entry", func(t *testing.T) {
		entry.DispositionNote = "To Be Deleted"
		if err := database.UpdateEntry(&entry); err != nil {
			t.Error(err)
		}

		entry2, err := database.FindEntry(entryID)
		if err != nil {
			t.Error(err)
		}

		if entry.DispositionNote != entry2.DispositionNote {
			t.Errorf("Wanted: %s, Got %s", entry.DispositionNote, entry2.DispositionNote)
		}
	})
}
