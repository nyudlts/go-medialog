package test

import (
	"testing"

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
}
