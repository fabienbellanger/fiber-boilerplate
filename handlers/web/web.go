package web

import "github.com/gofiber/fiber/v2"

// HealthCheckHandler returns status code 200.
func HealthCheckHandler() fiber.Handler {
	return func(c *fiber.Ctx) error {
		return c.SendStatus(fiber.StatusOK)
	}
}
