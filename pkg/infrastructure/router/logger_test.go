package router

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/zap/zapcore"
)

func TestLogLevel(t *testing.T) {
	cases := map[string]zapcore.Level{
		"debug": zapcore.DebugLevel,
		"info":  zapcore.InfoLevel,
		"warn":  zapcore.WarnLevel,
		"error": zapcore.ErrorLevel,
		"panic": zapcore.PanicLevel,
		"fatal": zapcore.FatalLevel,
		"":      zapcore.DebugLevel,
	}

	env := "development"
	for level, expected := range cases {
		assert.Equal(t, expected, getLoggerLevel(level, env))
	}

	env = "production"
	cases[""] = zapcore.WarnLevel
	for level, expected := range cases {
		assert.Equal(t, expected, getLoggerLevel(level, env))
	}
}

func TestLogOutputsWithOneOutput(t *testing.T) {
	appName := "go-url-shortener"
	filePath := "/tmp"
	outputs := []string{"stdout"}

	gottenOutputs, err := getLoggerOutputs(outputs, "", "")
	assert.Equal(t, []string{"stdout"}, gottenOutputs, "with stdout")
	assert.Nil(t, err)

	outputs = []string{"file"}
	gottenOutputs, err = getLoggerOutputs(outputs, appName, filePath)
	assert.Equal(t, []string{"/tmp/go-url-shortener.log"}, gottenOutputs, "with file")
	assert.Nil(t, err)

	gottenOutputs, err = getLoggerOutputs(outputs, "", filePath)
	assert.Equal(t, []string(nil), gottenOutputs, "with file and empty app name")
	assert.NotNil(t, err)

	gottenOutputs, err = getLoggerOutputs(outputs, appName, "")
	assert.Equal(t, []string{"./go-url-shortener.log"}, gottenOutputs, "with file and empty file path")
	assert.Nil(t, err)
}

func TestLogOutputsWithMoreThanOneOutput(t *testing.T) {
	appName := "go-url-shortener"
	filePath := "/tmp"
	outputs := []string{"stdout"}

	gottenOutputs, err := getLoggerOutputs(outputs, "", "")
	assert.Equal(t, []string{"stdout"}, gottenOutputs, "with stdout")
	assert.Nil(t, err)

	outputs = []string{"stdout", "file"}
	gottenOutputs, err = getLoggerOutputs(outputs, appName, filePath)
	assert.Equal(t, []string{"/tmp/go-url-shortener.log", "stdout"}, gottenOutputs, "with stdout")
	assert.Nil(t, err)
}
