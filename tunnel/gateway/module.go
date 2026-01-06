package gateway

import (
	"context"
	"crypto/tls"
	"fmt"
	"net"

	"github.com/quic-go/quic-go"
	"github.com/quic-go/quic-go/qlog"
	"github.com/structx/teapot"
	"go.uber.org/fx"
	"go.uber.org/multierr"
	"google.golang.org/grpc"
	"soft.structx.io/dino/setup"
	"soft.structx.io/dino/tunnel/transport"
)

// StreamServerInterceptor
type StreamServerInterceptor = grpc.StreamServerInterceptor

// TunnelTransport
type TunnelTransport struct {
	Service     any
	ServiceDesc *grpc.ServiceDesc
}

// Params
type Params struct {
	fx.In

	Lc fx.Lifecycle

	Logger *teapot.Logger

	Cfg *setup.Server

	Transport   *TunnelTransport
	Interceptor StreamServerInterceptor
}

// Module
var Module = fx.Module("tunnel_gateway", fx.Invoke(invokeModule))

func invokeModule(p Params) error {
	tlsConfig, err := newTlsConfig(p.Cfg.CertPath, p.Cfg.KeyPath)
	if err != nil {
		return fmt.Errorf("newTlsConfig: %w", err)
	}

	// lis, err := net.Listen("tcp", ":4242")
	// if err != nil {
	// 	return fmt.Errorf("net.Listen: %w", err)
	// }

	s := grpc.NewServer()
	s.RegisterService(p.Transport.ServiceDesc, p.Transport.Service)

	// grpclog.SetLoggerV2(zapgrpc.NewLogger(p.Logger))

	quicHostAndPort := net.JoinHostPort(p.Cfg.QuicHost, p.Cfg.QuicPort)
	quicListener, err := quic.ListenAddr(quicHostAndPort, tlsConfig, &quic.Config{
		Tracer: qlog.DefaultConnectionTracer,
	})
	if err != nil {
		return fmt.Errorf("quic.ListenAddr: %w", err)
	}
	grpcQuicListener := transport.New(quicListener)

	p.Lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			p.Logger.Info("start gRPC-QUIC server", teapot.Any("server_addr", grpcQuicListener.Addr()))
			go func() {
				if err := s.Serve(grpcQuicListener); err != nil {
					p.Logger.Fatal("unable to start gRPC-QUIC server", teapot.Error(err))
				}
			}()
			return nil
		},
		OnStop: func(ctx context.Context) error {
			var multiErr error
			p.Logger.Info("close quic listener")
			if err := grpcQuicListener.Close(); err != nil {
				multiErr = multierr.Append(err, fmt.Errorf("lis.Close: %w", err))
			}
			return multiErr
		},
	})

	return nil
}

func newTlsConfig(certFile, keyFile string) (*tls.Config, error) {
	cert, err := tls.LoadX509KeyPair(certFile, keyFile)
	if err != nil {
		return nil, fmt.Errorf("tls.LoadX509KeyPair: %w", err)
	}

	return &tls.Config{
		Certificates: []tls.Certificate{cert},
		NextProtos:   []string{"h3"},
	}, nil
}
