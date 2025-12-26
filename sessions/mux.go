package sessions

import (
	"context"
	"encoding/base64"
	"fmt"
	"sync"

	"github.com/structx/dino/internal/routes"
	"github.com/structx/dino/internal/tunnel"
	"github.com/structx/dino/pubsub"
	tunnelnet "github.com/structx/dino/tunnel/net"
	"go.uber.org/zap"
)

type Muxer interface {
	// RegisterTunnel
	RegisterTunnel(context.Context, tunnelnet.Conn, string) error

	// UnarySession
	UnarySession(context.Context, string) (tunnelnet.ReadCloser, tunnelnet.WriteCloser, bool)

	// SyncRouteConfig
	// SyncRouteConfig(context.Context) error

	// Cleanup session closer
	Cleanup(string)
}

type sessionMultiplexer struct {
	log *zap.Logger

	mtx     sync.Mutex
	tunnels map[string]*activeTunnel

	tunnelSvc tunnel.Service
	routeSvc  routes.Service

	broker pubsub.Broker
}

func newMux(logger *zap.Logger, tsvc tunnel.Service) Muxer {
	return &sessionMultiplexer{
		log:       logger,
		mtx:       sync.Mutex{},
		tunnels:   make(map[string]*activeTunnel),
		tunnelSvc: tsvc,
	}
}

// RegisterTunnel
func (m *sessionMultiplexer) RegisterTunnel(ctx context.Context, conn tunnelnet.Conn, tunnelUID string) error {

	m.mtx.Lock()
	defer m.mtx.Unlock()

	m.tunnels[tunnelUID] = &activeTunnel{
		streamID: tunnelUID,
		mtx:      sync.Mutex{},
		sessions: make(map[string]activeSession),
		stream:   conn,
	}
	go m.tunnels[tunnelUID].worker()

	return nil
}

// UnarySession
func (m *sessionMultiplexer) UnarySession(ctx context.Context, hostname string) (tunnelnet.ReadCloser, tunnelnet.WriteCloser, bool) {
	route, err := m.routeSvc.Active(ctx, hostname)
	if err != nil {
		return nil, nil, false
	}
	streamID := route.ID

	m.mtx.Lock()
	defer m.mtx.Unlock()

	tunnel, ok := m.tunnels[streamID]
	if !ok {
		return nil, nil, false
	}

	session := tunnel.registerSession()

	return &session, &session, true

}

// Cleanup
func (m *sessionMultiplexer) Cleanup(sessionID string) {
	// m.mtx.Lock()
	// defer m.mtx.Unlock()
	// if session, ok := m.sessions[sessionID]; ok {
	// 	close(session.inbound)
	// }
}

func (m *sessionMultiplexer) subscribe() {
	ch := m.broker.Subscribe("dino.routes")
	go m.subscription(ch)
}

func (m *sessionMultiplexer) subscription(ch chan string) {
	for msg := range ch {
		rcfg, err := decodeMessge(msg)
		if err != nil {
			m.log.Error("decode route config", zap.Error(err))
			continue
		}

		if t, ok := m.tunnels[rcfg.TunnelUID]; !ok {
			// tunnel is not active continue looping
			continue
		} else {
			if _, err := t.stream.Write(&tunnelnet.DataFrame{
				SessionID: t.streamID,
				RouteUpdate: &tunnelnet.RouteUpdate{
					Hostname:     rcfg.Hostname,
					DestProtocol: rcfg.DestProtocol,
					DestIP:       rcfg.DestAddr,
					DestPort:     rcfg.DestPort,
					IsDelete:     rcfg.IsDelete,
				},
			}); err != nil {
				m.log.Error("write route config", zap.Error(err))
			}
		}
	}
}

func decodeMessge(s string) (*pubsub.RouteConfig, error) {
	b64, err := base64.StdEncoding.DecodeString(s)
	if err != nil {
		return nil, fmt.Errorf("failed to decode string: %w", err)
	}

	var cfg pubsub.RouteConfig
	err = cfg.UnmarshalJSON(b64)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal message: %w", err)
	}

	return &cfg, nil
}
