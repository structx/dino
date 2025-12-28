package routes

import (
	"context"
	"time"

	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
	pb "soft.structx.io/dino/pb/routes/v1"
)

type routeServer struct {
	pb.UnimplementedRouteServiceServer

	log *zap.Logger
	svc Service
}

func newRouteServer(logger *zap.Logger, routeService Service) pb.RouteServiceServer {
	return &routeServer{
		log: logger.Named("route_server"),
		svc: routeService,
	}
}

// CreateRoute
func (rs *routeServer) CreateRoute(ctx context.Context, in *pb.CreateRouteRequest) (*pb.CreateRouteResponse, error) {
	rs.log.Debug("CreateRoute", zap.Any("request", in))
	args := RouteCreate{
		Tunnel:              in.Create.Tunnel,
		Hostname:            in.Create.Hostname,
		DestinationProtocol: in.Create.DestProtocol,
		DestinationIP:       in.Create.DestAddr,
		DestinationPort:     in.Create.DestPort,
	}

	route, err := rs.svc.Create(ctx, args)
	if err != nil {
		rs.log.Error("create route", zap.Error(err))
		return nil, status.Error(codes.Internal, codes.Internal.String())
	}

	return newCreateRouteResponse(route), nil
}

// DeleteRoute
func (rs *routeServer) DeleteRoute(ctx context.Context, in *pb.DeleteRouteRequest) (*pb.DeleteRouteResponse, error) {
	err := rs.svc.Delete(ctx, in.Hostname)
	if err != nil {
		rs.log.Error("delete route", zap.Error(err))
		return nil, status.Error(codes.Internal, codes.Internal.String())
	}
	return newDeleteRouteResponse(), nil
}

func newDeleteRouteResponse() *pb.DeleteRouteResponse {
	return &pb.DeleteRouteResponse{}
}

// GetRoute
func (rs *routeServer) GetRoute(ctx context.Context, in *pb.GetRouteRequest) (*pb.GetRouteResponse, error) {
	route, err := rs.svc.Get(ctx, in.Hostname)
	if err != nil {
		rs.log.Error("get route", zap.Error(err))
		return nil, status.Error(codes.Internal, codes.Internal.String())
	}
	return newGetRouteResponse(route), nil
}

// ListRoutes
func (rs *routeServer) ListRoutes(ctx context.Context, in *pb.ListRoutesRequest) (*pb.ListRoutesResponse, error) {
	partials, err := rs.svc.List(ctx, RouteList{
		Tunnel: in.Tunnel,
		Limit:  in.Limit,
		Offset: in.Offset,
	})
	if err != nil {
		rs.log.Error("list routes", zap.Error(err))
		return nil, status.Error(codes.Internal, codes.Internal.String())
	}
	return newListRoutesResponse(partials), err
}

// UpdateRoute
func (rs *routeServer) UpdateRoute(ctx context.Context, in *pb.UpdateRouteRequest) (*pb.UpdateRouteResponse, error) {

	args := RouteUpdate{
		ID:                  in.Update.Uid,
		Hostname:            in.Update.Hostname,
		DestinationProtocol: in.Update.DestProtocol,
		DestinationIP:       in.Update.DestAddr,
		DestinationPort:     in.Update.DestPort,
		Enabled:             in.Update.Enabled,
	}

	route, err := rs.svc.Update(ctx, args)
	if err != nil {
		rs.log.Error("update route", zap.Error(err))
		return nil, status.Error(codes.Internal, codes.Internal.String())
	}

	return newUpdateRouteResponse(route), nil
}

func pbRoute(r Route) *pb.Route {
	var updatedAt = time.Time{}
	if r.UpdatedAt != nil {
		updatedAt = *r.UpdatedAt
	}

	return &pb.Route{
		Uid:       r.ID,
		Tunnel:    r.Tunnel,
		Name:      r.Hostname,
		Enabled:   r.Enabled,
		CreatedAt: timestamppb.New(r.CreatedAt),
		UpdatedAt: timestamppb.New(updatedAt),
	}
}

func pbRoutePartial(p RoutePartial) *pb.RoutePartial {
	return &pb.RoutePartial{
		Uid:      p.ID,
		Hostname: p.Hostname,
	}
}

func newCreateRouteResponse(r Route) *pb.CreateRouteResponse {
	return &pb.CreateRouteResponse{
		Route: pbRoute(r),
	}
}

func newGetRouteResponse(r Route) *pb.GetRouteResponse {
	return &pb.GetRouteResponse{
		Route: pbRoute(r),
	}
}

func newListRoutesResponse(ps []RoutePartial) *pb.ListRoutesResponse {
	pbps := make([]*pb.RoutePartial, 0, len(ps))
	for _, p := range ps {
		pbps = append(pbps, pbRoutePartial(p))
	}
	return &pb.ListRoutesResponse{
		Partials: pbps,
	}
}

func newUpdateRouteResponse(r Route) *pb.UpdateRouteResponse {
	return &pb.UpdateRouteResponse{
		Route: pbRoute(r),
	}
}
