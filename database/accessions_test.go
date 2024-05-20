package database

import (
	"encoding/json"
	"flag"
	"testing"
	"time"

	"github.com/nyudlts/go-medialog/config"
	database "github.com/nyudlts/go-medialog/database"
	"github.com/nyudlts/go-medialog/models"
)

func TestAccessions(t *testing.T) {
	flag.Parse()
	env, _ = config.GetEnvironment(configuration, environment)
	database.ConnectDatabase(env.DatabaseLocation)

	var accessionID uint
	t.Run("Test create an accession", func(t *testing.T) {
		accession := models.Accession{}
		cid := 30
		accession.CollectionID = cid
		accession.AccessionNum = "yyy"
		accession.CreatedBy = 56
		accession.CreatedAt = time.Now()
		accession.CreatedBy = 56
		accession.CreatedAt = time.Now()
		var err error
		accessionID, err = database.InsertAccession(&accession)
		if err != nil {
			t.Error(err)
		}

		if cid != accession.CollectionID {
			t.Errorf("Wanted: %d, got: %d", cid, accession.CollectionID)
		}

	})

	var accession models.Accession
	t.Run("Test get an accession", func(t *testing.T) {
		var err error
		accession, err = database.FindAccession(accessionID)
		if err != nil {
			t.Error(err)
		}

		b, err := json.Marshal(accession)
		if err != nil {
			t.Error(err)
		}

		t.Logf("got accession: %s", string(b))
	})

	t.Run("Test update an accession", func(t *testing.T) {
		accession.AccessionNote = "test"
		accession.UpdatedAt = time.Now()

		if err := database.UpdateAccession(&accession); err != nil {
			t.Error(err)
		}

		accession2, err := database.FindAccession(accessionID)
		if err != nil {
			t.Error(err)
		}

		if accession2.AccessionNote != accession.AccessionNote {
			t.Errorf("Wanted: %s, Got: %s", accession.AccessionNote, accession2.AccessionNote)
		}

		t.Logf("Updated accession %d", accession.ID)
	})

	t.Run("Test delete an accession", func(t *testing.T) {
		if err := database.DeleteAccession(accessionID); err != nil {
			t.Error(err)
		}

		t.Logf("deleted accessions %d", accessionID)

		if _, err := database.FindAccession(accessionID); err == nil {
			t.Logf("Found deleted accession %d", accessionID)
		}
	})

}
