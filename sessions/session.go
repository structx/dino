package sessions

import (
	"fmt"
	"io"

	tunnelnet "soft.structx.io/dino/tunnel/net"
)

type activeSession struct {
	streamID  string
	sessionID string
	inbound   chan []byte
	outbound  tunnelnet.Conn
}

// Write implements net.WriteCloser.
func (a *activeSession) Write(p []byte) (n int, err error) {
	fmt.Println("outbound.Write")
	return a.outbound.Write(&tunnelnet.DataFrame{
		IsControlFrame: false,
		NewConn:        nil,
		RouteUpdate:    nil,
		SessionID:      a.sessionID,
		Payload:        p,
	})
}

// Close implements net.ReadCloser.
func (a *activeSession) Close() error {
	defer close(a.inbound)
	if _, err := a.outbound.Write(&tunnelnet.DataFrame{
		SessionID:      a.sessionID,
		IsControlFrame: true,
		CloseConn: &tunnelnet.CloseConn{
			Status: 1,
		},
	}); err != nil {
		return fmt.Errorf("outbound.Write: %w", err)
	}
	return nil
}

// Read implements net.ReadCloser.
func (a *activeSession) Read(p []byte) (n int, err error) {
	msg, ok := <-a.inbound
	if !ok {
		return 0, io.EOF
	}
	return copy(p, msg), nil
}

func (a *activeSession) openConn() error {
	_, err := a.outbound.Write(&tunnelnet.DataFrame{
		SessionID:      a.sessionID,
		IsControlFrame: true,
		NewConn: &tunnelnet.NewConn{
			Hostname: "whoami.dino.local",
		},
	})
	return err
}

var _ tunnelnet.ReadCloser = (*activeSession)(nil)
var _ tunnelnet.WriteCloser = (*activeSession)(nil)
