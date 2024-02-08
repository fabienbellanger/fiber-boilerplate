package router

import (
	"github.com/fabienbellanger/fiber-boilerplate/pkg/adapters/db"
	"github.com/fabienbellanger/fiber-boilerplate/pkg/adapters/stores"
	"github.com/fabienbellanger/fiber-boilerplate/pkg/domain/services"
	"github.com/fabienbellanger/fiber-boilerplate/pkg/domain/usecases"
	"github.com/fabienbellanger/fiber-boilerplate/pkg/infrastructure/handlers/api"
	"github.com/fabienbellanger/fiber-boilerplate/pkg/infrastructure/handlers/web"
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/basicauth"
	"github.com/gofiber/fiber/v2/middleware/filesystem"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

// Web routes
// ----------

func registerPublicWebRoutes(r fiber.Router, logger *zap.Logger) {
	// Basic Auth
	// ----------
	cfg := basicauth.Config{
		Users: map[string]string{
			viper.GetString("SERVER_BASICAUTH_USERNAME"): viper.GetString("SERVER_BASICAUTH_PASSWORD"),
		},
	}

	// API documentation
	doc := r.Group("/doc")
	doc.Use(basicauth.New(cfg))
	doc.Get("/api-v1", web.DocAPIv1())

	r.Get("/health-check", web.HealthCheck(logger))

	// Filesystem
	// ----------
	assets := r.Group("/assets")
	assets.Use(filesystem.New(filesystem.Config{
		Root:   http.Dir("./assets"),
		Browse: false,
		Index:  "index.html",
		MaxAge: 3600,
	}))
}

// API routes
// ----------

func registerPublicAPIRoutes(r fiber.Router, db *db.DB, logger *zap.Logger) {
	v1 := r.Group("/v1")
	userStore := stores.NewUserStore(db)
	userService := services.NewUser(userStore)
	userUserCase := usecases.NewUser(userService)

	// Login & password reset
	users := api.NewUser(v1, userUserCase, logger)
	users.UserPublicRoutes()
}

func registerProtectedAPIRoutes(r fiber.Router, db *db.DB, logger *zap.Logger) {
	v1 := r.Group("/v1")

	// Users
	registerUser(v1, db, logger)

	// Tasks
	registerTask(v1, db, logger)
}

func registerUser(r fiber.Router, db *db.DB, logger *zap.Logger) {
	userStore := stores.NewUserStore(db)
	userService := services.NewUser(userStore)
	userUserCase := usecases.NewUser(userService)

	// Users
	userGroup := r.Group("/users")
	users := api.NewUser(userGroup, userUserCase, logger)
	users.UserProtectedRoutes()
}

func registerTask(r fiber.Router, db *db.DB, logger *zap.Logger) {
	taskGroup := r.Group("/tasks")
	taskStore := stores.NewTaskStore(db)
	taskService := services.NewTask(taskStore)
	taskUserCase := usecases.NewTask(taskService)

	tasks := api.NewTask(taskGroup, taskUserCase, logger)
	tasks.TaskProtectedRoutes()
}
