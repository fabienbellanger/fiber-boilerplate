package server

import (
	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"

	"github.com/fabienbellanger/fiber-boilerplate/db"
	"github.com/fabienbellanger/fiber-boilerplate/handlers/api"
	"github.com/fabienbellanger/fiber-boilerplate/handlers/web"
	"github.com/fabienbellanger/fiber-boilerplate/ws"
)

// Web routes
// ----------

func registerPublicWebRoutes(r fiber.Router, logger *zap.Logger, hub *ws.Hub) {
	r.Get("/health-check", web.HealthCheck(logger))

	registerPublicWebSocketRoutes(r, hub)
}

func registerPublicWebSocketRoutes(r fiber.Router, hub *ws.Hub) {
	w := r.Group("/ws")

	// Access the websocket server: ws://localhost:3000/ws/123?v=1.0
	// Tests with: https://www.websocket.org/echo.html
	w.Get("/", func(c *fiber.Ctx) error {
		return ws.ServeWs(c, hub)
	})
}

// API routes
// ----------

func registerPublicAPIRoutes(r fiber.Router, db *db.DB) {
	registerAuth(r, db)

	v1 := r.Group("/v1")
	registerTask(v1, db)
}

func registerProtectedAPIRoutes(r fiber.Router, db *db.DB) {
	v1 := r.Group("/v1")

	registerUser(v1, db)
}

func registerAuth(r fiber.Router, db *db.DB) {
	r.Post("/login", api.Login(db))
	r.Post("/register", api.CreateUser(db))
}

func registerUser(r fiber.Router, db *db.DB) {
	users := r.Group("/users")

	users.Get("/", api.GetAllUsers(db))
	users.Get("/stream", api.StreamUsers(db))
	users.Get("/:id", api.GetUser(db))
	users.Delete("/:id", api.DeleteUser(db))
	users.Put("/:id", api.UpdateUser(db))
}

func registerTask(r fiber.Router, db *db.DB) {
	tasks := r.Group("/tasks")

	tasks.Get("/", api.GetAllTasks(db))
	tasks.Get("/stream", api.GetAllTasksStream(db))
	tasks.Post("/", api.CreateTask(db))
}
