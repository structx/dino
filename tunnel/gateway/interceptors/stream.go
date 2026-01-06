package interceptors

import (
	"github.com/structx/teapot"
	"google.golang.org/grpc/metadata"
	"soft.structx.io/dino/gateway"
	tunnelgateway "soft.structx.io/dino/tunnel/gateway"
	"soft.structx.io/dino/tunnel/verifier"
)

const (
	authorizationHeader = "authorization"
	tunnelHeader        = "tunnel-id"
)

type serverInterceptor struct {
	l *teapot.Logger
	v verifier.Verifier
}

type wrappedStream struct {
	gateway.ServerStream

	l *teapot.Logger
}

func newServerInterceptor(logger *teapot.Logger, verifier verifier.Verifier) tunnelgateway.StreamServerInterceptor {
	s := &serverInterceptor{l: logger, v: verifier}
	return s.StreamInterceptor
}

func newWrappedStream(s gateway.ServerStream, logger *teapot.Logger) gateway.ServerStream {
	return &wrappedStream{
		l:            logger,
		ServerStream: s,
	}
}

func (w *wrappedStream) RecvMsg(m any) error {
	w.l.Debug("recv msg", teapot.Any("msg", m))
	return w.ServerStream.RecvMsg(m)
}

func (w *wrappedStream) SendMsg(m any) error {
	w.l.Debug("send msg", teapot.Any("msg", m))
	return w.ServerStream.SendMsg(m)
}

// StreamInterceptor
func (si *serverInterceptor) StreamInterceptor(srv any, ss gateway.ServerStream, _ *gateway.StreamServerInfo, handler gateway.StreamHandler) error {
	ctx := ss.Context()

	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return gateway.ErrMissingMetadata
	}

	if !valid(md[authorizationHeader]) {
		return gateway.ErrInvalidToken
	}

	if !valid(md[tunnelHeader]) {
		return gateway.ErrInvalidTunnelID
	}

	// _, err := si.v.VerifyToken(ctx, md[tunnelHeader][0], md[authorizationHeader][0])
	// if err != nil {
	// 	si.l.Error("verify token", zap.Error(err))
	// 	return gateway.ErrVerifyToken
	// }

	err := handler(srv, newWrappedStream(ss, si.l))
	if err != nil {
		si.l.Error("RPC failed with error", teapot.Error(err))
	}

	return err
}

func valid(s []string) bool {
	return len(s) > 0
}
