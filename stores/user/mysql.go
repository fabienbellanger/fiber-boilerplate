package user

import (
	"crypto/sha512"
	"encoding/hex"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm/clause"

	"github.com/fabienbellanger/fiber-boilerplate/db"
	"github.com/fabienbellanger/fiber-boilerplate/entities"
)

// UserStore ...
type UserStore struct {
	db *db.DB
}

// New returns a new UserStore
func New(db *db.DB) UserStore {
	return UserStore{db: db}
}

// Login gets user from username and password.
func (u UserStore) Login(username, password string) (user entities.User, err error) {
	// Hash password
	// -------------
	passwordBytes := sha512.Sum512([]byte(password))
	password = hex.EncodeToString(passwordBytes[:])

	if result := u.db.Where(&entities.User{Username: username, Password: password}).First(&user); result.Error != nil {
		return user, result.Error
	}
	return user, err
}

// GetAll gets all users in database.
func (u UserStore) GetAll() ([]entities.User, error) {
	var users []entities.User

	if response := u.db.Find(&users); response.Error != nil {
		return users, response.Error
	}
	return users, nil
}

// Create adds user in database.
func (u UserStore) Create(user *entities.User) error {
	// UUID
	// ----
	user.ID = uuid.New().String()

	// Hash password
	// -------------
	passwordBytes := sha512.Sum512([]byte(user.Password))
	user.Password = hex.EncodeToString(passwordBytes[:])

	if result := u.db.Create(&user); result.Error != nil {
		return result.Error
	}
	return nil
}

// GetOne returns a user from its ID.
func (u UserStore) GetOne(id string) (user entities.User, err error) {
	if result := u.db.Find(&user, "id = ?", id); result.Error != nil {
		return user, result.Error
	}
	return user, err
}

// GetByUsername returns a user from its username.
func (u UserStore) GetByUsername(username string) (user entities.User, err error) {
	if result := u.db.Find(&user, "username = ?", username); result.Error != nil {
		return user, result.Error
	}
	return user, err
}

// Delete deletes a user from database.
func (u UserStore) Delete(id string) error {
	result := u.db.Delete(&entities.User{}, "id = ?", id)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

// Update updates user information.
func (u UserStore) Update(id string, userForm *entities.UserForm) (user entities.User, err error) {
	// Hash password
	// -------------
	hashedPassword := sha512.Sum512([]byte(userForm.Password))

	result := u.db.Model(&entities.User{}).Where("id = ?", id).Select("lastname", "firstname", "username", "password").Updates(entities.User{
		Lastname:  userForm.Lastname,
		Firstname: userForm.Firstname,
		Username:  userForm.Username,
		Password:  hex.EncodeToString(hashedPassword[:]),
	})
	if result.Error != nil {
		return user, result.Error
	}

	user, err = u.GetOne(id)
	if err != nil {
		return user, err
	}
	return user, err
}

// UpdatePassword updates user passwords.
func (u UserStore) UpdatePassword(id, currentPassword, password string) error {
	// Hash password
	// -------------
	hashedPassword := sha512.Sum512([]byte(password))

	result := u.db.Exec(`
		UPDATE users
		SET password = ?, updated_at = ?
		WHERE id = ?`,
		hex.EncodeToString(hashedPassword[:]),
		time.Now().UTC(),
		id,
	)
	if result.Error != nil {
		return result.Error
	}

	return nil
}

// GetIDFromPasswordReset update user password and delete password_resets line.
func (u UserStore) GetIDFromPasswordReset(token, password string) (string, string, error) {
	data := struct {
		ID       string
		Password string
	}{}

	result := u.db.Raw(`
			SELECT u.id AS id, u.password AS passwors
			FROM password_resets pr
				INNER JOIN users u ON u.id = pr.user_id AND u.deleted_at IS NULL
			WHERE pr.token = ?
				AND pr.expired_at >= ?`,
		token,
		time.Now().UTC()).Scan(&data)
	if result.Error != nil {
		return "", "", result.Error
	}

	return data.ID, data.Password, nil
}

// DeletePasswordReset deletes user password reset.
func (u UserStore) DeletePasswordReset(userId string) error {
	result := u.db.Where("user_id = ?", userId).Delete(&entities.PasswordResets{})
	if result.Error != nil {
		return result.Error
	}
	return nil
}

// CreateOrUpdatePasswordReset add a reset password request in database or update it if a line already exists.
func (u UserStore) CreateOrUpdatePasswordReset(passwordReset *entities.PasswordResets) error {
	result := u.db.Clauses(clause.OnConflict{
		UpdateAll: true,
	}).Create(&passwordReset)

	return result.Error
}
