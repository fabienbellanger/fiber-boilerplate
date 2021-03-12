package repositories

import (
	"crypto/sha512"
	"encoding/hex"

	"github.com/google/uuid"

	"github.com/fabienbellanger/fiber-boilerplate/db"
	"github.com/fabienbellanger/fiber-boilerplate/models"
)

// Login gets user from username and password.
func Login(db *db.DB, username, password string) (user models.User, err error) {
	// Hash password
	// -------------
	passwordBytes := sha512.Sum512([]byte(password))
	password = hex.EncodeToString(passwordBytes[:])

	if result := db.Where(&models.User{Username: username, Password: password}).First(&user); result.Error != nil {
		return user, result.Error
	}
	return user, err
}

// ListAllUsers gets all users in database.
func ListAllUsers(db *db.DB) ([]models.User, error) {
	var users []models.User

	if response := db.Find(&users); response.Error != nil {
		return users, response.Error
	}
	return users, nil
}

// CreateUser adds user in database.
func CreateUser(db *db.DB, user *models.User) error {
	// UUID
	// ----
	user.ID = uuid.New().String()

	// Hash password
	// -------------
	passwordBytes := sha512.Sum512([]byte(user.Password))
	user.Password = hex.EncodeToString(passwordBytes[:])

	if result := db.Create(&user); result.Error != nil {
		return result.Error
	}
	return nil
}

// GetUser returns a user from its ID.
func GetUser(db *db.DB, id string) (user models.User, err error) {
	if result := db.Find(&user, "id = ?", id); result.Error != nil {
		return user, result.Error
	}
	return user, err
}

// DeleteUser deletes a user from database.
func DeleteUser(db *db.DB, id string) error {
	result := db.Delete(&models.User{}, "id = ?", id)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

// UpdateUser updates user information.
func UpdateUser(db *db.DB, id string, userForm *models.UserForm) (user models.User, err error) {
	// Hash password
	// -------------
	hashedPassword := sha512.Sum512([]byte(userForm.Password))

	result := db.Model(&models.User{}).Where("id = ?", id).Select("lastname", "firstname", "username", "password").Updates(models.User{
		Lastname:  userForm.Lastname,
		Firstname: userForm.Firstname,
		Username:  userForm.Username,
		Password:  hex.EncodeToString(hashedPassword[:]),
	})
	if result.Error != nil {
		return user, result.Error
	}

	user, err = GetUser(db, id)
	if err != nil {
		return user, err
	}
	return user, err
}
