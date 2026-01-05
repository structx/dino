package gateway

import (
	"context"
	"errors"
	"fmt"
	"net"
	"net/http"
	"time"

	"go.uber.org/fx"
	"go.uber.org/multierr"
	"go.uber.org/zap"
	"go.uber.org/zap/zapgrpc"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/grpclog"
	"google.golang.org/grpc/status"
	"soft.structx.io/dino/setup"

	"google.golang.org/grpc/health"
	healthpb "google.golang.org/grpc/health/grpc_health_v1"
)

var (
	ErrMissingMetadata = status.Error(codes.InvalidArgument, "missing metadata")
	ErrVerifyToken     = status.Error(codes.Internal, "verify token")
	ErrInvalidTunnelID = status.Error(codes.InvalidArgument, "missing tunnel id")
	ErrInvalidToken    = status.Error(codes.Unauthenticated, "invalid token")
)

// UnaeryInterceptor
type UnaryInterceptor = grpc.UnaryServerInterceptor

// StreamInterceptor
type StreamInterceptor = grpc.StreamServerInterceptor

// ServerStream
type ServerStream = grpc.ServerStream

// StreamServerInfo
type StreamServerInfo = grpc.StreamServerInfo

// StreamHandler
type StreamHandler = grpc.StreamHandler

// Transport
type Transport struct {
	ServiceDesc *grpc.ServiceDesc
	Service     any
}

// Params
type Params struct {
	fx.In

	Lc fx.Lifecycle

	Logger *zap.Logger

	ServerConfig *setup.Server
	ProxyConfig  *setup.Proxy

	Proxy http.Handler

	Transports []Transport `group:"transport"`

	// UnaryInterceptors  []UnaryInterceptor
	// StreamInterceptors []StreamInterceptor
}

// Module
var Module = fx.Module("gateway", fx.Invoke(invokeModule))

func invokeModule(p Params) error {

	hostAndPort := net.JoinHostPort(p.ServerConfig.Host, p.ServerConfig.Port)

	p.Logger.Info("num transports", zap.Int("len", len(p.Transports)))

	gs := grpc.NewServer()
	for _, tr := range p.Transports {
		gs.RegisterService(tr.ServiceDesc, tr.Service)
	}

	healthcheck := health.NewServer()
	healthpb.RegisterHealthServer(gs, healthcheck)

	rpcLogger := zapgrpc.NewLogger(p.Logger)
	grpclog.SetLoggerV2(rpcLogger)

	mux := http.NewServeMux()
	mux.HandleFunc("/", p.Proxy.ServeHTTP)

	proxyAddr := net.JoinHostPort(p.ProxyConfig.Host, p.ProxyConfig.Port)
	h1 := &http.Server{
		Addr:              proxyAddr,
		Handler:           p.Proxy,
		ReadTimeout:       time.Second * p.ProxyConfig.ReadTimeout,
		ReadHeaderTimeout: time.Second * p.ProxyConfig.ReadHeaderTimeout,
		WriteTimeout:      time.Second * p.ProxyConfig.WriteTimeout,
		IdleTimeout:       time.Second * p.ProxyConfig.IdleTimeout,
	}

	h2h := h2c.NewHandler(gs, &http2.Server{})

	pr := new(http.Protocols)
	pr.SetHTTP1(true)
	pr.SetUnencryptedHTTP2(true)

	hls := &http.Server{
		Addr:      hostAndPort,
		Handler:   h2h,
		Protocols: pr,
	}

	p.Lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			p.Logger.Info("start http/1 proxy server", zap.String("server_addr", h1.Addr))
			go func() {
				if err := h1.ListenAndServe(); err != nil && errors.Is(err, http.ErrServerClosed) {
					p.Logger.Fatal("start http/1 proxy server", zap.Error(err))
				}
			}()

			p.Logger.Info("start hls server", zap.String("server_addr", hls.Addr))
			go func() {
				if err := hls.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
					p.Logger.Fatal("unable to start hls server", zap.Error(err))
				}
			}()

			p.Logger.Info("start gRPC healthceck server")
			go func() {
				next := healthpb.HealthCheckResponse_SERVING
				for {
					healthcheck.SetServingStatus("", next)
					if next == healthpb.HealthCheckResponse_NOT_SERVING {
						next = healthpb.HealthCheckResponse_SERVING
					} else {
						next = healthpb.HealthCheckResponse_NOT_SERVING
					}

					time.Sleep(time.Second * 3)
				}
			}()

			return nil
		},
		OnStop: func(ctx context.Context) error {
			var multiErr error

			p.Logger.Info("shutdown http/1 proxy server")
			if err := h1.Shutdown(ctx); err != nil {
				multiErr = multierr.Append(multiErr, fmt.Errorf("h1.Shutdown: %w", err))
			}

			p.Logger.Info("shutdown hls server")
			if err := hls.Shutdown(ctx); err != nil {
				multiErr = multierr.Append(multiErr, fmt.Errorf("hls.Shutdown: %w", err))
			}

			p.Logger.Info("shutdown gRPC-QUIC server")
			timer := time.AfterFunc(time.Second*10, func() {
				gs.Stop()
			})
			defer timer.Stop()
			gs.GracefulStop()

			return multiErr
		},
	})
	return nil
}
