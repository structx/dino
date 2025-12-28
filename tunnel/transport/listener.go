package transport

import (
	"context"
	"fmt"
	"net"

	"github.com/quic-go/quic-go"
)

type quicListener struct {
	qlis *quic.Listener
}

// interface compliance
var _ net.Listener = (*quicListener)(nil)

// New
func New(listener *quic.Listener) net.Listener {
	return &quicListener{
		qlis: listener,
	}
}

// Accept net listener implementation
func (l *quicListener) Accept() (net.Conn, error) {
	ctx := context.Background()
	conn, err := l.qlis.Accept(ctx)
	if err != nil {
		return nil, fmt.Errorf("qlis.Accept: %w", err)
	}

	stream, err := conn.AcceptStream(ctx)
	if err != nil {
		return nil, fmt.Errorf("conn.AcceptStream: %w", err)
	}

	return &QuicConn{connection: conn, stream: stream}, nil
}

// Close
func (l *quicListener) Close() error { return l.qlis.Close() }

// Addr
func (l *quicListener) Addr() net.Addr { return l.qlis.Addr() }
