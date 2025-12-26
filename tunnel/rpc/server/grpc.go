package server

import (
	"errors"
	"fmt"
	"io"

	pb "github.com/structx/dino/pb/rtunnel/v1"
	"github.com/structx/dino/sessions"
	tunnelnet "github.com/structx/dino/tunnel/net"
	"github.com/structx/dino/tunnel/verifier"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

type tunnelConn struct {
	str grpc.BidiStreamingServer[pb.TunnelMessage, pb.TunnelMessage]
}

// Read implements net.Conn.
func (t tunnelConn) Read() (*tunnelnet.DataFrame, error) {
	msg, err := t.str.Recv()
	if err != nil {
		if err == io.EOF {
			return nil, io.EOF
		}
		return nil, fmt.Errorf("str.Recv: %w", err)
	}

	if msg.GetCloseConnection() != nil {
		return nil, io.EOF
	} else if msg.GetData() != nil {
		return &tunnelnet.DataFrame{SessionID: msg.GetSessionId(), Payload: msg.GetData()}, nil
	}

	return nil, errors.New("unsupported message")
}

// Write implements net.Conn.
func (t tunnelConn) Write(df *tunnelnet.DataFrame) (int, error) {

	if df.RouteUpdate != nil {
		ru := df.RouteUpdate
		if err := t.str.Send(&pb.TunnelMessage{
			SessionId: "",
			Payload: &pb.TunnelMessage_RouteUpdates{
				RouteUpdates: &pb.Route{
					Hostname:            ru.Hostname,
					DestinationProtocol: ru.DestProtocol,
					DestinationIp:       ru.DestIP,
					DestinationPort:     ru.DestPort,
					IsDeleted:           ru.IsDelete,
				},
			},
		}); err != nil {
			return 0, fmt.Errorf("str.Send: %w", err)
		}

		return 0, nil
	}

	if err := t.str.Send(&pb.TunnelMessage{
		SessionId: df.SessionID,
		Payload: &pb.TunnelMessage_Data{
			Data: df.Payload,
		},
	}); err != nil {
		return 0, fmt.Errorf("str.Send: %w", err)
	}
	return len(df.Payload), nil
}

// Close implements net.Conn.
func (t tunnelConn) Close(sessionID string) error {
	return t.str.Send(&pb.TunnelMessage{
		SessionId: sessionID,
		Payload: &pb.TunnelMessage_CloseConnection{
			CloseConnection: &pb.CloseConnection{},
		},
	})
}

type reverseTunnelServer struct {
	pb.UnimplementedReverseTunnelServiceServer

	log *zap.Logger

	mux      sessions.Muxer
	verifier verifier.Verifier
}

// interface compliance
var _ pb.ReverseTunnelServiceServer = (*reverseTunnelServer)(nil)
var _ tunnelnet.Conn = (*tunnelConn)(nil)

func newReverseTunnelServer(logger *zap.Logger, sessionMux sessions.Muxer, verifier verifier.Verifier) pb.ReverseTunnelServiceServer {
	return &reverseTunnelServer{
		log:      logger.Named("rtunnel_server"),
		mux:      sessionMux,
		verifier: verifier,
	}
}

// EstablishTunnel
func (rts *reverseTunnelServer) EstablishTunnel(stream grpc.BidiStreamingServer[pb.TunnelMessage, pb.TunnelMessage]) error {
	rts.log.Info("establish tunnel")
	ctx := stream.Context()

	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		rts.log.Debug("missing request metadata")
		return status.Error(codes.InvalidArgument, codes.InvalidArgument.String())
	}

	authorizations := md.Get("authorization")
	if len(authorizations) < 1 {
		rts.log.Debug("missing authorization metadata")
		return status.Error(codes.InvalidArgument, codes.InvalidArgument.String())
	}

	ids := md.Get("tunnel-id")
	if len(ids) < 1 {
		rts.log.Debug("missing tunnel-id metadata")
		return status.Error(codes.InvalidArgument, codes.InvalidArgument.String())
	}

	claims, err := rts.verifier.VerifyToken(ctx, ids[0], authorizations[0])
	if err != nil {
		rts.log.Error("failed to verify token", zap.Error(err))
		return status.Error(codes.Unauthenticated, codes.Unauthenticated.String())
	}

	tc := tunnelConn{str: stream}
	if err := rts.mux.RegisterTunnel(ctx, tc, claims.ID); err != nil {
		rts.log.Error("sessionManager.RegisterTunnel", zap.Error(err))
		return status.Error(codes.Internal, codes.Internal.String())
	}
	// defer rts.sessionManager.

	<-ctx.Done()
	return nil
}
