package interceptors

import (
	"context"

	"go.uber.org/zap"
	"google.golang.org/grpc"
	"soft.structx.io/dino/gateway"
)

type loggerInterceptor struct {
	logger *zap.Logger
}

func newLoggerInterceptor(logger *zap.Logger) *loggerInterceptor {
	return &loggerInterceptor{
		logger: logger.Named("logging_interceptor"),
	}
}

// UnaryInterceptor
func (li *loggerInterceptor) UnaryInterceptor() gateway.UnaryInterceptor {
	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp any, err error) {
		li.logger.Debug("gRPC handler", zap.String("full_method", info.FullMethod))
		resp, err = handler(ctx, req)
		if err != nil {
			return nil, err
		}
		return resp, nil
	}
}

func (li *loggerInterceptor) StreamInterceptor() gateway.StreamInterceptor {
	return func(srv any, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		return nil
	}
}
