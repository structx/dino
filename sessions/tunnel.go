package sessions

import (
	"io"
	"sync"

	"github.com/google/uuid"
	tunnelnet "github.com/structx/dino/tunnel/net"
)

type activeTunnel struct {
	streamID string
	mtx      sync.Mutex
	sessions map[string]activeSession
	stream   tunnelnet.Conn
}

func (a *activeTunnel) worker() {
	for {
		df, err := a.stream.Read()
		if err != nil {
			if err == io.EOF {

			}
		}

		a.mtx.Lock()
		session, ok := a.sessions[df.SessionID]
		a.mtx.Unlock()

		if !ok {

		}

		session.inbound <- df.Payload
	}
}

func (a *activeTunnel) registerSession() activeSession {
	sessionID := uuid.New().String()
	ch := make(chan []byte)
	session := activeSession{
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
