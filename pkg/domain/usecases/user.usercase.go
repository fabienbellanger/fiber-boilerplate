package usecases

import (
	"github.com/fabienbellanger/fiber-boilerplate/pkg/domain/entities"
	"github.com/fabienbellanger/fiber-boilerplate/pkg/domain/requests"
	"github.com/fabienbellanger/fiber-boilerplate/pkg/domain/responses"
	"github.com/fabienbellanger/fiber-boilerplate/pkg/domain/services"
	"github.com/fabienbellanger/fiber-boilerplate/utils"
)

type User interface {
	Login(req requests.UserLogin) (responses.UserLogin, *utils.HTTPError)
	Create(req requests.UserEdit) (entities.User, *utils.HTTPError)
	GetAll() ([]entities.User, *utils.HTTPError)
	GetByID(id requests.UserByID) (entities.User, *utils.HTTPError)
	Delete(id requests.UserByID) *utils.HTTPError
}

type userUseCase struct {
	userService services.UserService
}

// NewUser returns a new CreateUser use case
func NewUser(userService services.UserService) User {
	return &userUseCase{userService}
}

// Login user
func (uc *userUseCase) Login(req requests.UserLogin) (responses.UserLogin, *utils.HTTPError) {
	return uc.userService.Login(req)
}

// Create user
func (uc *userUseCase) Create(req requests.UserEdit) (entities.User, *utils.HTTPError) {
	return uc.userService.Create(req)
}

// GetAll users
func (uc *userUseCase) GetAll() ([]entities.User, *utils.HTTPError) {
	return uc.userService.GetAll()
}

// GetByID user
func (uc *userUseCase) GetByID(id requests.UserByID) (entities.User, *utils.HTTPError) {
	return uc.userService.GetByID(id)
}

// Delete user
func (uc *userUseCase) Delete(id requests.UserByID) *utils.HTTPError {
	return uc.userService.Delete(id)
}
