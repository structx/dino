package pubsub

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"sync"
)

// Msg
type Msg interface {
	json.Marshaler
	json.Unmarshaler
}

// RouteConfig
type RouteConfig struct {
	TunnelUID    string `json:"tunnel_uid"`
	Hostname     string `json:"hostname"`
	DestProtocol string `json:"dest_protocol"`
	DestAddr     string `json:"dest_addr"`
	DestPort     uint32 `json:"dest_port"`
	IsDelete     bool   `json:"is_delete"`
}

// MarshalJSON implements [Msg].
func (r *RouteConfig) MarshalJSON() ([]byte, error) {
	return json.Marshal(r)
}

// UnmarshalJSON implements [Msg].
func (r *RouteConfig) UnmarshalJSON(b []byte) error {
	return json.Unmarshal(b, &r)
}

// interface compliance
var _ Msg = (*RouteConfig)(nil)

// Broker
type Broker interface {
	Publish(string, interface{}) error

	Subscribe(string) chan string
	Unsubsribe(string)
}

// inmemory implementation of message broker
type inmemoryBroker struct {
	mtx  sync.RWMutex
	subs map[string]chan string
}

// interface compliance
var _ Broker = (*inmemoryBroker)(nil)

func newBroker() Broker {
	return &inmemoryBroker{
		mtx:  sync.RWMutex{},
		subs: map[string]chan string{},
	}
}

// Publish implements [Broker].
func (i *inmemoryBroker) Publish(topic string, payload interface{}) error {
	jsonbytes, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("json.Marshal: %w", err)
	}

	encoded := base64.StdEncoding.EncodeToString(jsonbytes)

	i.mtx.Lock()
	defer i.mtx.Unlock()

	if ch, ok := i.subs[topic]; !ok {
		i.subs[topic] = make(chan string)
		i.subs[topic] <- encoded
	} else {
		ch <- encoded
	}

	return nil
}

// Subscribe implements [Broker].
func (i *inmemoryBroker) Subscribe(topic string) chan string {
	i.mtx.RLock()
	defer i.mtx.RUnlock()
	if ch, ok := i.subs[topic]; !ok {
		i.subs[topic] = make(chan string)
		return i.subs[topic]
	} else {
		return ch
	}
}

// Unsubsribe implements [Broker].
func (i *inmemoryBroker) Unsubsribe(topic string) {
	i.mtx.Lock()
	defer i.mtx.Unlock()
	delete(i.subs, topic)
}
