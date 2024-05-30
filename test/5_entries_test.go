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
		entry.CollectionID = int(resourceID)
		entry.RepositoryID = int(repositoryID)
		entry.AccessionID = int(accessionID)
		entry.ImagedBy = "Donald Mennerich"
		entry.CreatedBy = int(userID)
		entry.UpdatedBy = int(userID)
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

	t.Run("Test Unique Media ID in Resource", func(t *testing.T) {
		mediaID := 5
		resourceID := resourceID

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

		mediaID := 789
		resourceID := resourceID

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
