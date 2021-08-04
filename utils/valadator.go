package utils

import (
	"github.com/go-playground/validator/v10"
)

// ValidatorError represents error validation struct.
type ValidatorError struct {
	FailedField string
	Tag         string
	Value       string
}

// ValidateStruct checks if a struct is valid and returns an array of errors
// if it is not valid.
func ValidateStruct(task interface{}) (errors []*ValidatorError) {
	validate := validator.New()
	err := validate.Struct(task)
	if err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			errors = append(errors, &ValidatorError{
				FailedField: err.Field(),
				Tag:         err.Tag(),
				Value:       err.Param(),
			})
		}
	}
	return
}
