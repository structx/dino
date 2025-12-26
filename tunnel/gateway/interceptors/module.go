package interceptors

import (
	"github.com/structx/dino/tunnel/gateway"
	"github.com/structx/dino/tunnel/verifier"
	"go.uber.org/fx"
	"go.uber.org/zap"
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
