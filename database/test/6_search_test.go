package test

import (
	"testing"

	"github.com/nyudlts/go-medialog/database"
)

var query = "Rusty"

func TestSearch(t *testing.T) {
	t.Run("Test search for record", func(t *testing.T) {
		entries, err := database.SearchEntries(query)
		if err != nil {
			t.Error(err)
		}

		want := 1
		got := len(entries)
		if want != got {
			t.Errorf("Wanted %d entry, Got %d", want, got)
		}
	})
}
