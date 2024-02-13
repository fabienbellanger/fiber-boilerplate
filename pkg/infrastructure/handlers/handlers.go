package handlers

import (
	"errors"
	"github.com/fabienbellanger/fiber-boilerplate/utils"
	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

func ManageError(err *utils.HTTPError, c *fiber.Ctx, logger *zap.Logger) error {
	if errors.Is(err, utils.HTTPError{}) && err.Err != nil {
		if details, ok := err.Details.(string); ok {
			return utils.NewError(c, logger, err.Message, details, err.Err)
		}
	}
	return c.Status(err.Code).JSON(err)
}
