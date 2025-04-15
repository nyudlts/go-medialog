package test

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/nyudlts/go-medialog/database"
	"github.com/nyudlts/go-medialog/models"
)

var resourceID uint

func TestResources(t *testing.T) {

	t.Run("Test Create A Resource", func(t *testing.T) {
		resource := models.Resource{}
		resource.PartnerCode = "fales"
		resource.CollectionCode = "mss.1000"
		resource.Title = "Test Resource"
		resource.CreatedBy = int(userID)
		resource.CreatedAt = time.Now()
		resource.UpdatedBy = int(userID)
		resource.RepositoryID = repositoryID

		var err error
		resourceID, err = database.InsertResource(&resource)
		if err != nil {
			t.Error(err)
		}
		t.Logf("Created resource %d", resourceID)
	})

	var resource models.Resource
	t.Run("Test Get A Resource", func(t *testing.T) {
		var err error
		resource, err = database.FindResource(resourceID)
		if err != nil {
			t.Error(err)
		}

		want := "mss.1000"
		got := resource.CollectionCode

		if want != got {
			t.Errorf("Wanted: %s, Got: %s", want, got)
		}

		if _, err = json.Marshal(resource); err != nil {
			t.Error(err)
		}
	})

	t.Run("Test update a resource", func(t *testing.T) {
		resource.Title = "updated title"
		if err := database.UpdateResource(&resource); err != nil {
			t.Error(err)
		}

		t.Logf("Updated resource %d", resource.ID)

		resource2, err := database.FindResource(resource.ID)
		if err != nil {
			t.Error(err)
		}

		if resource2.Title != resource.Title {
			t.Errorf("Got: %s, Wanted: %s", resource2.Title, resource.Title)
		}
	})
}
