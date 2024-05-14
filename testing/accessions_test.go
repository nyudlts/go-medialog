package test

import (
	"testing"
	"time"

	database "github.com/nyudlts/go-medialog/database"
	"github.com/nyudlts/go-medialog/models"
)

func TestAccessions(t *testing.T) {

	t.Run("Get Accession", func(t *testing.T) {
		_, err := database.FindAccession(30)
		if err != nil {
			t.Error(err)
		}
	})

	t.Run("Insert Accession", func(t *testing.T) {
		accession := models.Accession{}
		cid := 30
		accession.CollectionID = cid
		accession.AccessionNum = "yyy"
		accession.CreatedBy = 56
		accession.CreatedAt = time.Now()
		accession.CreatedBy = 56
		accession.CreatedAt = time.Now()
		_, err := database.InsertAccession(&accession)
		if err != nil {
			t.Error(err)
		}

		if cid != accession.CollectionID {
			t.Error("NO")
		}

		t.Log(accession.ID, accession.CollectionID)

	})
}
