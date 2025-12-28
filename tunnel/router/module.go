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
	mux := newMux()
	mux.Add("whoami.dino.local", "http", "127.0.0.1", "8888")
	return Result{
		Mux: mux,
	}
}
