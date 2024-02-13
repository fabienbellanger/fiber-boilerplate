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
	Create(req requests.UserCreation) (entities.User, *utils.HTTPError)
	GetAll(req requests.Pagination) (responses.UsersListPaginated, *utils.HTTPError)
	GetByID(id requests.UserByID) (entities.User, *utils.HTTPError)
	Delete(id requests.UserByID) *utils.HTTPError
	Update(req requests.UserUpdate) (entities.User, *utils.HTTPError)
	UpdatePassword(req requests.UserPasswordUpdate) *utils.HTTPError
	ForgottenPassword(req requests.UserForgotPassword) (entities.PasswordResets, *utils.HTTPError)
}

type userUseCase struct {
	userService services.UserService
}

// NewUser returns a new User use case
func NewUser(userService services.UserService) User {
	return &userUseCase{userService}
}

// Login user
func (uc *userUseCase) Login(req requests.UserLogin) (responses.UserLogin, *utils.HTTPError) {
	return uc.userService.Login(req)
}

// Create user
func (uc *userUseCase) Create(req requests.UserCreation) (entities.User, *utils.HTTPError) {
	return uc.userService.Create(req)
}

// GetAll users
func (uc *userUseCase) GetAll(req requests.Pagination) (responses.UsersListPaginated, *utils.HTTPError) {
	return uc.userService.GetAll(req)
}

// GetByID user
func (uc *userUseCase) GetByID(id requests.UserByID) (entities.User, *utils.HTTPError) {
	return uc.userService.GetByID(id)
}

// Delete user
func (uc *userUseCase) Delete(id requests.UserByID) *utils.HTTPError {
	return uc.userService.Delete(id)
}

// Update user
func (uc *userUseCase) Update(req requests.UserUpdate) (entities.User, *utils.HTTPError) {
	return uc.userService.Update(req)
}

// UpdatePassword user
func (uc *userUseCase) UpdatePassword(req requests.UserPasswordUpdate) *utils.HTTPError {
	return uc.userService.UpdatePassword(req)
}

// ForgottenPassword user
func (uc *userUseCase) ForgottenPassword(req requests.UserForgotPassword) (entities.PasswordResets, *utils.HTTPError) {
	return uc.userService.ForgottenPassword(req)
}
