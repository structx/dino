package server

import (
	pb "github.com/structx/dino/pb/rtunnel/v1"
	"github.com/structx/dino/sessions"
	"github.com/structx/dino/tunnel/gateway"
	"github.com/structx/dino/tunnel/verifier"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

// Paramss
type Params struct {
	fx.In

	Logger *zap.Logger

	Mux      sessions.Muxer
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
