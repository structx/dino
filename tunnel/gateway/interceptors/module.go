package interceptors

import (
	"go.uber.org/fx"
	"go.uber.org/zap"
	"soft.structx.io/dino/tunnel/gateway"
	"soft.structx.io/dino/tunnel/verifier"
)

// Params
type Params struct {
	fx.In

	Logger *zap.Logger

	Verifier verifier.Verifier
}

// Result
type Result struct {
	fx.Out

	Interceptor gateway.StreamServerInterceptor
}

// Module
var Module = fx.Module("tunnel_gateway_interceptors", fx.Provide(newModule))

func newModule(p Params) Result {
	return Result{
		Interceptor: newServerInterceptor(p.Logger, p.Verifier),
	}
}
