package services

import (
	"errors"
	"github.com/fabienbellanger/fiber-boilerplate/pkg/domain/entities"
	"github.com/fabienbellanger/fiber-boilerplate/pkg/domain/repositories"
	"github.com/fabienbellanger/fiber-boilerplate/pkg/domain/requests"
	"github.com/fabienbellanger/fiber-boilerplate/pkg/domain/responses"
	"github.com/fabienbellanger/fiber-boilerplate/utils"
	"github.com/spf13/viper"
	"gorm.io/gorm"
)

type UserService interface {
	Login(req requests.UserLogin) (responses.UserLogin, *utils.HTTPError)
	Create(req requests.UserEdit) (entities.User, *utils.HTTPError)
	GetAll() ([]entities.User, *utils.HTTPError)
	GetByID(id requests.UserByID) (entities.User, *utils.HTTPError)
	Delete(id requests.UserByID) *utils.HTTPError
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
func (us userService) Create(req requests.UserEdit) (entities.User, *utils.HTTPError) {
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
