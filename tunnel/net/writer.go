package net

import (
	"fmt"
	"io"
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
func NewClientWriter(clientConn ClientConn, sessionID string) WriteCloser {
	return &tunnelWriter{}
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
		return 0, io.EOF
	}

	if len(p) == 0 {
		return 0, nil
	}

	return tw.conn.Write(&DataFrame{
		SessionID: tw.sessionID,
		Payload:   p,
	})
}
