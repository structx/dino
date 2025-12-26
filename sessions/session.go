package sessions

import (
	"io"

	tunnelnet "github.com/structx/dino/tunnel/net"
)

type activeSession struct {
	streamID  string
	sessionID string
	inbound   chan []byte
	outbound  tunnelnet.Conn
}

// Write implements net.WriteCloser.
func (a *activeSession) Write(p []byte) (n int, err error) {
	return a.outbound.Write(&tunnelnet.DataFrame{
		SessionID: a.sessionID,
		Payload:   p,
	})
}

// Close implements net.ReadCloser.
func (a *activeSession) Close() error {
	close(a.inbound)
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

var _ tunnelnet.ReadCloser = (*activeSession)(nil)
var _ tunnelnet.WriteCloser = (*activeSession)(nil)
