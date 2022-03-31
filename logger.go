package server

import (
	"errors"
	"fmt"
	"os"
	"path"
	"time"

	"github.com/fabienbellanger/goutils"
	"github.com/gofiber/fiber/v2"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// InitLogger initializes custom Zap logger configuration.
func InitLogger() (*zap.Logger, error) {
	// Logs outputs
	// ------------
	outputs, err := getLoggerOutputs(viper.GetStringSlice("LOG_OUTPUTS"), viper.GetString("APP_NAME"), viper.GetString("LOG_PATH"))
	if err != nil {
		return nil, err
	}

	// Level
	// -----
	level := getLoggerLevel(viper.GetString("LOG_LEVEL"), viper.GetString("APP_ENV"))

	cfg := zap.Config{
		Encoding:         "json",
		Level:            zap.NewAtomicLevelAt(level),
		OutputPaths:      outputs,
		ErrorOutputPaths: outputs,
		EncoderConfig: zapcore.EncoderConfig{
			MessageKey:   "message",
			LevelKey:     "level",
			EncodeLevel:  zapcore.CapitalLevelEncoder,
			TimeKey:      "time",
			EncodeTime:   zapcore.RFC3339TimeEncoder,
			CallerKey:    "caller",
			EncodeCaller: zapcore.ShortCallerEncoder,
		},
	}
	logger, err := cfg.Build()
	if err != nil {
		return zap.NewProduction()
	}

	return logger, nil
}

// getLoggerOutputs returns an array with the log outputs.
// Outputs can be stdout and/or file.
func getLoggerOutputs(logOutputs []string, appName, filePath string) (outputs []string, err error) {
	if goutils.StringInSlice("file", logOutputs) {
		logPath := path.Clean(filePath)
		_, err := os.Stat(logPath)
		if err != nil {
			return nil, err
		}

		if appName == "" {
			return nil, errors.New("no APP_NAME variable defined")
		}

		outputs = append(outputs, fmt.Sprintf("%s/%s.log",
			logPath,
			appName))
	}
	if goutils.StringInSlice("stdout", logOutputs) {
		outputs = append(outputs, "stdout")
	}
	return
}

// getLoggerLevel returns the minimum log level.
// If nothing is specified in the environment variable LOG_LEVEL,
// The level is DEBUG in development mode and WARN in others cases.
func getLoggerLevel(l string, env string) (level zapcore.Level) {
	switch l {
	case "debug":
		level = zapcore.DebugLevel
	case "info":
		level = zapcore.InfoLevel
	case "warn":
		level = zapcore.WarnLevel
	case "error":
		level = zapcore.ErrorLevel
	case "panic":
		level = zapcore.PanicLevel
	case "fatal":
		level = zapcore.FatalLevel
	default:
		if env == "development" {
			level = zapcore.DebugLevel
		} else {
			level = zapcore.WarnLevel
		}
	}
	return
}

// zapLogger is a middleware and zap to provide an "access log" like logging for each request.
func zapLogger(log *zap.Logger) fiber.Handler {
	return func(c *fiber.Ctx) error {
		start := time.Now().UTC()

		err := c.Next()
		if err != nil {
			return err
		}

		stop := time.Since(start)
		req := c.Request()
		fields := []zapcore.Field{
			zap.Int("code", c.Response().StatusCode()),
			zap.String("method", string(req.Header.Method())),
			zap.String("path", string(req.RequestURI())),
			zap.String("url", string(req.URI().FullURI())),
			zap.String("ip", c.IP()),
			zap.String("userAgent", c.Get(fiber.HeaderUserAgent)),
			zap.String("latency", stop.String()),
			zap.String("requestId", fmt.Sprintf("%s", c.Locals("requestid"))),
		}
		go func() {
			log.Info("", fields...)
		}()

		return nil
	}
}
