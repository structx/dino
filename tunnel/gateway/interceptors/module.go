package interceptors

import (
	"github.com/structx/teapot"
	"go.uber.org/fx"
	"soft.structx.io/dino/tunnel/gateway"
	"soft.structx.io/dino/tunnel/verifier"
)

// Params
type Params struct {
	fx.In

	Logger *teapot.Logger

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
