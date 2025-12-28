package client

import (
	"errors"
	"fmt"
	"io"

	"google.golang.org/grpc"
	pb "soft.structx.io/dino/pb/rtunnel/v1"
	tunnelnet "soft.structx.io/dino/tunnel/net"
)

type clientConn struct {
	str grpc.BidiStreamingClient[pb.TunnelMessage, pb.TunnelMessage]
}

var _ tunnelnet.Conn = (*clientConn)(nil)

// Close implements net.Conn.
func (c *clientConn) Close(sessionID string) error {
	return c.str.Send(&pb.TunnelMessage{
		SessionId: sessionID,
		Payload: &pb.TunnelMessage_CloseConnection{
			CloseConnection: &pb.CloseConnection{
				StatusCode: 1,
			},
		},
	})
}

// Read implements net.Conn.
func (c *clientConn) Read() (*tunnelnet.DataFrame, error) {
	msg, err := c.str.Recv()
	if err != nil {
		if err == io.EOF {
			// stream was terminated
			return nil, io.EOF
		}
		return nil, fmt.Errorf("str.Recv: %w", err)
	}

	if msg.GetCloseConnection() != nil {
		return &tunnelnet.DataFrame{
			SessionID:      msg.GetSessionId(),
			IsControlFrame: true,
			CloseConn: &tunnelnet.CloseConn{
				Status: int(msg.GetCloseConnection().GetStatusCode()),
			},
		}, nil
	} else if msg.GetData() != nil {
		return &tunnelnet.DataFrame{
			IsControlFrame: false,
			SessionID:      msg.GetSessionId(),
			Payload:        msg.GetData(),
		}, nil
	} else if msg.GetNewConnection() != nil {
		return &tunnelnet.DataFrame{
			SessionID:      msg.GetSessionId(),
			IsControlFrame: true,
			NewConn: &tunnelnet.NewConn{
				Hostname: msg.GetNewConnection().GetDestination(),
			},
		}, nil
	}

	return nil, errors.New("unsupported message")
}

// Write implements net.Conn.
func (c *clientConn) Write(df *tunnelnet.DataFrame) (int, error) {

	if df.IsControlFrame {
		if df.CloseConn != nil {
			return 0, c.str.Send(&pb.TunnelMessage{
				SessionId: df.SessionID,
				Payload: &pb.TunnelMessage_CloseConnection{
					CloseConnection: &pb.CloseConnection{
						StatusCode: uint32(df.CloseConn.Status),
					},
				},
			})
		}
	}

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
