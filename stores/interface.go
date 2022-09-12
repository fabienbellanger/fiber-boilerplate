package stores

import (
	"database/sql"

	"github.com/fabienbellanger/fiber-boilerplate/entities"
)

// UserStorer interface
type UserStorer interface {
	Login(username, password string) (entities.User, error)
	Create(user *entities.User) error
	GetAll() ([]entities.User, error)
	GetOne(id string) (entities.User, error)
	Delete(id string) error
	Update(id string, userForm *entities.UserForm) (entities.User, error)
}

// TaskStorer interface
type TaskStorer interface {
	ListAll() ([]entities.Task, error)
	ListAllRows() (*sql.Rows, error)
	Create(task *entities.Task) error
}
