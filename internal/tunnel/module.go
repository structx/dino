package tunnel

import (
	"github.com/structx/teapot"
	"go.uber.org/fx"
	"soft.structx.io/dino/auth"
	"soft.structx.io/dino/database"
	"soft.structx.io/dino/gateway"
	pb "soft.structx.io/dino/pb/tunnels/v1"
)

// Params
type Params struct {
	fx.In

	Logger *teapot.Logger

	DB database.DBTX

	Auth auth.Authenticator
}

// Result
type Result struct {
	fx.Out

	Transport gateway.Transport `group:"transport"`

	TunnelService Service `name:"tunnel_service"`
}

// Module
var Module = fx.Module("tunnel", fx.Provide(newModule))

func newModule(p Params) Result {
	svc := newService(p.DB)
	s := newGrpcServer(p.Logger, svc, p.Auth)
	return Result{
		Transport: gateway.Transport{
			ServiceDesc: &pb.TunnelService_ServiceDesc,
			Service:     s,
		},
		TunnelService: svc,
	}
}
