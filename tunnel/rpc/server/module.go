package server

import (
	"go.uber.org/fx"
	"go.uber.org/zap"
	pb "soft.structx.io/dino/pb/rtunnel/v1"
	"soft.structx.io/dino/sessions"
	"soft.structx.io/dino/tunnel/gateway"
	"soft.structx.io/dino/tunnel/verifier"
)

// Paramss
type Params struct {
	fx.In

	Logger *zap.Logger

	Mux      sessions.Multiplexer
	Verifier verifier.Verifier
}

// Result
type Result struct {
	fx.Out

	Transport *gateway.TunnelTransport
}

// Module
var Module = fx.Module("rtunnel_rpc", fx.Provide(newModule))

func newModule(p Params) Result {
	rts := newReverseTunnelServer(p.Logger, p.Mux, p.Verifier)
	return Result{
		Transport: &gateway.TunnelTransport{
			ServiceDesc: &pb.ReverseTunnelService_ServiceDesc,
			Service:     rts,
		},
	}
}
