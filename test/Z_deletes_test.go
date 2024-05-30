package test

import (
	"testing"

	"github.com/nyudlts/go-medialog/database"
)

func TestDeleteObjects(t *testing.T) {

	//delete the Entry
	t.Run("Test delete an entry", func(t *testing.T) {
		if err := database.DeleteEntry(entryID); err != nil {
			t.Error(err)
		}

		t.Logf("deleted entry %d", entryID)

		if _, err := database.FindEntry(entryID); err == nil {
			t.Logf("Found deleted entry %d", entryID)
		}
	})

	//delete the accession
	t.Run("Test delete an accession", func(t *testing.T) {
		if err := database.DeleteAccession(accessionID); err != nil {
			t.Error(err)
		}

		if _, err := database.FindAccession(accessionID); err == nil {
			t.Errorf("found deleted resource %d", resourceID)
		}

	})

	//delete the resource
	t.Run("Test delete a resource", func(t *testing.T) {
		if err := database.DeleteResource(resourceID); err != nil {
			t.Error(err)
		}

		if _, err := database.FindResource(resourceID); err == nil {
			t.Errorf("found deleted resource %d", resourceID)
		}
	})

	//delete the repository
	t.Run("test delete a repository", func(t *testing.T) {

		if err := database.DeleteRepository(repositoryID); err != nil {
			t.Error(err)
		}
		t.Logf("Repository %d Deleted", repositoryID)

		if _, err := database.FindRepository(repositoryID); err == nil {
			t.Errorf("found deleted repository")
		}
	})

	//delete the user
	t.Run("Test delete a user", func(t *testing.T) {
		if err := database.DeleteUser(userID); err != nil {
			t.Error(err)
		}

		t.Logf("deleted user %d", userID)

		if _, err := database.FindUser(userID); err == nil {
			t.Logf("Found deleted user %d", userID)
		}
	})
}
