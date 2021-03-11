package web

import (
	"github.com/gofiber/fiber/v2"

	"github.com/fabienbellanger/fiber-boilerplate/logger"
)

// HealthCheck returns status code 200.
func HealthCheck() fiber.Handler {
	return func(c *fiber.Ctx) error {
		return &logger.Error{
			Code: 1,
			Err: &logger.Error{
				Code: 234,
				Err:  c.SendStatus(fiber.StatusOK),
			},
		}
	}
}
