package router

import "sync"

type Route struct {
	Protocol string
	IP       string
	Port     string
}

type Mux interface {
	Add(string, string, string, string)
	Del(string)
	Match(string) (Route, bool)
}

type tunnelRouter struct {
	mtx    sync.RWMutex
	routes map[string]Route
}

// interface compliance
var _ Mux = (*tunnelRouter)(nil)

func newMux() Mux {
	return &tunnelRouter{
		mtx:    sync.RWMutex{},
		routes: map[string]Route{},
	}
}

// Add implements Router.
func (t *tunnelRouter) Add(hostname, protocol, ip string, port string) {
	t.mtx.Lock()
	defer t.mtx.Unlock()
	t.routes[hostname] = Route{Protocol: protocol, IP: ip, Port: port}
}

// Del implements Router.
func (t *tunnelRouter) Del(hostname string) {
	t.mtx.Lock()
	defer t.mtx.Unlock()
	delete(t.routes, hostname)
}

// Match implements Router.
func (t *tunnelRouter) Match(hostname string) (Route, bool) {
	t.mtx.RLock()
	defer t.mtx.RUnlock()
	route, ok := t.routes[hostname]
	if !ok {
		return Route{}, false
	}
	return route, true
}
