package web

import (
	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

// HealthCheck returns status code 200.
func HealthCheck(logger *zap.Logger) fiber.Handler {
	return func(c *fiber.Ctx) error {
		return c.SendStatus(fiber.StatusOK)
	}
}

// DocAPIv1 show API v1 documentation.
func DocAPIv1() fiber.Handler {
	return func(c *fiber.Ctx) error {
		return c.Render("doc_api_v1", fiber.Map{})
	}
}
