package sessions

import (
	"context"
	"encoding/base64"
	"errors"
	"fmt"
	"sync"

	"github.com/structx/teapot"
	"soft.structx.io/dino/internal/routes"
	"soft.structx.io/dino/internal/tunnel"
	"soft.structx.io/dino/pubsub"
	tunnelnet "soft.structx.io/dino/tunnel/net"
)

// Multiplexer
type Multiplexer interface {
	// RegisterTunnel
	RegisterTunnel(context.Context, tunnelnet.Conn, string) error

	// UnarySession
	UnarySession(context.Context, string) (tunnelnet.ReadCloser, tunnelnet.WriteCloser, func(), bool)

	// SyncRoutes
	SyncRoutes(context.Context, string) error
}

type sessionMultiplexer struct {
	ctx      context.Context
	cancelFn context.CancelFunc

	log *teapot.Logger

	mtx     sync.Mutex
	tunnels map[string]*activeTunnel

	tunnelSvc tunnel.Service
	routeSvc  routes.Service

	broker pubsub.Broker
}

func newMux(
	logger *teapot.Logger,
	broker pubsub.Broker,
	tsvc tunnel.Service,
	rsvc routes.Service,
) *sessionMultiplexer {
	ctx, cancel := context.WithCancel(context.Background())
	return &sessionMultiplexer{
		log:       logger,
		ctx:       ctx,
		cancelFn:  cancel,
		mtx:       sync.Mutex{},
		tunnels:   make(map[string]*activeTunnel),
		tunnelSvc: tsvc,
		routeSvc:  rsvc,
		broker:    broker,
	}
}

// RegisterTunnel
func (m *sessionMultiplexer) RegisterTunnel(ctx context.Context, conn tunnelnet.Conn, tunnelUID string) error {
	m.mtx.Lock()
	defer m.mtx.Unlock()

	m.tunnels[tunnelUID] = &activeTunnel{
		streamID: tunnelUID,
		mtx:      sync.Mutex{},
		sessions: make(map[string]*activeSession),
		stream:   conn,
	}
	go m.tunnels[tunnelUID].worker()

	return nil
}

// UnarySession
func (m *sessionMultiplexer) UnarySession(ctx context.Context, tunnelUID string) (tunnelnet.ReadCloser, tunnelnet.WriteCloser, func(), bool) {
	m.mtx.Lock()
	defer m.mtx.Unlock()

	tunnel, ok := m.tunnels[tunnelUID]
	if !ok {
		return nil, nil, nil, false
	}

	session := tunnel.registerSession()
	err := session.openConn()
	if err != nil {
		m.log.Error("open session connection", teapot.Error(err))
		return nil, nil, nil, false
	}

	cleanup := func() {
		if err := tunnel.deregisterSession(session.sessionID); err != nil {
			m.log.Error("deregister session", teapot.Error(err))
		}
	}

	return session, session, cleanup, true

}

// SyncRoutes
func (m *sessionMultiplexer) SyncRoutes(ctx context.Context, tunnelUID string) error {
	routes, err := m.routeSvc.Sync(ctx, tunnelUID)
	if err != nil {
		return fmt.Errorf("failed to sync routes: %w", err)
	}

	m.mtx.Lock()
	tunnel, ok := m.tunnels[tunnelUID]
	m.mtx.Unlock()
	if !ok {
		return errors.New("tunnel is not active")
	}

	for _, r := range routes {
		if _, err := tunnel.stream.Write(&tunnelnet.DataFrame{
			SessionID:      tunnelUID,
			IsControlFrame: true,
			RouteUpdate: &tunnelnet.RouteUpdate{
				Hostname:     r.Hostname,
				DestProtocol: r.DestinationProtocol,
				DestIP:       r.DestinationIP,
				DestPort:     r.DestinationPort,
				IsDelete:     !r.Enabled,
			},
			NewConn:   nil,
			CloseConn: nil,
			Payload:   nil,
		}); err != nil {
			return fmt.Errorf("failed to write route config: %w", err)
		}
	}

	return nil
}

func (m *sessionMultiplexer) start(_ context.Context) error {
	ch := m.broker.Subscribe("dino.routes.*")
	go m.subscription(ch)
	return nil
}

func (m *sessionMultiplexer) stop(_ context.Context) error {
	m.cancelFn()
	return nil
}

func (m *sessionMultiplexer) subscription(ch chan string) {
	for {
		select {
		case <-m.ctx.Done():
			return
		case msg := <-ch:
			rcfg, err := decodeMessge(msg)
			if err != nil {
				m.log.Error("decode route config", teapot.Error(err))
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
					m.log.Error("write route config", teapot.Error(err))
				}
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
