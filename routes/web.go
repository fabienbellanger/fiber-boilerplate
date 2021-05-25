package routes

import (
	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"

	"github.com/fabienbellanger/fiber-boilerplate/handlers/web"
	"github.com/fabienbellanger/fiber-boilerplate/ws"
)

// RegisterPublicWebRoutes lists all public Web routes.
func RegisterPublicWebRoutes(r fiber.Router, logger *zap.Logger) {
	r.Get("/health-check", web.HealthCheck(logger))
}

// RegisterPublicWebSocketRoutes lists all public Web Socket routes.
func RegisterPublicWebSocketRoutes(r fiber.Router, hub *ws.Hub) {
	w := r.Group("/ws")

	// Access the websocket server: ws://localhost:3000/ws/123?v=1.0
	// https://www.websocket.org/echo.html
	w.Get("/", func(c *fiber.Ctx) error {
		return ws.ServeWs(c, hub)
	})
}
