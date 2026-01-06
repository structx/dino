package sessions

import (
	"github.com/structx/teapot"
	"go.uber.org/fx"
	"soft.structx.io/dino/internal/routes"
	"soft.structx.io/dino/internal/tunnel"
	"soft.structx.io/dino/pubsub"
)

// Params
type Params struct {
	fx.In

	Lc fx.Lifecycle

	Logger *teapot.Logger

	TunnelService tunnel.Service `name:"tunnel_service"`
	RouteService  routes.Service

	Broker pubsub.Broker
}

// Result
type Result struct {
	fx.Out

	Mux Multiplexer
}

// Module
var Module = fx.Module("sessions", fx.Provide(newModule))

func newModule(p Params) Result {
	mux := newMux(p.Logger, p.Broker, p.TunnelService, p.RouteService)

	p.Lc.Append(fx.Hook{
		OnStart: mux.start,
		OnStop:  mux.stop,
	})

	return Result{
		Mux: mux,
	}
}
