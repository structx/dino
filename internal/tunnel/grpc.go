package tunnel

import (
	"context"

	"github.com/structx/teapot"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"
	"soft.structx.io/dino/auth"
	pb "soft.structx.io/dino/pb/tunnels/v1"
)

type grpcServer struct {
	pb.UnimplementedTunnelServiceServer

	l *teapot.Logger
	s Service
	a auth.Authenticator
}

// interface compliance
var _ pb.TunnelServiceServer = (*grpcServer)(nil)

func newGrpcServer(logger *teapot.Logger, tunnelService Service, auth auth.Authenticator) pb.TunnelServiceServer {
	return &grpcServer{
		l: logger,
		s: tunnelService,
		a: auth,
	}
}

// CreateTunnel
func (g *grpcServer) CreateTunnel(ctx context.Context, in *pb.CreateTunnelRequest) (*pb.CreateTunnelResponse, error) {
	args := TunnelCreate{
		Name: in.GetTunnelName(),
	}

	tunnel, sharedSecret, err := g.s.Create(ctx, args)
	if err != nil {
		g.l.Error("create tunnel", teapot.Error(err))
		return nil, status.Error(codes.Internal, codes.Internal.String())
	}

	token, err := g.a.GenerateJWT(tunnel.Name, tunnel.ID, sharedSecret.Secret)
	if err != nil {
		g.l.Error("generate jwt", teapot.Error(err))
		return nil, status.Error(codes.Internal, codes.Internal.String())
	}

	return newCreateTunnelResponse(tunnel, token), nil
}

// DeleteTunnel
func (g *grpcServer) DeleteTunnel(ctx context.Context, in *pb.DeleteTunnelRequest) (*pb.DeleteTunnelResponse, error) {
	err := g.s.Delete(ctx, in.GetName())
	if err != nil {
		g.l.Error("delete tunnel", teapot.Error(err))
		return nil, status.Error(codes.Internal, codes.Internal.String())
	}
	return newDeleteTunnelResponse(), nil
}

// GetTunnel
func (g *grpcServer) GetTunnel(ctx context.Context, in *pb.GetTunnelRequest) (*pb.GetTunnelResponse, error) {
	tunnel, err := g.s.Get(ctx, in.GetName())
	if err != nil {
		g.l.Error("get tunnel", teapot.Error(err))
		return nil, status.Error(codes.Internal, codes.Internal.String())
	}
	return newGetTunnelReply(tunnel), nil
}

// ListTunnels
func (g *grpcServer) ListTunnels(ctx context.Context, in *pb.ListTunnelsRequest) (*pb.ListTunnelsResponse, error) {
	partials, err := g.s.List(ctx, in.Limit, in.Offset)
	if err != nil {
		g.l.Error("list tunnels", teapot.Error(err))
		return nil, status.Error(codes.Internal, codes.Internal.String())
	}
	return newListTunnelsReply(partials), nil
}

// UpdateTunnel
func (g *grpcServer) UpdateTunnel(ctx context.Context, in *pb.UpdateTunnelRequest) (*pb.UpdateTunnelResponse, error) {
	args := TunnelUpdate{
		OldName: in.TunnelUpdate.GetOldName(),
		Name:    in.TunnelUpdate.GetNewName(),
	}
	tunnel, err := g.s.Update(ctx, args)
	if err != nil {
		g.l.Error("update tunnel", teapot.Error(err))
		return nil, status.Error(codes.Internal, codes.Internal.String())
	}
	return newUpdateTunnelResponse(tunnel), nil
}

func newUpdateTunnelResponse(t Tunnel) *pb.UpdateTunnelResponse {
	return &pb.UpdateTunnelResponse{
		Tunnel: pbTunnel(t),
	}
}

func pbTunnel(t Tunnel) *pb.Tunnel {
	return &pb.Tunnel{
		Id:        t.ID,
		Name:      t.Name,
		CreatedAt: timestamppb.New(t.CreatedAt),
	}
}

func newCreateTunnelResponse(t Tunnel, sharedSecret string) *pb.CreateTunnelResponse {
	return &pb.CreateTunnelResponse{
		Tunnel: pbTunnel(t),
		AuthDetails: &pb.CreateTunnelResponse_SecretKey{
			SecretKey: sharedSecret,
		},
	}
}

func newGetTunnelReply(t Tunnel) *pb.GetTunnelResponse {
	return &pb.GetTunnelResponse{
		Tunnel: pbTunnel(t),
	}
}

func newListTunnelsReply(p []TunnelPartial) *pb.ListTunnelsResponse {
	ps := make([]*pb.TunnelPartial, 0, len(p))
	for _, t := range p {
		ps = append(ps, &pb.TunnelPartial{Name: t.Name})
	}
	return &pb.ListTunnelsResponse{
		Tunnels: ps,
	}
}

func newDeleteTunnelResponse() *pb.DeleteTunnelResponse {
	return &pb.DeleteTunnelResponse{
		Empty: &emptypb.Empty{},
	}
}
