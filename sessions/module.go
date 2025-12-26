package sessions

import (
	"github.com/structx/dino/internal/tunnel"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

// Params
type Params struct {
	fx.In

	Logger *zap.Logger

	TunnelService tunnel.Service `name:"tunnel_service"`
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
		Mux: newMux(p.Logger, p.TunnelService),
	}
}
