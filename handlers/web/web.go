package web

import (
	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

// HealthCheck returns status code 200.
func HealthCheck(logger *zap.Logger) fiber.Handler {
	return func(c *fiber.Ctx) error {
		logger.Debug("Health check OK")
		return c.SendStatus(fiber.StatusOK)
	}
}
