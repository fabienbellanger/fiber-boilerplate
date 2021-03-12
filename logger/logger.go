package logger

import (
	"fmt"

	"github.com/spf13/viper"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// TODO: Pour rappel
// _, file, line, ok := runtime.Caller(1)
// if ok {
// file = goutils.SubPath(file, viper.GetString("APP_NAME"))
// }
// log.Printf("file=%v, line=%v\n", file, line)

// Init initializes custom Zap logger.
func Init() (*zap.Logger, error) {
	cfg := zap.Config{
		Encoding:         "json",
		Level:            zap.NewAtomicLevelAt(zapcore.DebugLevel),
		OutputPaths:      []string{"stderr", fmt.Sprintf("/tmp/%s.log", viper.GetString("APP_NAME"))},
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