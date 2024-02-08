package repositories

import (
	"database/sql"
	"github.com/fabienbellanger/fiber-boilerplate/pkg/domain/entities"
)

// TaskRepository is the interface that wraps the basic task repository methods.
type TaskRepository interface {
	GetAll(page, limit, sorts string) ([]entities.Task, int64, error)
	GetAllRows() (*sql.Rows, error)
	Create(task *entities.Task) error
	ScanRow(rows *sql.Rows, task *entities.Task) error
}
