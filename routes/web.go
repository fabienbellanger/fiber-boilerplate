package routes

import "github.com/gofiber/fiber/v2"

// RegisterPublicWebRoutes lists all public Web routes.
func RegisterPublicWebRoutes(r fiber.Router) {
	// TODO: Mettre dans un handler proprement
	r.Get("/health_check", func(c *fiber.Ctx) error {
		return c.SendStatus(fiber.StatusOK)
	})
}

// RegisterProtectedWebRoutes lists all protected Web routes.
func RegisterProtectedWebRoutes(r fiber.Router) {

}
