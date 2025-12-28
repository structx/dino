package sessions

import (
	"go.uber.org/fx"
	"go.uber.org/zap"
)

// Params
type Params struct {
	fx.In

	Lc fx.Lifecycle

	Logger *zap.Logger
}

// Result
type Result struct {
	fx.Out

	SessionMultiplexer Mux
}

// Module
var Module = fx.Module("tunnel_sessions", fx.Provide(newModule))

func newModule(p Params) Result {
	sessionMux := newMux(p.Logger)
	p.Lc.Append(fx.Hook{
		OnStart: sessionMux.start,
		OnStop:  sessionMux.stop,
	})

	return Result{
		SessionMultiplexer: sessionMux,
	}
}
