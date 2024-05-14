package test

import (
	"encoding/json"
	"testing"
	"time"

	database "github.com/nyudlts/go-medialog/database"
	"github.com/nyudlts/go-medialog/models"
)

func TestRepositories(t *testing.T) {

	var repositoryID uint
	t.Run("Test Create A Repository", func(t *testing.T) {
		repository := models.Repository{}
		repository.Title = "Test Repository"
		repository.Slug = "test"
		repository.CreatedAt = time.Now()
		repository.CreatedBy = 56
		repository.UpdatedAt = time.Now()
		repository.UpdatedBy = 56

		var err error
		repositoryID, err = database.CreateRepository(&repository)
		if err != nil {
			t.Error(err)
		}
		t.Logf("Created repository %d", repositoryID)

	})

	var repo models.Repository
	t.Run("Test Get A Repository", func(t *testing.T) {
		var err error
		repo, err = database.FindRepository(repositoryID)
		if err != nil {
			t.Error(err)
		}

		b, err := json.Marshal(repo)
		if err != nil {
			t.Error(err)
		}
		t.Log("returned repository " + string(b))
	})

	t.Run("Test Update a Repository", func(t *testing.T) {
		repo.Slug = "tests"
		if err := database.UpdateRepository(&repo); err != nil {
			t.Error(err)
		}
		t.Logf("Repository %d Updated", repo.ID)

		repo2, err := database.FindRepository(repo.ID)
		if err != nil {
			t.Error(err)
		}

		if repo2.Slug != repo.Slug {
			t.Errorf("Got %s, wanted %s", repo2.Slug, repo.Slug)
		}
	})

	t.Run("Test Delete a Repository", func(t *testing.T) {
		if err := database.DeleteRepository(repo.ID); err != nil {
			t.Error(err)
		}
		t.Logf("Repository %d Deleted", repo.ID)

		if _, err := database.FindRepository(repo.ID); err == nil {
			t.Logf("Found deleted repository")
		}
	})
}
