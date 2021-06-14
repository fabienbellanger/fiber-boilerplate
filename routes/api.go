package routes

import (
	"github.com/fabienbellanger/fiber-boilerplate/db"
	"github.com/fabienbellanger/fiber-boilerplate/handlers/api"
	"github.com/gofiber/fiber/v2"
)

// RegisterPublicAPIRoutes lists all public API routes.
func RegisterPublicAPIRoutes(r fiber.Router, db *db.DB) {
	registerAuth(r, db)
}

// RegisterProtectedAPIRoutes lists all protected API routes.
func RegisterProtectedAPIRoutes(r fiber.Router, db *db.DB) {
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
