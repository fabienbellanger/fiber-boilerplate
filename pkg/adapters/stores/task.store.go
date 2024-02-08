package stores

import (
	"database/sql"
	"github.com/fabienbellanger/fiber-boilerplate/pkg/adapters/db"
	"github.com/fabienbellanger/fiber-boilerplate/pkg/domain/entities"

	"github.com/google/uuid"
)

// TaskStore ...
type TaskStore struct {
	db *db.DB
}

// NewTaskStore returns a new TaskStore
func NewTaskStore(db *db.DB) TaskStore {
	return TaskStore{db: db}
}

// GetAll gets all users in database.
func (t TaskStore) GetAll(page, limit, sorts string) (tasks []entities.Task, total int64, err error) {
	// Total rows
	t.db.Model(&tasks).Count(&total)

	q := t.db.Scopes(db.Paginate(page, limit))
	q.Scopes(db.Order(sorts))
	if response := q.Find(&tasks); response.Error != nil {
		return tasks, total, response.Error
	}
	return tasks, total, nil
}

// GetAllRows gets all tasks in database.
func (t TaskStore) GetAllRows() (*sql.Rows, error) {
	return t.db.Model(&entities.Task{}).Where("deleted_at IS NULL").Rows()
}

// Create a new task in database.
func (t TaskStore) Create(task *entities.Task) error {
	// UUID
	// ----
	task.ID = uuid.NewString()

	if result := t.db.Create(&task); result.Error != nil {
		return result.Error
	}
	return nil
}

func (t TaskStore) ScanRow(rows *sql.Rows, task *entities.Task) error {
	return t.db.ScanRows(rows, &task)
}
