package proxy

import (
	"io"
	"net"
	"net/http"
	"sync"

	"go.uber.org/fx"
	"go.uber.org/zap"
	"soft.structx.io/dino/internal/routes"
	"soft.structx.io/dino/sessions"
)

// Handler
type Handler interface {
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

type handler struct {
	log *zap.Logger

	routeSvc routes.Service

	mux sessions.Muxer
}

// Module
var Module = fx.Module("proxy", fx.Provide(newModule))

func newModule(p Params) Result {
	return Result{
		Handler: &handler{
			log:      p.Logger,
			routeSvc: p.RouteService,
			mux:      p.Mux,
		},
	}
}

// ServeHTTP
func (h *handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	hostname := r.Header.Get("host")
	if hostname == "" {
		hostname = r.Host
	}

	if hostname == "" {
		hostname = r.URL.Host
	}

	if hostname == "" {
		h.log.Error("missing host header")
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	host, _, err := net.SplitHostPort(hostname)
	if err != nil {
		h.log.Error("split host and port", zap.Error(err))
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	activeUID, err := h.routeSvc.Active(ctx, host)
	if err != nil {
		h.log.Error("active route", zap.Error(err))
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}

	hijacker, ok := w.(http.Hijacker)
	if !ok {
		h.log.Error("hijacking not supported")
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	// hijack end user connection
	conn, buf, err := hijacker.Hijack()
	if err != nil {
		h.log.Error("hijacker.Hijack", zap.Error(err))
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	defer func() { _ = conn.Close() }()

	rc, wc, cleanup, ok := h.mux.UnarySession(ctx, activeUID)
	if !ok {
		h.log.Error("unary session invalid")
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	defer cleanup()

	if err := r.Write(wc); err != nil {
		h.log.Error("r.Write", zap.Error(err))
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	var wg sync.WaitGroup

	wg.Add(2)
	go func() {
		defer wg.Done()
		go func() {
			if _, err := io.Copy(wc, buf); err != nil {
				h.log.Error("io.Copy", zap.Error(err))
				http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
				return
			}
		}()
	}()

	go func() {
		defer wg.Done()
		if _, err := io.Copy(conn, rc); err != nil {
			h.log.Error("io.Copy", zap.Error(err))
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
	}()

	wg.Wait()
}
