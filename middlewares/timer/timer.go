package timer

import (
	"fmt"
	"time"

	"github.com/gofiber/fiber/v2"
)

// Config defines the configuration for middleware.
type Config struct {
	// DisplaySeconds indicates the process time in seconds.
	//
	// Optional. Default value true.
	DisplaySeconds bool

	// DisplayMilliseconds indicates the process time in milliseconds.
	//
	// Optional. Default value false.
	DisplayMilliseconds bool

	// Prefix indicates prefix for header name.
	//
	// Optional. Default value "X-Process-Time".
	Prefix string
}

// ConfigDefault is the default configuration.
var ConfigDefault = Config{
	DisplaySeconds:      true,
	DisplayMilliseconds: false,
	Prefix:              "x-process-time",
}

// New creates a new instance of middleware handler.
func New(config ...Config) func(*fiber.Ctx) error {
	// Default configuration
	cfg := ConfigDefault

	// Override configuration if provided
	if len(config) > 0 {
		cfg = config[0]

		if cfg.Prefix == "" {
			cfg.Prefix = ConfigDefault.Prefix
		}
	}

	return func(c *fiber.Ctx) error {
		now := time.Now()

		err := c.Next()

		duration := time.Since(now)
		if cfg.DisplayMilliseconds {
			c.Set(cfg.Prefix+"-ms", fmt.Sprintf("%d", duration.Milliseconds()))
		}
		if cfg.DisplaySeconds {
			c.Set(cfg.Prefix+"-sec", fmt.Sprintf("%.6f", duration.Seconds()))
		}

		return err
	}
}
