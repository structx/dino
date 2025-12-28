package net

import (
	"fmt"
	"io"
	"net"
	"sync"
)

// Writer
type Writer interface {
	io.Writer
}

type tunnelWriter struct {
	conn      Conn
	sessionID string
	isClosed  bool
	close     sync.Once
}

// NewWriter
func NewWriter(conn Conn, sessionID string) WriteCloser {
	return &tunnelWriter{
		conn:      conn,
		isClosed:  false,
		close:     sync.Once{},
		sessionID: sessionID,
	}
}

// NewClientWriter
func NewClientWriter(conn Conn, sessionID string) WriteCloser {
	return &tunnelWriter{
		conn:      conn,
		sessionID: sessionID,
		isClosed:  false,
		close:     sync.Once{},
	}
}

// Close implements WriteCloser.
func (tw *tunnelWriter) Close() error {
	err := tw.conn.Close(tw.sessionID)
	if err != nil {
		return fmt.Errorf("conn.Close: %w", err)
	}
	tw.close.Do(func() {
		tw.isClosed = true
	})
	return nil
}

// Write implements WriteCloser.
func (tw *tunnelWriter) Write(p []byte) (n int, err error) {
	if tw.isClosed {
		return 0, net.ErrClosed
	}

	return tw.conn.Write(&DataFrame{
		SessionID:      tw.sessionID,
		NewConn:        nil,
		RouteUpdate:    nil,
		IsControlFrame: false,
		Payload:        p,
	})
}
