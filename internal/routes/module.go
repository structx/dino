package routes

import (
	"go.uber.org/fx"
	"go.uber.org/zap"
	"soft.structx.io/dino/database"
	"soft.structx.io/dino/gateway"
	pb "soft.structx.io/dino/pb/routes/v1"
	"soft.structx.io/dino/pubsub"
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

	RouteService Service

	Transport gateway.Transport `group:"transport"`
}

// Module
var Module = fx.Module("routes_module", fx.Provide(newModule))

func newModule(p Params) Result {
	svc := newService(p.DBTX, p.Broker)
	rs := newRouteServer(p.Logger, svc)
	return Result{
		RouteService: svc,
		Transport: gateway.Transport{
			ServiceDesc: &pb.RouteService_ServiceDesc,
			Service:     rs,
		},
	}
}
