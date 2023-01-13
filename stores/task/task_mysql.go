package task

import (
	"database/sql"

	"github.com/fabienbellanger/fiber-boilerplate/db"
	"github.com/fabienbellanger/fiber-boilerplate/entities"
	"github.com/google/uuid"
)

// TaskStore ...
type TaskStore struct {
	db *db.DB
}

// New returns a new TaskStore
func New(db *db.DB) TaskStore {
	return TaskStore{db: db}
}

// ListAll gets all users in database.
func (t TaskStore) ListAll() ([]entities.Task, error) {
	var tasks []entities.Task

	if response := t.db.Find(&tasks); response.Error != nil {
		return tasks, response.Error
	}
	return tasks, nil
}

// ListAllRows gets all users in database.
func (t TaskStore) ListAllRows() (*sql.Rows, error) {
	return t.db.Model(&entities.Task{}).Where("deleted_at IS NULL").Rows()
}

// Create a new task in database.
func (t TaskStore) Create(task *entities.Task) error {
	// UUID
	// ----
	task.ID = uuid.New().String()

	if result := t.db.Create(&task); result.Error != nil {
		return result.Error
	}
	return nil
}

func (t TaskStore) ScanRow(rows *sql.Rows, task *entities.Task) error {
	return t.db.ScanRows(rows, &task)
}
