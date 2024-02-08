package usecases

import (
	"database/sql"
	"github.com/fabienbellanger/fiber-boilerplate/pkg/domain/entities"
	"github.com/fabienbellanger/fiber-boilerplate/pkg/domain/requests"
	"github.com/fabienbellanger/fiber-boilerplate/pkg/domain/responses"
	"github.com/fabienbellanger/fiber-boilerplate/pkg/domain/services"
	"github.com/fabienbellanger/fiber-boilerplate/utils"
)

type Task interface {
	GetAll(req requests.Pagination) (responses.TaskGetAll, *utils.HTTPError)
	Create(req requests.TaskCreation) (entities.Task, *utils.HTTPError)
	GetAllStream() (*sql.Rows, *utils.HTTPError)
	ScanTask(rows *sql.Rows, task *entities.Task) *utils.HTTPError
}

type taskUseCase struct {
	taskService services.TaskService
}

// NewTask returns a new Task use case
func NewTask(taskService services.TaskService) Task {
	return &taskUseCase{taskService}
}

// GetAll tasks
func (uc *taskUseCase) GetAll(req requests.Pagination) (responses.TaskGetAll, *utils.HTTPError) {
	return uc.taskService.GetAll(req)
}

// Create task
func (uc *taskUseCase) Create(req requests.TaskCreation) (entities.Task, *utils.HTTPError) {
	return uc.taskService.Create(req)
}

// GetAllStream tasks
func (uc *taskUseCase) GetAllStream() (*sql.Rows, *utils.HTTPError) {
	return uc.taskService.GetAllStream()
}

// ScanTask tasks
func (uc *taskUseCase) ScanTask(rows *sql.Rows, task *entities.Task) *utils.HTTPError {
	return uc.taskService.ScanTask(rows, task)
}
