package database

import "github.com/nyudlts/go-medialog/models"

func FindUserByID(id int) (models.User, error) {
	user := models.User{}
	if err := db.Where("id = ?", id).First(&user).Error; err != nil {
		return user, err
	}
	return user, nil
}

func FindUserEmailByID(id int) (string, error) {
	user := models.User{}
	if err := db.Where("id = ?", id).First(&user).Error; err != nil {
		return "", err
	}
	return user.Email, nil
}

func UpdateUser(user *models.User) error {
	if err := db.Save(&user).Error; err != nil {
		return err
	}
	return nil
}

func FindUsers() ([]models.User, error) {
	users := []models.User{}
	if err := db.Order("is_active desc, is_admin desc, email").Find(&users).Error; err != nil {
		return users, err
	}
	return users, nil
}

func FindUser(id uint) (models.User, error) {
	user := models.User{}
	if err := db.Where("id = ?", id).First(&user).Error; err != nil {
		return models.User{}, err
	}
	return user, nil
}

func FindUserByEmail(email string) (models.User, error) {
	user := models.User{}
	if err := db.Where("email = ?", email).First(&user).Error; err != nil {
		return user, err
	}
	return user, nil
}

func InsertUser(user *models.User) (uint, error) {
	if err := db.Create(&user).Error; err != nil {
		return 0, err
	}
	return user.ID, nil
}

type EntryUser struct {
	ID    int
	Email string
}

type EntryUsers struct {
	CreateUser EntryUser
	UpdateUser EntryUser
}

func FindEntryUsers(createUserID int, modUserID int) (EntryUsers, error) {

	var createUser EntryUser
	if createUserID > 0 {
		var cUser = models.User{}
		if err := db.Where("id = ?", createUserID).First(&cUser).Error; err != nil {
			return EntryUsers{}, err
		}
		createUser = EntryUser{createUserID, cUser.Email}
	} else {
		createUser = EntryUser{0, "admin"}
	}

	var modUser EntryUser
	if modUserID > 0 {
		var mUser = models.User{}
		if err := db.Where("id = ?", modUserID).First(&mUser).Error; err != nil {
			return EntryUsers{}, err
		}
		modUser = EntryUser{modUserID, mUser.Email}
	} else {
		modUser = EntryUser{0, "admin"}
	}

	return EntryUsers{createUser, modUser}, nil
}

func DeleteUser(id uint) error {
	if err := db.Delete(models.User{}, id).Error; err != nil {
		return err
	}
	return nil
}
