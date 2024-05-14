package test

import (
	"testing"

	database "github.com/nyudlts/go-medialog/database"
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
}
