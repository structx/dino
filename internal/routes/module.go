package routes

import (
	"github.com/structx/dino/database"
	"github.com/structx/dino/gateway"
	pb "github.com/structx/dino/pb/routes/v1"
	"github.com/structx/dino/pubsub"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

// Params
type Params struct {
	fx.In

	Logger *zap.Logger

	DBTX database.DBTX

	Broker pubsub.Broker
}

// Result
type Result struct {
	fx.Out

	Transport gateway.Transport `group:"transport"`
}

// Module
var Module = fx.Module("routes_module", fx.Provide(newModule))

func newModule(p Params) Result {
	svc := newService(p.DBTX, p.Broker)
	rs := newRouteServer(p.Logger, svc)
	return Result{
		Transport: gateway.Transport{
			ServiceDesc: &pb.RouteService_ServiceDesc,
			Service:     rs,
		},
	}
}
