package repositories

import (
	"github.com/fabienbellanger/fiber-boilerplate/pkg/domain/entities"
	"github.com/fabienbellanger/fiber-boilerplate/pkg/domain/requests"
)

// UserRepository is the interface that wraps the basic user repository methods.
type UserRepository interface {
	Login(username, password string) (entities.User, error)
	Create(user *entities.User) error
	GetAll() ([]entities.User, error)
	GetByID(id string) (entities.User, error)
	GetByUsername(username string) (user entities.User, err error)
	Delete(id string) error
	Update(id string, user requests.UserEdit) (entities.User, error)
	UpdatePassword(id, currentPassword, password string) error
	GetIDFromPasswordReset(token, password string) (string, string, error)
	DeletePasswordReset(userId string) error
	CreateOrUpdatePasswordReset(passwordReset entities.PasswordResets) error
}
