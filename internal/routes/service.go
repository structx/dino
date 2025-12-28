package routes

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"soft.structx.io/dino/database"
	"soft.structx.io/dino/internal/routes/queries"
	"soft.structx.io/dino/pubsub"
)

// RouteCreate
type RouteCreate struct {
	Tunnel string

	Hostname string

	DestinationProtocol string
	DestinationIP       string
	DestinationPort     uint32
}

// Route
type Route struct {
	ID                  string
	Tunnel              string
	Hostname            string
	DestinationProtocol string
	DestinationIP       string
	DestinationPort     uint32
	Enabled             bool
	CreatedAt           time.Time
	UpdatedAt           *time.Time
}

// RouteUpdate
type RouteUpdate struct {
	ID                  string
	Hostname            string
	DestinationProtocol string
	DestinationIP       string
	DestinationPort     uint32
	Enabled             bool
}

// RoutePartial
type RoutePartial struct {
	ID        string
	Hostname  string
	Enabled   bool
	CreatedAt time.Time
}

// RouteList
type RouteList struct {
	Tunnel        string
	Limit, Offset int32
}

// Service
//
//go:generate mockgen -source service.go -destination mock_service_test.go -package routes Service
type Service interface {
	// Create
	Create(context.Context, RouteCreate) (Route, error)
	// Get
	Get(context.Context, string) (Route, error)
	// List
	List(context.Context, RouteList) ([]RoutePartial, error)
	// Update
	Update(context.Context, RouteUpdate) (Route, error)
	// Delete
	Delete(context.Context, string) error

	// Active
	Active(context.Context, string) (string, error)
	// Sync
	Sync(context.Context, string) ([]Route, error)
}

type serviceImpl struct {
	db database.DBTX
	br pubsub.Broker
}

// interface compliance
var _ Service = (*serviceImpl)(nil)

func newService(dbtx database.DBTX, broker pubsub.Broker) Service {
	return &serviceImpl{
		db: dbtx,
		br: broker,
	}
}

// Create
func (s *serviceImpl) Create(ctx context.Context, create RouteCreate) (Route, error) {
	timeout, cancel := context.WithTimeout(ctx, time.Second*3)
	defer cancel()

	params := queries.InsertRouteParams{
		TunnelName:          create.Tunnel,
		Hostname:            create.Hostname,
		DestinationProtocol: create.DestinationProtocol,
		DestinationIp:       create.DestinationIP,
		DestinationPort:     int32(create.DestinationPort),
	}

	sqlRoute, err := queries.New(s.db).InsertRoute(timeout, params)
	if err != nil {
		return Route{}, fmt.Errorf("failed to execute insert route query: %w", err)
	}

	if err := s.br.Publish("dino.routes", &pubsub.RouteConfig{
		TunnelUID:    sqlRoute.TunnelName,
		Hostname:     sqlRoute.Hostname,
		DestProtocol: sqlRoute.DestinationProtocol,
		DestAddr:     sqlRoute.DestinationIp,
		DestPort:     uint32(sqlRoute.DestinationPort),
		IsDelete:     false,
	}); err != nil {
		return Route{}, fmt.Errorf("failed to publish route creation: %w", err)
	}

	return dtoRoute(sqlRoute), nil
}

// Get
func (s *serviceImpl) Get(ctx context.Context, routeID string) (Route, error) {
	timeout, cancel := context.WithTimeout(ctx, time.Second*3)
	defer cancel()

	routeUID, err := uuid.Parse(routeID)
	if err != nil {
		return Route{}, fmt.Errorf("uuid.Parse: %w", err)
	}

	sqlRoute, err := queries.New(s.db).SelectRoute(timeout, routeUID)
	if err != nil {
		return Route{}, fmt.Errorf("failed to execute select route query: %w", err)
	}

	return dtoRoute(sqlRoute), nil
}

// List
func (s *serviceImpl) List(ctx context.Context, args RouteList) ([]RoutePartial, error) {
	timeout, cancel := context.WithTimeout(ctx, time.Second*3)
	defer cancel()

	params := queries.SelectManyRoutesParams{
		TunnelName: args.Tunnel,
		Limit:      args.Limit,
		Offset:     args.Offset,
	}

	rows, err := queries.New(s.db).SelectManyRoutes(timeout, params)
	if err != nil {
		return nil, fmt.Errorf("failed to execute select many routes queries: %w", err)
	}

	return dtoPartials(rows), nil
}

// Update
func (s *serviceImpl) Update(ctx context.Context, args RouteUpdate) (Route, error) {
	timeout, cancel := context.WithTimeout(ctx, time.Second*3)
	defer cancel()

	routeUID, err := uuid.Parse(args.ID)
	if err != nil {
		return Route{}, fmt.Errorf("uuid.Parse: %w", err)
	}

	params := queries.UpdateRouteParams{
		ID:              routeUID,
		Hostname:        args.Hostname,
		DestinationIp:   args.DestinationIP,
		DestinationPort: int32(args.DestinationPort),
	}

	sqlRoute, err := queries.New(s.db).UpdateRoute(timeout, params)
	if err != nil {
		return Route{}, fmt.Errorf("failed to execute update route query: %w", err)
	}

	if err := s.br.Publish("dino.routes", &pubsub.RouteConfig{
		TunnelUID:    sqlRoute.TunnelName,
		Hostname:     sqlRoute.Hostname,
		DestProtocol: sqlRoute.DestinationProtocol,
		DestAddr:     sqlRoute.DestinationIp,
		DestPort:     uint32(sqlRoute.DestinationPort),
		IsDelete:     false,
	}); err != nil {
		return Route{}, fmt.Errorf("failed to publish route creation: %w", err)
	}

	return dtoRoute(sqlRoute), nil
}

// Delete
func (s *serviceImpl) Delete(ctx context.Context, routeID string) error {
	timeout, cancel := context.WithTimeout(ctx, time.Second*3)
	defer cancel()

	routeUID, err := uuid.Parse(routeID)
	if err != nil {
		return fmt.Errorf("uuid.Parse: %w", err)
	}

	// if err := s.br.Publish("dino.routes", &pubsub.RouteConfig{
	// 	TunnelUID:    sqlRoute.TunnelName,
	// 	Hostname:     sqlRoute.Hostname,
	// 	DestProtocol: sqlRoute.DestinationProtocol,
	// 	DestAddr:     sqlRoute.DestinationIp,
	// 	DestPort:     uint32(sqlRoute.DestinationPort),
	// 	IsDelete:     false,
	// }); err != nil {
	// 	return Route{}, fmt.Errorf("failed to publish route creation: %w", err)
	// }

	return queries.New(s.db).DeleteRoute(timeout, routeUID)
}

// Active
func (s *serviceImpl) Active(ctx context.Context, hostname string) (string, error) {
	timeout, cancel := context.WithTimeout(ctx, time.Second*3)
	defer cancel()

	activeUID, err := queries.New(s.db).SelectActiveRoute(timeout, hostname)
	if err != nil {
		return "", fmt.Errorf("failed to execute select active route query: %w", err)
	}
	return activeUID.String(), nil
}

// Sync
func (s *serviceImpl) Sync(ctx context.Context, tunnel string) ([]Route, error) {
	timeout, cancel := context.WithTimeout(ctx, time.Second*3)
	defer cancel()

	rows, err := queries.New(s.db).SelectRoutesMany(timeout, tunnel)
	if err != nil {
		return nil, fmt.Errorf("failed to execute select routes many query: %w", err)
	}

	return dtoTunnelRoutes(rows), nil
}

func dtoRoute(r queries.DinoRoute) Route {
	return Route{
		ID:                  r.ID.String(),
		Tunnel:              r.TunnelName,
		Hostname:            r.Hostname,
		DestinationProtocol: r.DestinationProtocol,
		DestinationIP:       r.DestinationIp,
		DestinationPort:     uint32(r.DestinationPort),
		Enabled:             r.IsActive,
		CreatedAt:           r.CreatedAt.Time,
	}
}

func dtoPartials(rs []queries.SelectManyRoutesRow) []RoutePartial {
	s := make([]RoutePartial, 0, len(rs))
	for _, r := range rs {
		s = append(s, dtoPartial(r))
	}
	return s
}

func dtoPartial(r queries.SelectManyRoutesRow) RoutePartial {
	return RoutePartial{
		ID:        r.ID.String(),
		Hostname:  r.Hostname,
		Enabled:   r.IsActive,
		CreatedAt: r.CreatedAt.Time,
	}
}

func dtoTunnelRoutes(s []queries.SelectRoutesManyRow) []Route {
	rs := make([]Route, 0, len(s))
	for _, r := range s {
		rs = append(rs, dtoTunnelRoute(r))
	}
	return rs
}

func dtoTunnelRoute(r queries.SelectRoutesManyRow) Route {
	return Route{
		ID:                  r.ID.String(),
		Tunnel:              r.ID_2.String(),
		Hostname:            r.Hostname,
		DestinationProtocol: r.DestinationProtocol,
		DestinationIP:       r.DestinationIp,
		DestinationPort:     uint32(r.DestinationPort),
		Enabled:             r.IsActive,
		CreatedAt:           r.CreatedAt.Time,
		UpdatedAt:           &r.UpdatedAt.Time,
	}
}
