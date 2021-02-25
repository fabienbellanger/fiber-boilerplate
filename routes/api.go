package routes

import (
	"github.com/fabienbellanger/fiber-boilerplate/db"
	"github.com/fabienbellanger/fiber-boilerplate/handlers/api"
	"github.com/gofiber/fiber/v2"
)

// RegisterPublicAPIRoutes lists all public API routes.
func RegisterPublicAPIRoutes(r fiber.Router, db *db.DB) {
	registerUser(r, db)
}

// RegisterProtectedAPIRoutes lists all protected API routes.
func RegisterProtectedAPIRoutes(r fiber.Router, db *db.DB) {

}

func registerUser(r fiber.Router, db *db.DB) {
	users := r.Group("users")

	users.Get("", api.GetAllUsers(db))
	users.Post("", api.CreateUser(db))
}
