package responses

import "github.com/fabienbellanger/fiber-boilerplate/pkg/domain/entities"

// TaskGetAll response
type TaskGetAll struct {
	Data  []entities.Task `json:"data"`
	Total int64           `json:"total"`
}
