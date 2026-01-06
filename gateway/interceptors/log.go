package interceptors

import (
	"context"

	"github.com/structx/teapot"
	"google.golang.org/grpc"
	"soft.structx.io/dino/gateway"
)

type loggerInterceptor struct {
	l *teapot.Logger
}

func newLoggerInterceptor(logger *teapot.Logger) *loggerInterceptor {
	return &loggerInterceptor{
		l: logger,
	}
}

// UnaryInterceptor
func (li *loggerInterceptor) UnaryInterceptor() gateway.UnaryInterceptor {
	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp any, err error) {
		li.l.Debug("gRPC handler", teapot.String("full_method", info.FullMethod))
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
