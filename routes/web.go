package routes

import (
	"github.com/gofiber/fiber/v2"

	"github.com/fabienbellanger/fiber-boilerplate/handlers/web"
)

// RegisterPublicWebRoutes lists all public Web routes.
func RegisterPublicWebRoutes(r fiber.Router) {
	r.Get("/health-check", web.HealthCheck())
}
