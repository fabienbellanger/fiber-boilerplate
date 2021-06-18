package repositories

import (
	"database/sql"

	"github.com/fabienbellanger/fiber-boilerplate/db"
	"github.com/fabienbellanger/fiber-boilerplate/models"
	"github.com/google/uuid"
)

// ListAllTasks gets all users in database.
func ListAllTasks(db *db.DB) ([]models.Task, error) {
	var tasks []models.Task

	if response := db.Find(&tasks); response.Error != nil {
		return tasks, response.Error
	}
	return tasks, nil
}

// ListAllTasksRows gets all users in database.
func ListAllTasksRows(db *db.DB) (*sql.Rows, error) {
	return db.Model(&models.Task{}).Where("deleted_at IS NULL").Rows()
}

// CreateTask adds task in database.
func CreateTask(db *db.DB, task *models.Task) error {
	// UUID
	// ----
	task.ID = uuid.New().String()

	if result := db.Create(&task); result.Error != nil {
		return result.Error
	}
	return nil
}
