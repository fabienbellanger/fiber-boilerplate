package utils

import (
	"fmt"
)

// Error represents a custom error.
type Error struct {
	Code int
	Err  error
}

func (e *Error) Error() string {
	return fmt.Sprintf("Code: %d, Err=%v", e.Code, e.Err)
}
