package routes

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/websocket/v2"

	"github.com/fabienbellanger/fiber-boilerplate/handlers/web"
	"github.com/fabienbellanger/fiber-boilerplate/ws"
)

// RegisterPublicWebRoutes lists all public Web routes.
func RegisterPublicWebRoutes(r fiber.Router) {
	r.Get("/health-check", web.HealthCheck())
}

// RegisterPublicWebSocketRoutes lists all public Web Socket routes.
func RegisterPublicWebSocketRoutes(r fiber.Router, hub *ws.Hub) {
	w := r.Group("/ws")

	r.Use("/ws", func(c *fiber.Ctx) error {
		// IsWebSocketUpgrade returns true if the client
		// requested upgrade to the WebSocket protocol.
		if websocket.IsWebSocketUpgrade(c) {
			c.Locals("allowed", true)
			return c.Next()
		}
		return fiber.ErrUpgradeRequired
	})

	// Access the websocket server: ws://localhost:3000/ws/123?v=1.0
	// https://www.websocket.org/echo.html
	w.Get("/:id", websocket.New(func(c *websocket.Conn) {
		// c.Locals is added to the *websocket.Conn
		log.Printf("allowed: %v, params: %v, query: %v\n", c.Locals("allowed"), c.Params("id"), c.Query("v"))

		ws.ServeWs(hub, c)
	}))
}
