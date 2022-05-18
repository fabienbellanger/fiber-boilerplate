package web

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

// HealthCheck returns status code 200.
func HealthCheck(logger *zap.Logger) fiber.Handler {
	return func(c *fiber.Ctx) error {
		log.Printf("%v\n", c.IP())
		return c.SendStatus(fiber.StatusOK)
	}
}

// Hello is a test for pkger template.
func Hello() fiber.Handler {
	return func(c *fiber.Ctx) error {
		name := c.Params("name")
		return c.Render("hello", fiber.Map{
			"name": name,
		})
	}
}
