package test

import (
	"encoding/json"
	"testing"

	"github.com/google/uuid"
	database "github.com/nyudlts/go-medialog/database"
	"github.com/nyudlts/go-medialog/models"
)

func TestEntries(t *testing.T) {

	t.Run("Test Unique Media ID in Resource", func(t *testing.T) {
		mediaID := 5
		resourceID := uint(30)

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

		mediaID := 1
		resourceID := uint(30)

		got, err := database.IsMediaIDUniqueInResource(mediaID, resourceID)
		if err != nil {
			t.Error(err)
		}

		want := false
		if want != got {
			t.Errorf("Wanted %v, got %v", want, got)
		}

	})

	var entryID uuid.UUID
	t.Run("Test create an entry", func(t *testing.T) {
		uid, _ := uuid.NewUUID()
		entry := models.Entry{}
		entry.ID = uid
		entry.MediaID = 789
		entry.ImagedBy = "Donald Mennerich"
		var err error
		entryID, err = database.InsertEntry(&entry)
		if err != nil {
			t.Error(err)
		}

		t.Logf("Created entry %s", entryID.String())
	})

	var entry models.Entry
	t.Run("Test get an entry", func(t *testing.T) {
		var err error
		entry, err = database.FindEntry(entryID)
		if err != nil {
			t.Error(err)
		}

		b, err := json.Marshal(entry)
		if err != nil {
			t.Error(err)
		}

		t.Logf("got accession: %s", string(b))

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

	t.Run("Test delete an entry", func(t *testing.T) {
		if err := database.DeleteEntry(entryID); err != nil {
			t.Error(err)
		}

		t.Logf("deleted entry %d", entryID)

		if _, err := database.FindEntry(entryID); err == nil {
			t.Logf("Found deleted entry %d", entryID)
		}
	})
}