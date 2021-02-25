package routes

import (
	"github.com/fabienbellanger/fiber-boilerplate/db"
	"github.com/fabienbellanger/fiber-boilerplate/handlers/web"
	"github.com/gofiber/fiber/v2"
)

// RegisterPublicWebRoutes lists all public Web routes.
func RegisterPublicWebRoutes(r fiber.Router, db *db.DB) {
	r.Get("/health_check", web.HealthCheckHandler())
}

// RegisterProtectedWebRoutes lists all protected Web routes.
func RegisterProtectedWebRoutes(r fiber.Router, db *db.DB) {

}
