package utils

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

// HTTPError represents an HTTP error.
type HTTPError struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Details interface{} `json:"details,omitempty"`
}

// NewError returns a fiber error and log the error.
func NewError(c *fiber.Ctx, logger *zap.Logger, desc, msg string, err error) *fiber.Error {
	logger.Error(desc, zap.Error(err), zap.String("requestId", fmt.Sprintf("%v", c.Locals("requestid"))))

	return fiber.NewError(fiber.StatusInternalServerError, msg)
}
