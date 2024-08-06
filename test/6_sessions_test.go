package test

import (
	"testing"

	"github.com/nyudlts/go-medialog/database"
)

func TestSessions(t *testing.T) {
	t.Run("Test get sessions", func(t *testing.T) {
		sessions, err := database.GetSessions()
		if err != nil {
			t.Error(err)
		}

		for _, session := range sessions {
			t.Logf("%s", session.ID)
		}
	})
}
