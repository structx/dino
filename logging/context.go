package logging

import (
	"context"

	"go.uber.org/zap"
)

type contextKey string

const consoleLoggerKey contextKey = "zap_console_logger"

// WithContext
func WithContext(ctx context.Context, logger *zap.Logger) context.Context {
	return context.WithValue(ctx, consoleLoggerKey, logger)
}

// FromContext
func FromContext(ctx context.Context) *zap.Logger {
	logger, ok := ctx.Value(consoleLoggerKey).(*zap.Logger)
	if !ok {
		return nil
	}
	return logger
}
