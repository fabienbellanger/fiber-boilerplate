package utils

import (
	"fmt"
)

// TODO: Improve error handling

// Error represents a custom error.
type Error struct {
	Code int
	Msg  string
	Err  *error
}

func (e *Error) Error() string {
	return fmt.Sprintf("Code: %d, Message: %s, Err=%v", e.Code, e.Msg, *e.Err)
}
