package interceptors

import (
	"github.com/structx/teapot"
	"go.uber.org/fx"
	"soft.structx.io/dino/gateway"
)

// Params
type Params struct {
	fx.In

	Logger *teapot.Logger
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
