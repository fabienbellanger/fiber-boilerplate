package values_objects

import "github.com/fabienbellanger/fiber-boilerplate/utils"

// Email represents an email value object
type Email struct {
	Value string `validate:"required,email"`
}

// String returns the email value
func (e *Email) String() string {
	return e.Value
}

// NewEmail creates a new email
func NewEmail(value string) (Email, error) {
	p := Email{Value: value}

	err := p.Validate()
	if err != nil {
		return Email{}, err
	}

	return Email{Value: value}, nil
}

// Validate checks if a struct is valid and returns an array of errors
func (e *Email) Validate() utils.ValidatorErrors {
	return utils.ValidateStruct(e)
}
