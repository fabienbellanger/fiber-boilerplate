package services

import (
	"database/sql"
	"github.com/fabienbellanger/fiber-boilerplate/pkg/domain/entities"
	"github.com/fabienbellanger/fiber-boilerplate/pkg/domain/repositories"
	"github.com/fabienbellanger/fiber-boilerplate/pkg/domain/requests"
	"github.com/fabienbellanger/fiber-boilerplate/pkg/domain/responses"
	"github.com/fabienbellanger/fiber-boilerplate/utils"
)

type TaskService interface {
	GetAll(req requests.Pagination) (responses.TasksListPaginated, *utils.HTTPError)
	Create(req requests.TaskCreation) (entities.Task, *utils.HTTPError)
	GetAllStream() (*sql.Rows, *utils.HTTPError)
	ScanTask(rows *sql.Rows, task *entities.Task) *utils.HTTPError
}

type taskService struct {
	taskRepository repositories.TaskRepository
}

// NewTask returns a new user service
func NewTask(repo repositories.TaskRepository) TaskService {
	return &taskService{repo}
}

// GetAll tasks
func (ts taskService) GetAll(req requests.Pagination) (responses.TasksListPaginated, *utils.HTTPError) {
	tasks, total, err := ts.taskRepository.GetAll(req.Page, req.Limit, req.Sorts)
	if err != nil {
		return responses.TasksListPaginated{}, utils.NewHTTPError(utils.StatusInternalServerError, "Database error", "Error during tasks list", err)
	}

	return responses.TasksListPaginated{
		Data:  tasks,
		Total: total,
	}, nil
}

// Create task
func (ts taskService) Create(req requests.TaskCreation) (entities.Task, *utils.HTTPError) {
	validateReq := utils.ValidateStruct(req)
	if validateReq != nil {
		return entities.Task{}, utils.NewHTTPError(utils.StatusBadRequest, "Invalid parameters", validateReq, nil)
	}

	if req.Name == "" {
		return entities.Task{}, utils.NewHTTPError(utils.StatusBadRequest, "Name cannot be empty", validateReq, nil)
	}

	newTask := entities.Task{
		Name:        req.Name,
		Description: req.Description,
	}

	if err := ts.taskRepository.Create(&newTask); err != nil {
		return entities.Task{}, utils.NewHTTPError(utils.StatusInternalServerError, "Database error", "Error during task creation", err)
	}

	return newTask, nil
}

// GetAllStream tasks list
func (ts taskService) GetAllStream() (*sql.Rows, *utils.HTTPError) {
	rows, err := ts.taskRepository.GetAllRows()
	if err != nil {
		return nil, utils.NewHTTPError(utils.StatusInternalServerError, "Database error", "Error during tasks list with stream", err)
	}

	return rows, nil
}

// ScanTask scans a row
func (ts taskService) ScanTask(rows *sql.Rows, task *entities.Task) *utils.HTTPError {
	if err := ts.taskRepository.ScanRow(rows, task); err != nil {
		return utils.NewHTTPError(utils.StatusInternalServerError, "Database error", "Error during task scan", err)
	}

	return nil
}
