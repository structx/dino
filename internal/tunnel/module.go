package tunnel

import (
	"github.com/structx/dino/auth"
	"github.com/structx/dino/database"
	"github.com/structx/dino/gateway"
	pb "github.com/structx/dino/pb/tunnels/v1"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

// Params
type Params struct {
	fx.In

	Logger *zap.Logger

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
