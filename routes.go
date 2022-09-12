package server

import (
	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"

	"github.com/fabienbellanger/fiber-boilerplate/db"
	"github.com/fabienbellanger/fiber-boilerplate/deliveries/task"
	"github.com/fabienbellanger/fiber-boilerplate/deliveries/user"
	"github.com/fabienbellanger/fiber-boilerplate/deliveries/web"
	storeTask "github.com/fabienbellanger/fiber-boilerplate/stores/task"
	storeUser "github.com/fabienbellanger/fiber-boilerplate/stores/user"
)

// Web routes
// ----------

func registerPublicWebRoutes(r fiber.Router, logger *zap.Logger) {
	r.Get("/health-check", web.HealthCheck(logger))
	r.Get("/hello/:name", web.Hello())
}

// API routes
// ----------

func registerPublicAPIRoutes(r fiber.Router, db *db.DB, logger *zap.Logger) {
	v1 := r.Group("/v1")
	userStore := storeUser.New(db)

	// Login
	authGroup := v1.Group("")
	auth := user.New(authGroup, userStore, logger)
	authGroup.Post("/login", auth.Login)

	// Tasks
	registerTask(v1, db, logger)
}

func registerProtectedAPIRoutes(r fiber.Router, db *db.DB, logger *zap.Logger) {
	v1 := r.Group("/v1")
	userStore := storeUser.New(db)

	// Register
	registerGroup := v1.Group("")
	register := user.New(registerGroup, userStore, logger)
	registerGroup.Post("/register", register.Create)

	// Users
	userGroup := r.Group("/users")
	users := user.New(userGroup, userStore, logger)
	users.Routes()
}

func registerTask(r fiber.Router, db *db.DB, logger *zap.Logger) {
	taskGroup := r.Group("/tasks")
	taskStore := storeTask.New(db)
	tasks := task.New(taskGroup, taskStore, db, logger)
	tasks.Routes()
}
