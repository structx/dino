package client

import (
	"context"
	"crypto/tls"
	"fmt"
	"net"

	"go.uber.org/fx"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	pb "soft.structx.io/dino/pb/rtunnel/v1"
	"soft.structx.io/dino/setup"
	tunnelnet "soft.structx.io/dino/tunnel/net"
	"soft.structx.io/dino/tunnel/router"
	"soft.structx.io/dino/tunnel/sessions"
	"soft.structx.io/dino/tunnel/transport"
)

// Params
type Params struct {
	fx.In

	Lc fx.Lifecycle

	Logger *zap.Logger

	Cfg *setup.Tunnel

	SessionsManager sessions.Mux
	Mux             router.Mux
}

type Result struct {
	fx.Out
}

// Tunneler
type Tunneler interface {
	establishTunnel(context.Context) error
}

type tunnelClient struct {
	log *zap.Logger

	target   string
	tunnelID string
	token    string

	conn     *grpc.ClientConn
	sessions sessions.Mux
	mux      router.Mux
}

// interface compliance
var _ Tunneler = (*tunnelClient)(nil)

// Module
var Module = fx.Module("tunnel_client", fx.Invoke(newModule))

func newModule(p Params) error {

	tlsConfig := &tls.Config{
		InsecureSkipVerify: true,
		NextProtos:         []string{"h3"},
	}

	dial := transport.NewQuicDialer(tlsConfig)
	creds := transport.NewCredentials(tlsConfig)

	opts := []grpc.DialOption{
		// grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithTransportCredentials(creds),
		grpc.WithContextDialer(dial),
	}

	conn, err := grpc.NewClient("tunnel.dino.local:4242", opts...)
	if err != nil {
		return fmt.Errorf("grpc.NewClient: %w", err)
	}

	tc := &tunnelClient{
		log:      p.Logger,
		sessions: p.SessionsManager,
		mux:      p.Mux,
		target:   p.Cfg.Endpoint,
		tunnelID: p.Cfg.ID,
		token:    p.Cfg.Token,
		conn:     conn,
	}

	p.Lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			return tc.establishTunnel(ctx)
		},
		OnStop: func(ctx context.Context) error {
			return nil
		},
	})

	return nil
}

// EstablishTunnel implements Tunneler.
func (t *tunnelClient) establishTunnel(ctx context.Context) error {
	cli := pb.NewReverseTunnelServiceClient(t.conn)

	md := metadata.New(map[string]string{
		"tunnel-id":     t.tunnelID,
		"authorization": t.token,
	})
	ctx = metadata.NewOutgoingContext(ctx, md)

	ctx = context.WithoutCancel(ctx)
	stream, err := cli.EstablishTunnel(ctx)
	if err != nil {
		return fmt.Errorf("cli.EstablishTunnel: %w", err)
	}
	cc := &clientConn{str: stream}
	go t.worker(cc)

	return nil
}

func (t *tunnelClient) worker(conn tunnelnet.Conn) {
	for {
		df, err := conn.Read()
		if err != nil {
			t.log.Error("read from conn", zap.Error(err))
			continue
		}

		t.log.Debug("received msg", zap.Any("dataframe", df))

		if df.NewConn != nil {
			r, ok := t.mux.Match(df.NewConn.Hostname)
			if !ok {
				// send back missing route error
				t.log.Error("no matched route", zap.String("hostname", df.NewConn.Hostname))
				continue
			}

			hostAndPort := net.JoinHostPort(r.IP, r.Port)
			if err := t.sessions.InitSession(conn, df.SessionID, r.Protocol, hostAndPort); err != nil {
				t.log.Error("init session", zap.Error(err))
				continue
			}
		}
		if df.CloseConn != nil {
			if err := t.sessions.CloseSession(df.SessionID); err != nil {
				t.log.Error("close session", zap.Error(err))
				continue
			}
		}
		if df.RouteUpdate != nil {

			if df.RouteUpdate.IsDelete {
				t.mux.Del(df.RouteUpdate.Hostname)
				continue
			}

			port, err := validateAndParsePort(df.RouteUpdate.DestPort)
			if err != nil {
				t.log.Error("invalid port", zap.Error(err))
				continue
			}

			t.log.Info("config update", zap.Any("route", df.RouteUpdate))
			t.mux.Add(df.RouteUpdate.Hostname, df.RouteUpdate.DestProtocol, df.RouteUpdate.DestIP, port)
			continue
		}

		err = t.sessions.RouteMsg(df.SessionID, df.Payload)
		if err != nil {
			t.log.Error("route message", zap.Error(err))
		}
	}
}
