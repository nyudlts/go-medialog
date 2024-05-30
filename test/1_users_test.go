package test

import (
	"crypto/sha512"
	"encoding/hex"
	"encoding/json"
	"testing"
	"time"

	"github.com/nyudlts/go-medialog/controllers"
	"github.com/nyudlts/go-medialog/database"
	"github.com/nyudlts/go-medialog/models"
)

var userID uint

func TestUsers(t *testing.T) {

	var email = "test@nyu.edu"
	var password = "parallel"

	t.Run("Test create a user", func(t *testing.T) {
		user := models.User{}
		user.Email = email
		user.CreatedAt = time.Now()
		user.CreatedBy = 56
		user.Salt = controllers.GenerateStringRunes(16)
		hash := sha512.Sum512([]byte(password + user.Salt))
		user.EncryptedPassword = hex.EncodeToString(hash[:])
		user.IsActive = true
		user.IsAdmin = false

		var err error
		userID, err = database.InsertUser(&user)
		if err != nil {
			t.Error(err)
		}

		t.Logf("Created user %d", userID)
	})

	var user models.User
	t.Run("Test get a user", func(t *testing.T) {
		var err error
		user, err = database.FindUser(userID)
		if err != nil {
			t.Error(err)
		}

		b, err := json.Marshal(user)
		if err != nil {
			t.Error(err)
		}

		t.Logf("got user %s", string(b))
	})

	t.Run("Test authenticate a user", func(t *testing.T) {

		hash := sha512.Sum512([]byte(password + user.Salt))
		userSHA512 := hex.EncodeToString(hash[:])

		if userSHA512 != user.EncryptedPassword {
			t.Errorf("Got %s, wanted %s", userSHA512, user.EncryptedPassword)
		}

		t.Logf("Authenticated user %d", userID)
	})

	t.Run("Test update a user", func(t *testing.T) {
		user.IsActive = false
		if err := database.UpdateUser(&user); err != nil {
			t.Error(err)
		}

		user2, err := database.FindUser(userID)
		if err != nil {
			t.Error(err)
		}

		if user2.IsActive {
			t.Errorf("Wanted false, Got True")
		}
	})
}
