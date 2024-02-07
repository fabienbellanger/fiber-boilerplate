package router

import (
	"github.com/fabienbellanger/fiber-boilerplate/pkg/adapters/db"
	"github.com/fabienbellanger/fiber-boilerplate/pkg/adapters/stores"
	"github.com/fabienbellanger/fiber-boilerplate/pkg/domain/services"
	"github.com/fabienbellanger/fiber-boilerplate/pkg/domain/usecases"
	"github.com/fabienbellanger/fiber-boilerplate/pkg/infrastructure/handlers"
	"github.com/fabienbellanger/fiber-boilerplate/pkg/infrastructure/handlers/task"
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

	// Login
	auth := handlers.NewUser(v1, userUserCase, userStore, logger)
	v1.Post("/login", auth.Login)

	// Password reset
	v1.Post("/forgotten-password/:email", auth.ForgottenPassword)
	v1.Patch("/update-password/:token", auth.UpdatePassword)
}

func registerProtectedAPIRoutes(r fiber.Router, db *db.DB, logger *zap.Logger) {
	v1 := r.Group("/v1")

	// Register & Users
	registerUser(v1, db, logger)

	// Tasks
	registerTask(v1, db, logger)
}

func registerUser(r fiber.Router, db *db.DB, logger *zap.Logger) {
	userStore := stores.NewUserStore(db)
	userService := services.NewUser(userStore)
	userUserCase := usecases.NewUser(userService)

	// Register
	registerGroup := r.Group("")
	register := handlers.NewUser(registerGroup, userUserCase, userStore, logger)
	registerGroup.Post("/register", register.Create)

	// Users
	userGroup := r.Group("/users")
	users := handlers.NewUser(userGroup, userUserCase, userStore, logger)
	users.UserRoutes()
}

func registerTask(r fiber.Router, db *db.DB, logger *zap.Logger) {
	taskGroup := r.Group("/tasks")
	taskStore := stores.NewTaskStore(db)
	tasks := task.New(taskGroup, taskStore, logger)
	tasks.Routes()
}
