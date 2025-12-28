package sessions

import (
	"go.uber.org/fx"
	"go.uber.org/zap"
	"soft.structx.io/dino/internal/routes"
	"soft.structx.io/dino/internal/tunnel"
	"soft.structx.io/dino/pubsub"
)

// Params
type Params struct {
	fx.In

	Logger *zap.Logger

	TunnelService tunnel.Service `name:"tunnel_service"`
	RouteService  routes.Service

	Broker pubsub.Broker
}

// Result
type Result struct {
	fx.Out

	Mux Muxer
}

// Module
var Module = fx.Module("sessions", fx.Provide(newModule))

func newModule(p Params) Result {
	return Result{
		Mux: newMux(p.Logger, p.Broker, p.TunnelService, p.RouteService),
	}
}
