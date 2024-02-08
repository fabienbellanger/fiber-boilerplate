package services

import (
	"bytes"
	"crypto/sha512"
	"encoding/hex"
	"errors"
	"fmt"
	"github.com/fabienbellanger/fiber-boilerplate/pkg/domain/entities"
	"github.com/fabienbellanger/fiber-boilerplate/pkg/domain/repositories"
	"github.com/fabienbellanger/fiber-boilerplate/pkg/domain/requests"
	"github.com/fabienbellanger/fiber-boilerplate/pkg/domain/responses"
	"github.com/fabienbellanger/fiber-boilerplate/utils"
	"github.com/fabienbellanger/goutils/mail"
	"github.com/google/uuid"
	"github.com/spf13/viper"
	"gorm.io/gorm"
	"html/template"
	"time"
)

type UserService interface {
	Login(req requests.UserLogin) (responses.UserLogin, *utils.HTTPError)
	Create(req requests.UserCreation) (entities.User, *utils.HTTPError)
	GetAll() ([]entities.User, *utils.HTTPError)
	GetByID(id requests.UserByID) (entities.User, *utils.HTTPError)
	Delete(id requests.UserByID) *utils.HTTPError
	Update(req requests.UserUpdate) (entities.User, *utils.HTTPError)
	UpdatePassword(req requests.UserPasswordUpdate) *utils.HTTPError
	ForgottenPassword(req requests.UserForgotPassword) (entities.PasswordResets, *utils.HTTPError)
}

type userService struct {
	userRepository repositories.UserRepository
}

// NewUser returns a new user service
func NewUser(repo repositories.UserRepository) UserService {
	return &userService{repo}
}

// Login user
func (us userService) Login(req requests.UserLogin) (responses.UserLogin, *utils.HTTPError) {
	loginErrors := utils.ValidateStruct(req)
	if loginErrors != nil {
		return responses.UserLogin{}, utils.NewHTTPError(utils.StatusBadRequest, "Invalid body", loginErrors, nil)
	}

	user, err := us.userRepository.Login(req.Username, req.Password)
	if err != nil {
		var e *utils.HTTPError
		if errors.Is(err, gorm.ErrRecordNotFound) {
			e = utils.NewHTTPError(utils.StatusUnauthorized, "Unauthorized", nil, nil)
		} else {
			e = utils.NewHTTPError(utils.StatusInternalServerError, "Internal server error", "Error during authentication", err)
		}
		return responses.UserLogin{}, e
	}

	// Create token
	token, expiresAt, err := user.GenerateJWT(
		viper.GetDuration("JWT_LIFETIME"),
		viper.GetString("JWT_ALGO"),
		viper.GetString("JWT_SECRET"))
	if err != nil {
		return responses.UserLogin{}, utils.NewHTTPError(utils.StatusInternalServerError, "Internal server error", "Error during token generation", err)
	}

	return responses.UserLogin{
		User:      user,
		Token:     token,
		ExpiresAt: expiresAt.Format("2006-01-02T15:04:05.000Z"),
	}, nil
}

// Create user
func (us userService) Create(req requests.UserCreation) (entities.User, *utils.HTTPError) {
	creationErrors := utils.ValidateStruct(req)
	if creationErrors != nil {
		return entities.User{}, utils.NewHTTPError(utils.StatusBadRequest, "Invalid body", creationErrors, nil)
	}

	newUser := entities.User{
		Lastname:  req.Lastname,
		Firstname: req.Firstname,
		Password:  req.Password,
		Username:  req.Username,
	}

	if err := us.userRepository.Create(&newUser); err != nil {
		return entities.User{}, utils.NewHTTPError(utils.StatusInternalServerError, "Database error", "Error during user creation", err)
	}

	return newUser, nil
}

// GetAll returns all users
func (us userService) GetAll() ([]entities.User, *utils.HTTPError) {
	users, err := us.userRepository.GetAll()
	if err != nil {
		return []entities.User{}, utils.NewHTTPError(utils.StatusInternalServerError, "Database error", "Error when getting all users", err)
	}

	return users, nil
}

// GetByID returns a user from its ID
func (us userService) GetByID(req requests.UserByID) (entities.User, *utils.HTTPError) {
	validateID := utils.ValidateStruct(req)
	if validateID != nil {
		return entities.User{}, utils.NewHTTPError(utils.StatusBadRequest, "Invalid parameters", validateID, nil)
	}

	user, err := us.userRepository.GetByID(req.ID)
	if err != nil {
		return entities.User{}, utils.NewHTTPError(utils.StatusInternalServerError, "Database error", "Error when getting user by id", err)
	}

	if user.ID == "" {
		return entities.User{}, utils.NewHTTPError(utils.StatusNotFound, "No user found", nil, nil)
	}

	return user, nil
}

// Delete user
func (us userService) Delete(req requests.UserByID) *utils.HTTPError {
	validateID := utils.ValidateStruct(req)
	if validateID != nil {
		return utils.NewHTTPError(utils.StatusBadRequest, "Invalid parameters", validateID, nil)
	}

	if err := us.userRepository.Delete(req.ID); err != nil {
		return utils.NewHTTPError(utils.StatusInternalServerError, "Database error", "Error when deleting the user", err)
	}

	return nil
}

// Update user
func (us userService) Update(req requests.UserUpdate) (entities.User, *utils.HTTPError) {
	validateID := utils.ValidateStruct(req)
	if validateID != nil {
		return entities.User{}, utils.NewHTTPError(utils.StatusBadRequest, "Invalid parameters", validateID, nil)
	}

	user := entities.User{
		ID:        req.ID,
		Lastname:  req.Lastname,
		Firstname: req.Firstname,
		Password:  req.Password,
		Username:  req.Username,
	}

	if err := us.userRepository.Update(&user); err != nil {
		return entities.User{}, utils.NewHTTPError(utils.StatusInternalServerError, "Database error", "Error during user update", err)
	}

	if user.ID == "" {
		return entities.User{}, utils.NewHTTPError(utils.StatusNotFound, "No user found", nil, nil)
	}

	return user, nil
}

// UpdatePassword updates user password
func (us userService) UpdatePassword(req requests.UserPasswordUpdate) *utils.HTTPError {
	validateReq := utils.ValidateStruct(req)
	if validateReq != nil {
		return utils.NewHTTPError(utils.StatusBadRequest, "Invalid parameters", validateReq, nil)
	}

	// Update user password
	userID, currentPassword, err := us.userRepository.GetIDFromPasswordReset(req.Token, req.Password)
	if err != nil {
		return utils.NewHTTPError(utils.StatusInternalServerError, "Database error", "Error when searching user", err)
	}
	if userID == "" {
		return utils.NewHTTPError(utils.StatusNotFound, "No user found", nil, nil)
	}

	// Change by the same password is forbidden
	hashedPassword := sha512.Sum512([]byte(req.Password))
	if hex.EncodeToString(hashedPassword[:]) == currentPassword {
		return utils.NewHTTPError(utils.StatusBadRequest, "New password cannot be the same as the current one", nil, nil)
	}

	err = us.userRepository.UpdatePassword(userID, currentPassword, req.Password)
	if err != nil {
		return utils.NewHTTPError(utils.StatusInternalServerError, "Database error", "Error when updating user password", err)
	}

	// Delete password reset
	err = us.userRepository.DeletePasswordReset(userID)
	if err != nil {
		return utils.NewHTTPError(utils.StatusInternalServerError, "Database error", "Error when deleting user password reset", err)
	}

	return nil
}

// ForgottenPassword save a forgotten password request
func (us userService) ForgottenPassword(req requests.UserForgotPassword) (entities.PasswordResets, *utils.HTTPError) {
	validateReq := utils.ValidateStruct(req)
	if validateReq != nil {
		return entities.PasswordResets{}, utils.NewHTTPError(utils.StatusBadRequest, "Invalid parameters", validateReq, nil)
	}

	// Find user
	user, err := us.userRepository.GetByUsername(req.Email)
	if err != nil {
		return entities.PasswordResets{}, utils.NewHTTPError(utils.StatusInternalServerError, "Database error", "Error when retrieving user", err)
	}
	if user.ID == "" {
		return entities.PasswordResets{}, utils.NewHTTPError(utils.StatusNotFound, "No user found", nil, nil)
	}

	// Create password reset
	passwordReset := entities.PasswordResets{
		UserID:    user.ID,
		Token:     uuid.NewString(),
		ExpiredAt: time.Now().Add(viper.GetDuration("FORGOTTEN_PASSWORD_EXPIRATION_DURATION") * time.Hour).UTC(),
	}
	err = us.userRepository.CreateOrUpdatePasswordReset(passwordReset)
	if err != nil {
		return entities.PasswordResets{}, utils.NewHTTPError(utils.StatusInternalServerError, "Database error", "Error when requesting new password", err)
	}

	// Send email with link
	to := make([]string, 1)
	to[0] = user.Username
	subject := fmt.Sprintf("[%s] Forgotten password", viper.GetString("APP_NAME"))
	var body bytes.Buffer

	tp, err := template.ParseFiles("templates/forgotten_password.gohtml")
	if err != nil {
		return entities.PasswordResets{}, utils.NewHTTPError(utils.StatusInternalServerError, "Email error", "Error when creating password reset email", err)
	}
	err = tp.Execute(&body, struct {
		Title string
		Link  string
	}{
		Title: fmt.Sprintf("%s - Forgotten password", viper.GetString("APP_NAME")),
		Link:  fmt.Sprintf("%s/%s", viper.GetString("FORGOTTEN_PASSWORD_BASE_URL"), passwordReset.Token),
	})
	if err != nil {
		return entities.PasswordResets{}, utils.NewHTTPError(utils.StatusInternalServerError, "Email error", "Error when creating password reset email", err)
	}

	err = mail.Send(
		viper.GetString("FORGOTTEN_PASSWORD_EMAIL_FROM"),
		to,
		nil,
		nil,
		subject,
		body.String(),
		"",
		"",
		viper.GetString("SMTP_HOST"),
		viper.GetInt("SMTP_PORT"))
	if err != nil {
		return entities.PasswordResets{}, utils.NewHTTPError(utils.StatusInternalServerError, "Email error", "Error when sending password reset email", err)
	}

	return passwordReset, nil
}
