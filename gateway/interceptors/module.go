package interceptors

import (
	"go.uber.org/fx"
	"go.uber.org/zap"
	"soft.structx.io/dino/gateway"
)

// Params
type Params struct {
	fx.In

	Logger *zap.Logger
}

// Result
type Result struct {
	fx.Out

	UnaryInterceptors []gateway.UnaryInterceptor
	StreamInterceptor []gateway.StreamInterceptor
}

// Module
var Module = fx.Module("interceptors", fx.Provide(newModule))

func newModule(p Params) Result {
	logger := newLoggerInterceptor(p.Logger)
	_ = logger
	return Result{
		UnaryInterceptors: []gateway.UnaryInterceptor{},
		StreamInterceptor: []gateway.StreamInterceptor{},
	}
}
