package web

import "github.com/gofiber/fiber/v2"

// HealthCheck returns status code 200.
func HealthCheck() fiber.Handler {
	return func(c *fiber.Ctx) error {
		return c.SendStatus(fiber.StatusOK)
	}
}
