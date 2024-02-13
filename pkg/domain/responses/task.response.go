package responses

import "github.com/fabienbellanger/fiber-boilerplate/pkg/domain/entities"

// TasksListPaginated response
type TasksListPaginated struct {
	Data  []entities.Task `json:"data"`
	Total int64           `json:"total"`
}
