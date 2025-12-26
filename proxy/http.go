package proxy

import (
	"io"
	"net/http"
	"sync"

	"github.com/structx/dino/internal/routes"
	"github.com/structx/dino/sessions"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

// ProxyHandler
type ProxyHandler interface {
	ServeHTTP(http.ResponseWriter, *http.Request)
}

// Params
type Params struct {
	fx.In

	Logger *zap.Logger

	Mux sessions.Muxer

	RouteService routes.Service
}

// Result
type Result struct {
	fx.Out

	Handler http.Handler
}

type proxyHandler struct {
	log *zap.Logger

	routeSvc routes.Service

	mux sessions.Muxer
}

var Module = fx.Module("proxy", fx.Provide(newModule))

func newModule(p Params) Result {
	return Result{
		Handler: &proxyHandler{
			log:      p.Logger,
			routeSvc: p.RouteService,
			mux:      p.Mux,
		},
	}
}

// ServeHTTP
func (ph *proxyHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	hostname := r.Header.Get("Host")
	if hostname == "" {
		return
	}

	route, err := ph.routeSvc.Active(ctx, hostname)
	if err != nil {
		return
	}

	hijacker, ok := w.(http.Hijacker)
	if !ok {
		ph.log.Error("hijacking not supported")
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	// hijack end user connection
	conn, _, err := hijacker.Hijack()
	if err != nil {
		ph.log.Error("hijacker.Hijack", zap.Error(err))
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	defer func() { _ = conn.Close() }()

	rc, wc, ok := ph.mux.UnarySession(ctx, route.Tunnel)
	if !ok {
		ph.log.Error("unary session invalid")
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		if _, err := io.Copy(wc, conn); err != nil {
			ph.log.Error("failed to copy from connection to tunnel", zap.Error(err))
		}
	}()
	wg.Add(1)
	go func() {
		defer wg.Done()
		if _, err := io.Copy(conn, rc); err != nil {
			ph.log.Error("failed to copy from tunnel to connection", zap.Error(err))
		}
	}()
	wg.Wait()
}
