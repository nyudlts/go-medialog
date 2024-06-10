package test

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/nyudlts/go-medialog/database"
	"github.com/nyudlts/go-medialog/models"
)

var accessionID uint

func TestAccessions(t *testing.T) {
	t.Run("Test Create a Accession", func(t *testing.T) {
		accession := models.Accession{}
		accession.ResourceID = resourceID
		accession.AccessionNum = "test.acc"
		accession.CreatedBy = int(userID)
		accession.UpdatedBy = int(userID)
		var err error
		accessionID, err = database.InsertAccession(&accession)
		if err != nil {
			t.Error(err)
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
		t.Log("returned accession: " + string(b))

	})

	t.Run("Test update an accession", func(t *testing.T) {
		var err error
		accession, err = database.FindAccession(accessionID)
		if err != nil {
			t.Error(err)
		}

		accession.CreatedAt = time.Now()
		accession.UpdatedAt = time.Now()

		if err := database.UpdateAccession(&accession); err != nil {
			t.Error(err)
		}

		t.Logf("Updated accession %d", accessionID)

	})
}
