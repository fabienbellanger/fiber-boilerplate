package repositories

import (
	"crypto/sha512"
	"encoding/hex"
	"errors"

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
func GetUser(db *db.DB, id uint) (user models.User, err error) {
	if result := db.Find(&user, id); result.Error != nil {
		return user, result.Error
	}
	return user, err
}

// DeleteUser deletes a user from database.
func DeleteUser(db *db.DB, id uint) error {
	result := db.Delete(&models.User{}, id)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("No user")
	}
	return nil
}
