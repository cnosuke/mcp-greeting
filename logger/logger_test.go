package logger

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

func TestInitLogger_InvalidLevel(t *testing.T) {
	err := InitLogger("invalid", "", false)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "invalid log level")
}

func TestInitLogger_SuppressConsoleWithNoPath(t *testing.T) {
	err := InitLogger("info", "", true)
	assert.NoError(t, err)

	// Logger should work without panicking (NopCore)
	zap.S().Info("test message")
}

func TestInitLogger_SuppressConsoleWithStdoutPath(t *testing.T) {
	err := InitLogger("info", "stdout", true)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "conflicts with stdio mode")
}

func TestInitLogger_SuppressConsoleWithStderrPath(t *testing.T) {
	err := InitLogger("info", "stderr", true)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "conflicts with stdio mode")
}

func TestInitLogger_ValidLevels(t *testing.T) {
	levels := []string{"debug", "info", "warn", "error"}
	for _, level := range levels {
		t.Run(level, func(t *testing.T) {
			err := InitLogger(level, "", true)
			assert.NoError(t, err)
		})
	}
}

func TestInitLogger_FileLogging(t *testing.T) {
	logFile := filepath.Join(t.TempDir(), "test.log")

	err := InitLogger("info", logFile, true)
	assert.NoError(t, err)

	zap.S().Infow("file logging test", "key", "value")
	_ = Sync()

	data, err := os.ReadFile(logFile)
	assert.NoError(t, err)
	assert.Contains(t, string(data), "file logging test")
}
