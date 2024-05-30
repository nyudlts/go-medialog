//go:build exclude

package test

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/nyudlts/go-medialog/models"
)

func TestResources(t *testing.T) {

	var resourceID uint
	t.Run("Test Create A Resource", func(t *testing.T) {
		resource := models.Collection{}
		resource.PartnerCode = "fales"
		resource.CollectionCode = "mss.1000"
		resource.Title = "Test Resource"
		resource.CreatedBy = 56
		resource.CreatedAt = time.Now()
		resource.RepositoryID = 3

		var err error
		resourceID, err = InsertResource(&resource)
		if err != nil {
			t.Error(err)
		}
		t.Logf("Created resource %d", resourceID)
	})

	var resource models.Collection
	t.Run("Test Get A Resource", func(t *testing.T) {
		var err error
		resource, err = FindResource(resourceID)
		if err != nil {
			t.Error(err)
		}

		b, err := json.Marshal(resource)
		if err != nil {
			t.Error(err)
		}
		t.Log("returned resource: " + string(b))
	})

	t.Run("Test update a resource", func(t *testing.T) {
		resource.Title = "updated title"
		if err := UpdateResource(&resource); err != nil {
			t.Error(err)
		}

		t.Logf("Updated resource %d", resource.ID)

		resource2, err := FindResource(resource.ID)
		if err != nil {
			t.Error(err)
		}

		if resource2.Title != resource.Title {
			t.Errorf("Got: %s, Wanted: %s", resource2.Title, resource.Title)
		}
	})

	t.Run("Test Delete A Resource", func(t *testing.T) {
		if err := DeleteResource(resource.ID); err != nil {
			t.Error(err)
		}

		if _, err := FindResource(resource.ID); err == nil {
			t.Logf("Found deleted resource %d", resource.ID)
		}
	})
}
