package server

import (
	"fmt"
	"os"
	"os/signal"
	"strings"
	"time"

	"github.com/ansrivas/fiberprometheus/v2"
	"github.com/fabienbellanger/goutils"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/basicauth"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/favicon"
	"github.com/gofiber/fiber/v2/middleware/filesystem"
	"github.com/gofiber/fiber/v2/middleware/limiter"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/monitor"
	"github.com/gofiber/fiber/v2/middleware/pprof"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/fiber/v2/middleware/requestid"
	jwtware "github.com/gofiber/jwt/v2"
	"github.com/gofiber/template/django"
	"github.com/markbates/pkger"
	"github.com/spf13/viper"
	"go.uber.org/zap"

	"github.com/fabienbellanger/fiber-boilerplate/db"
	"github.com/fabienbellanger/fiber-boilerplate/middlewares/timer"
	"github.com/fabienbellanger/fiber-boilerplate/utils"
	"github.com/fabienbellanger/fiber-boilerplate/ws"
)

// Run starts HTTP server.
func Run(db *db.DB, hub *ws.Hub, logger *zap.Logger) {
	app := fiber.New(initConfig(logger))

	initMiddlewares(app)
	initTools(app)

	// Routes
	// ------
	web := app.Group("")
	api := app.Group("api")

	// Public routes
	// -------------
	registerPublicWebRoutes(web, logger, hub)
	registerPublicAPIRoutes(api, db)

	// Protected routes
	// ----------------
	initJWT(app)
	registerProtectedAPIRoutes(api, db)

	// Custom 404 (after all routes but not available because of JWT)
	// --------------------------------------------------------------
	app.Use(func(ctx *fiber.Ctx) error {
		return ctx.Status(fiber.StatusNotFound).JSON(utils.HTTPError{
			Code:    fiber.StatusNotFound,
			Message: "Resource Not Found",
		})
	})

	// Close any connections on interrupt signal
	// -----------------------------------------
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		<-c
		app.Shutdown()
	}()

	// Run fiber server
	// ----------------
	err := app.Listen(fmt.Sprintf("%s:%s", viper.GetString("APP_ADDR"), viper.GetString("APP_PORT")))
	if err != nil {
		app.Shutdown()
	}
}

func initConfig(logger *zap.Logger) fiber.Config {
	// Initialize standard Go html template engine
	engine := django.NewFileSystem(pkger.Dir("/public/templates"), ".django")

	return fiber.Config{
		Prefork:               viper.GetBool("SERVER_PREFORK"),
		DisableStartupMessage: false,
		StrictRouting:         true,
		Views:                 engine,
		// Errors handling
		// ---------------
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			code := fiber.StatusInternalServerError

			// Fiber error
			// -----------
			e, ok := err.(*fiber.Error)
			if ok {
				code = e.Code
			}

			// Request ID
			// ----------
			requestID := c.Locals("requestid")

			// Custom Fiber error
			// ------------------
			if e != nil {
				logger.Error(
					"HTTP error",
					zap.Error(e),
					zap.Int("code", code),
					zap.String("method", c.Method()),
					zap.String("path", c.Path()),
					zap.ByteString("body", c.Body()),
					zap.String("url", c.OriginalURL()),
					zap.String("host", c.BaseURL()),
					zap.String("ip", c.IP()),
					zap.String("requestId", fmt.Sprintf("%v", requestID)))

				return c.Status(code).JSON(e)
			}

			// Internal Server Error
			// ---------------------
			if code == fiber.StatusInternalServerError {
				logger.Error(
					"Internal server error",
					zap.Error(e),
					zap.Int("code", code),
					zap.String("method", c.Method()),
					zap.String("path", c.Path()),
					zap.ByteString("body", c.Body()),
					zap.String("url", c.OriginalURL()),
					zap.String("host", c.BaseURL()),
					zap.String("ip", c.IP()),
					zap.String("requestId", fmt.Sprintf("%v", requestID)))

				return c.Status(code).JSON(utils.HTTPError{
					Code:    code,
					Message: "Internal Server Error",
				})
			}
			return nil
		},
	}
}

func initMiddlewares(s *fiber.App) {
	// CORS
	// ----
	s.Use(cors.New(cors.Config{
		AllowOrigins:     strings.Join(viper.GetStringSlice("CORS_ALLOW_ORIGINS"), ", "),
		AllowMethods:     strings.Join(viper.GetStringSlice("CORS_ALLOW_METHODS"), ", "),
		AllowHeaders:     strings.Join(viper.GetStringSlice("CORS_ALLOW_HEADERS"), ", "),
		ExposeHeaders:    strings.Join(viper.GetStringSlice("CORS_EXPOSE_HEADERS"), ", "),
		AllowCredentials: viper.GetBool("CORS_ALLOW_CREDENTIALS"),
		MaxAge:           int(12 * time.Hour),
	}))

	// Favicon
	// -------
	s.Use(favicon.New(favicon.Config{
		File: "favicon.png",
	}))

	// Logger
	// ------
	if viper.GetString("APP_ENV") == "development" || viper.GetBool("ENABLE_ACCESS_LOG") {
		s.Use(logger.New(logger.Config{
			Next:         nil,
			Format:       "${time} | ${status} | ${method} | ${path} | ${protocol}://${host}${url} | ${latency} | ${locals:requestid}\n",
			TimeFormat:   "2006-01-02 15:04:05",
			TimeZone:     "Local",
			TimeInterval: 500 * time.Millisecond,
			Output:       os.Stderr,
		}))
	}

	// Recover
	// -------
	s.Use(recover.New())

	// Request ID
	// ----------
	s.Use(requestid.New())

	// Timer
	// -----
	s.Use(timer.New(timer.Config{
		DisplayMilliseconds: false,
		DisplaySeconds:      true,
	}))

	// Limiter
	// -------
	if viper.GetBool("LIMITER_ENABLE") {
		s.Use(limiter.New(limiter.Config{
			Next: func(c *fiber.Ctx) bool {
				excludedIP := viper.GetStringSlice("LIMITER_EXCLUDE_IP")
				if len(excludedIP) == 0 {
					return false
				}
				return goutils.StringInSlice(c.IP(), excludedIP)
			},
			Max:        viper.GetInt("LIMITER_MAX"),
			Expiration: viper.GetDuration("LIMITER_EXPIRATION") * time.Second,
			KeyGenerator: func(c *fiber.Ctx) string {
				return c.IP()
			},
			LimitReached: func(c *fiber.Ctx) error {
				return c.Status(fiber.StatusTooManyRequests).JSON(utils.HTTPError{
					Code:    fiber.StatusTooManyRequests,
					Message: "Too Many Requests",
				})
			},
		}))
	}
}

func initTools(s *fiber.App) {
	// Pkger
	// -----
	s.Use("/assets", filesystem.New(filesystem.Config{
		Root: pkger.Dir("/public/assets"),
	}))

	// Basic Auth
	// ----------
	cfg := basicauth.Config{
		Users: map[string]string{
			viper.GetString("SERVER_BASICAUTH_USERNAME"): viper.GetString("SERVER_BASICAUTH_PASSWORD"),
		},
	}

	// Prometheus
	// ----------
	if viper.GetBool("SERVER_PROMETHEUS") {
		prometheus := fiberprometheus.New("go-fiber")
		prometheus.RegisterAt(s, "/metrics")
		s.Use(prometheus.Middleware)
	}

	// Pprof
	// -----
	if viper.GetBool("SERVER_PPROF") {
		private := s.Group("/debug/pprof")
		private.Use(basicauth.New(cfg))
		s.Use(pprof.New())
	}

	// Monitor
	// -------
	// Consumes memory periodically
	if viper.GetBool("SERVER_MONITOR") {
		tools := s.Group("/tools")
		tools.Use(basicauth.New(cfg))
		tools.Get("/monitor", monitor.New())
	}
}

func initJWT(s *fiber.App) {
	s.Use(jwtware.New(jwtware.Config{
		SigningMethod: "HS512",
		SigningKey:    []byte(viper.GetString("JWT_SECRET")),
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			return c.Status(fiber.StatusUnauthorized).JSON(utils.HTTPError{
				Code:    fiber.StatusUnauthorized,
				Message: "Invalid or expired JWT",
			})
		},
	}))
}
