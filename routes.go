package server

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/basicauth"
	"github.com/gofiber/fiber/v2/middleware/filesystem"
	"github.com/spf13/viper"
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
	userStore := storeUser.New(db)

	// Login
	authGroup := v1.Group("")
	auth := user.New(authGroup, userStore, logger)
	authGroup.Post("/login", auth.Login)

	// Password reset
	authGroup.Post("/forgotten-password/:email", auth.ForgottenPassword)
	authGroup.Patch("/update-password/:token", auth.UpdatePassword)

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
	userGroup := v1.Group("/users")
	users := user.New(userGroup, userStore, logger)
	users.Routes()
}

func registerTask(r fiber.Router, db *db.DB, logger *zap.Logger) {
	taskGroup := r.Group("/tasks")
	taskStore := storeTask.New(db)
	tasks := task.New(taskGroup, taskStore, db, logger)
	tasks.Routes()
}
