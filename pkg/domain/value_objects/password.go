package values_objects

import (
	"github.com/fabienbellanger/fiber-boilerplate/utils"
)

// Password represents an password value object
type Password struct {
	Value string `validate:"required,min=8"`
}

// String returns the password value
func (p *Password) String() string {
	return p.Value
}

// NewPassword creates a new password
func NewPassword(value string) (Password, error) {
	p := Password{Value: value}

	err := p.Validate()
	if err != nil {
		return Password{}, err
	}

	return p, nil
}

// Validate checks if a struct is valid and returns an array of errors
func (p *Password) Validate() utils.ValidatorErrors {
	return utils.ValidateStruct(p)
}
