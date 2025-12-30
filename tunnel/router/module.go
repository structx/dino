package router

import "go.uber.org/fx"

// Result
type Result struct {
	fx.Out

	Mux Mux
}

// Module
var Module = fx.Module("tunnel_router", fx.Provide(newModule))

func newModule() Result {
	return Result{
		Mux: newMux(),
	}
}
