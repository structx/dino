package sessions

import (
	"github.com/structx/teapot"
	"go.uber.org/fx"
)

// Params
type Params struct {
	fx.In

	Lc fx.Lifecycle

	Logger *teapot.Logger
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
