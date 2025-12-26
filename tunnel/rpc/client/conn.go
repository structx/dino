package client

import (
	"errors"
	"fmt"
	"io"

	pb "github.com/structx/dino/pb/rtunnel/v1"
	tunnelnet "github.com/structx/dino/tunnel/net"
	"google.golang.org/grpc"
)

type clientConn struct {
	str grpc.BidiStreamingClient[pb.TunnelMessage, pb.TunnelMessage]
}

var _ tunnelnet.ClientConn = (*clientConn)(nil)

// Close implements net.Conn.
func (c *clientConn) Close(sessionID string) error {
	return c.str.Send(&pb.TunnelMessage{
		SessionId: sessionID,
		Payload:   &pb.TunnelMessage_CloseConnection{},
	})
}

// Read implements net.Conn.
func (c *clientConn) Read() (*tunnelnet.ClientFrame, error) {
	msg, err := c.str.Recv()
	if err != nil {
		if err == io.EOF {
			// stream was terminated
			return nil, io.EOF
		}
		return nil, fmt.Errorf("str.Recv: %w", err)
	}

	if msg.GetCloseConnection() != nil {
		return nil, io.EOF
	} else if msg.GetData() != nil {
		return &tunnelnet.ClientFrame{
			SessionID: msg.GetSessionId(),
			Payload:   msg.GetData(),
		}, nil
	}

	return nil, errors.New("unsupported message")
}

// Write implements net.Conn.
func (c *clientConn) Write(df *tunnelnet.ClientFrame) (int, error) {
	if err := c.str.Send(&pb.TunnelMessage{
		SessionId: df.SessionID,
		Payload: &pb.TunnelMessage_Data{
			Data: df.Payload,
		},
	}); err != nil {
		return 0, fmt.Errorf("str.Send: %w", err)
	}
	return len(df.Payload), nil
}
