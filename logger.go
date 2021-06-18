package server

import (
	"fmt"
	"os"
	"path"

	"github.com/spf13/viper"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// InitLogger initializes custom Zap logger.
func InitLogger() (*zap.Logger, error) {
	logPath := path.Clean(viper.GetString("LOG_PATH"))
	_, err := os.Stat(logPath)
	if err != nil {
		return nil, err
	}

	cfg := zap.Config{
		Encoding: "json",
		Level:    zap.NewAtomicLevelAt(zapcore.DebugLevel),
		OutputPaths: []string{
			"stderr",
			fmt.Sprintf("%s/%s.log",
				logPath,
				viper.GetString("APP_NAME"))},
		ErrorOutputPaths: []string{"stderr"},
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
