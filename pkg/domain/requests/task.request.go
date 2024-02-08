package requests

// TaskCreation request to create a task
type TaskCreation struct {
	Name        string `json:"name" xml:"name" form:"name" validate:"required,min=3,max=127"`
	Description string `json:"description" xml:"description" form:"description"`
}
