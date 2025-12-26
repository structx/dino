package verifier

import (
	"github.com/structx/dino/internal/tunnel"
	"go.uber.org/fx"
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
