package verifier

import (
	"go.uber.org/fx"
	"soft.structx.io/dino/internal/tunnel"
)

// Params
type Params struct {
	fx.In

	TunnelService tunnel.Service `name:"tunnel_service"`
}

type Result struct {
	fx.Out

	Verifier Verifier
}

// Module
var Module = fx.Module("verifier", fx.Provide(newModule))

func newModule(p Params) Result {
	return Result{
		Verifier: newVerifier(p.TunnelService),
	}
}
