package logging

import (
	"os"
	"strings"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// NewConsoleLogger
func NewConsoleLogger(lvl string) *zap.Logger {
	encoderConfig := zapcore.EncoderConfig{
		EncodeLevel: zapcore.CapitalColorLevelEncoder,
	}

	var l zapcore.Level
	switch strings.ToLower(lvl) {
	case "error":
		l = zapcore.ErrorLevel
	case "info":
		l = zapcore.InfoLevel
	default:
		l = zapcore.DebugLevel
	}

	core := zapcore.NewCore(
		zapcore.NewConsoleEncoder(encoderConfig),
		zapcore.AddSync(os.Stdout),
		l,
	)

	logger := zap.New(core)
	return logger
}
