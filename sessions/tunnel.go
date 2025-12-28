package sessions

import (
	"fmt"
	"io"
	"sync"

	"github.com/google/uuid"
	tunnelnet "soft.structx.io/dino/tunnel/net"
)

type activeTunnel struct {
	streamID string
	mtx      sync.Mutex
	sessions map[string]*activeSession
	stream   tunnelnet.Conn
}

func (a *activeTunnel) worker() {
	for {
		df, err := a.stream.Read()
		if err != nil {
			if err == io.EOF {
				return
			}
			fmt.Println(err)
			return
		}

		a.mtx.Lock()
		session, ok := a.sessions[df.SessionID]
		a.mtx.Unlock()

		if !ok {
			fmt.Println("no session")
			return
		}

		session.inbound <- df.Payload
	}
}

func (a *activeTunnel) registerSession() *activeSession {
	sessionID := uuid.New().String()
	ch := make(chan []byte)
	session := &activeSession{
		streamID:  a.streamID,
		sessionID: sessionID,
		inbound:   ch,
		outbound:  a.stream,
	}

	a.mtx.Lock()
	a.sessions[sessionID] = session
	a.mtx.Unlock()

	return session
}

func (a *activeTunnel) deregisterSession(sessionID string) error {
	a.mtx.Lock()

	session, ok := a.sessions[sessionID]
	if !ok {
		// session is not active
		return nil
	}
	a.mtx.Unlock()

	// close inbound channel
	defer close(session.inbound)

	// remove session from map
	defer delete(a.sessions, sessionID)

	// send close connection signal
	if _, err := a.stream.Write(&tunnelnet.DataFrame{
		SessionID:      sessionID,
		IsControlFrame: true,
		CloseConn: &tunnelnet.CloseConn{
			Status: 1,
		},
		NewConn:     nil,
		Payload:     nil,
		RouteUpdate: nil,
	}); err != nil {
		return fmt.Errorf("stream.Write: %w", err)
	}

	return nil
}
