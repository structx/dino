package sessions

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net"
	"net/netip"
	"sync"

	"go.uber.org/multierr"
	"go.uber.org/zap"
	tunnelnet "soft.structx.io/dino/tunnel/net"
)

type actor interface {
	routeIncoming([]byte)
	handleConn()
	close() error
}

type sessionActor struct {
	sessionID string
	incoming  chan []byte
	errCh     chan error
	outbound  io.WriteCloser
	localConn net.Conn
	cleanup   func()
}

// close implements actor.
func (s *sessionActor) close() error {
	var result error
	if s.localConn != nil {
		if err := s.localConn.Close(); err != nil {
			result = multierr.Append(result, fmt.Errorf("failed to close local conn: %w", err))
		}
	}

	if s.outbound != nil {
		if err := s.outbound.Close(); err != nil {
			result = multierr.Append(result, fmt.Errorf("failed to close outbound conn: %w", err))
		}
	}

	close(s.incoming)

	return result
}

// routeIncoming implements actor.
func (s *sessionActor) routeIncoming(msg []byte) {
	s.incoming <- msg
}

// interface compliance
var _ actor = (*sessionActor)(nil)

type Mux interface {
	InitSession(tunnelnet.Conn, string, string, string) error
	RouteMsg(string, []byte) error
	CloseSession(string) error

	start(context.Context) error
	stop(context.Context) error
}

type sessionMultiplexer struct {
	log *zap.Logger

	mtx    sync.RWMutex
	actors map[string]actor
	errCh  chan error
}

// interface compliance
var _ Mux = (*sessionMultiplexer)(nil)

func newMux(logger *zap.Logger) Mux {
	return &sessionMultiplexer{
		log:    logger.Named("session_multiplexer"),
		mtx:    sync.RWMutex{},
		errCh:  make(chan error),
		actors: map[string]actor{},
	}
}

// CloseSession implements Manager.
func (s *sessionMultiplexer) CloseSession(sessionID string) error {
	s.mtx.Lock()
	defer delete(s.actors, sessionID)
	defer s.mtx.Unlock()
	if _, ok := s.actors[sessionID]; ok {
		return s.actors[sessionID].close()
	}
	return nil
}

// RouteMsg implements Manager.
func (s *sessionMultiplexer) RouteMsg(sessionID string, payload []byte) error {
	s.mtx.RLock()
	actor, ok := s.actors[sessionID]
	s.mtx.RUnlock()

	if !ok {
		return errors.New("session not active")
	}

	go actor.routeIncoming(payload)

	return nil
}

// initConn implements Manager.
func (s *sessionMultiplexer) InitSession(conn tunnelnet.Conn, sessionID, protocol, addr string) error {
	s.mtx.Lock()
	defer s.mtx.Unlock()

	ch := make(chan []byte)

	addrPort, err := netip.ParseAddrPort(addr)
	if err != nil {
		return fmt.Errorf("netip.ParseAddrPort: %w", err)
	}

	var outboundConn io.WriteCloser
	var localConn net.Conn

	switch protocol {
	// case "udp":
	// 	dfn = net.Dialer{
	// 		LocalAddr: net.UDPAddrFromAddrPort(addrPort),
	// 	}
	case "http", "https":
		tcpAddr := net.TCPAddrFromAddrPort(addrPort)
		localConn, err = net.DialTCP("tcp", nil, tcpAddr)
		if err != nil {
			return fmt.Errorf("net.DialTCP: %w", err)
		}

		outboundConn = tunnelnet.NewClientWriter(conn, sessionID)
	default:
		return fmt.Errorf("unsupported protocol %s", protocol)
	}

	actor := &sessionActor{
		sessionID: sessionID,
		localConn: localConn,
		errCh:     s.errCh,
		incoming:  ch,
		outbound:  outboundConn,
		cleanup: func() {
			if _, err := conn.Write(&tunnelnet.DataFrame{
				SessionID:      sessionID,
				IsControlFrame: true,
				CloseConn: &tunnelnet.CloseConn{
					Status: 1,
				},
			}); err != nil {
				s.log.Error("conn.Write", zap.Error(err))
			}
		},
	}

	s.actors[sessionID] = actor
	go actor.handleConn()

	return nil
}

func (s *sessionMultiplexer) start(_ context.Context) error {
	go s.worker()
	return nil
}

func (s *sessionMultiplexer) stop(_ context.Context) error {
	// TODO
	// implement stop functionality
	var result error
	for _, a := range s.actors {
		if err := a.close(); err != nil {
			result = multierr.Append(result, fmt.Errorf("failed to close actor: %w", err))
		}
	}
	close(s.errCh)
	return result
}

func (s *sessionMultiplexer) worker() {
	for err := range s.errCh {
		s.log.Error("session error occurred", zap.Error(err))
	}
}

func (sa *sessionActor) handleConn() {
	defer sa.cleanup()

	var (
		wg sync.WaitGroup

		rc = tunnelnet.NewReader(sa.incoming)
	)

	wg.Add(2)
	go func() {
		defer wg.Done()
		// read from gRPC stream and write to local conn
		if _, err := io.Copy(sa.localConn, rc); err != nil && !errors.Is(err, net.ErrClosed) {
			sa.errCh <- fmt.Errorf("failed to copy from gRPC stream to local conn: %w", err)
		}
	}()

	go func() {
		defer wg.Done()
		// read from local conn and write to gRPC stream
		if _, err := io.Copy(sa.outbound, sa.localConn); err != nil && !errors.Is(err, net.ErrClosed) {
			sa.errCh <- fmt.Errorf("failed to copy from local conn to gRPC stream: %w", err)
		}
	}()

	wg.Wait()
}
