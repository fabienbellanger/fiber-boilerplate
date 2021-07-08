package server

import (
	"errors"
	"fmt"
	"os"
	"path"

	"github.com/fabienbellanger/goutils"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// InitLogger initializes custom Zap logger.
func InitLogger() (*zap.Logger, error) {
	// Logs outputs
	// ------------
	logOutputs := viper.GetStringSlice("LOG_OUTPUTS")
	var outputs []string
	if goutils.StringInSlice("file", logOutputs) {
		logPath := path.Clean(viper.GetString("LOG_PATH"))
		_, err := os.Stat(logPath)
		if err != nil {
			return nil, err
		}

		appName := viper.GetString("APP_NAME")
		if appName == "" {
			return nil, errors.New("no APP_NAME variable defined")
		}

		outputs = append(outputs, fmt.Sprintf("%s/%s.log",
			logPath,
			appName))
	}
	if goutils.StringInSlice("stderr", logOutputs) {
		outputs = append(outputs, "stderr")
	}
	if goutils.StringInSlice("stdout", logOutputs) {
		outputs = append(outputs, "stdout")
	}

	// Level
	// -----
	level := getLoggerLevel(viper.GetString("LOG_LEVEL"))

	cfg := zap.Config{
		Encoding:         "console",
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

func getLoggerOutputs(outputs []string, appName, path string) (logOutputs []string) {
	// logOutputs := viper.GetStringSlice("LOG_OUTPUTS")
	// var outputs []string
	// if goutils.StringInSlice("file", logOutputs) {
	// 	logPath := path.Clean(viper.GetString("LOG_PATH"))
	// 	_, err := os.Stat(logPath)
	// 	if err != nil {
	// 		return nil, err
	// 	}

	// 	appName := viper.GetString("APP_NAME")
	// 	if appName == "" {
	// 		return nil, errors.New("no APP_NAME variable defined")
	// 	}

	// 	outputs = append(outputs, fmt.Sprintf("%s/%s.log",
	// 		logPath,
	// 		appName))
	// }
	// if goutils.StringInSlice("stderr", logOutputs) {
	// 	outputs = append(outputs, "stderr")
	// }
	// if goutils.StringInSlice("stdout", logOutputs) {
	// 	outputs = append(outputs, "stdout")
	// }
	return
}

func getLoggerLevel(l string) (level zapcore.Level) {
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
		level = zapcore.InfoLevel
	}
	return
}
