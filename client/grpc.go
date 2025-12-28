package client

import (
	"context"
	"crypto/tls"
	"fmt"
	"net"
	"time"

	"github.com/google/uuid"
	"golang.org/x/net/http2"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	pbroutes "soft.structx.io/dino/pb/routes/v1"
	pbtunnels "soft.structx.io/dino/pb/tunnels/v1"
)

type TunnelAdd struct {
	Name string
}

type Tunnel struct {
	UID       uuid.UUID
	Name      string
	CreatedAt time.Time
	UpdatedAt *time.Time
}

type TunnelList struct {
	Limit        int32
	Offset       int32
	Autocomplete bool
	CompleteMe   string
}

type TunnelPartial struct {
	UID  uuid.UUID
	Name string
}

type TunnelUpdate struct {
	OldName string
	Name    string
}

type RouteAdd struct {
	Tunnel              string
	Hostname            string
	DestinationProtocol string
	DestinationIP       string
	DestinationPort     uint32
	Enabled             bool
}

type RouteDel struct{}

type Route struct {
	Name    string
	Enabled bool
	Tunnel  string
}

type RoutePartial struct {
	UID       uuid.UUID
	Hostname  string
	Enabled   bool
	CreatedAt time.Time
}

type RouteUpdate struct {
	UID                 uuid.UUID
	Hostname            string
	DestinationProtocol string
	DestinationIP       string
	DestinationPort     uint32
	Enabled             bool
}

type Auth interface{}

type SharedSecret struct {
	Key string
}

type Client interface {
	AddTunnel(context.Context, TunnelAdd) (Tunnel, Auth, error)
	GetTunnel(context.Context, string) (Tunnel, error)
	ListTunnels(context.Context, TunnelList) ([]TunnelPartial, error)
	UpdateTunnel(context.Context, TunnelUpdate) (Tunnel, error)
	DelTunnel(context.Context, string) error

	AddRoute(context.Context, RouteAdd) (Route, error)
	GetRoute(context.Context, string) (Route, error)
	ListRoutes(context.Context, string, int32, int32) ([]RoutePartial, error)
	UpdateRoute(context.Context, RouteUpdate) (Route, error)
	DelRoute(context.Context, string) error

	// Close client conn
	Close() error
}

type clientImpl struct {
	target string

	conn *grpc.ClientConn
}

// interface compliance
var _ Client = (*clientImpl)(nil)

// ClientOption
type ClientOption func(*clientImpl)

// WithTarget
func WithTarget(t string) ClientOption {
	return func(c *clientImpl) {
		c.target = t
	}
}

// New
func New(options ...ClientOption) (Client, error) {
	cli := &clientImpl{}

	for _, opt := range options {
		opt(cli)
	}

	tr := &http2.Transport{
		AllowHTTP: false,
		DialTLSContext: func(ctx context.Context, network, addr string, cfg *tls.Config) (net.Conn, error) {
			return net.Dial(network, addr)
		},
	}

	conn, err := grpc.NewClient(cli.target,
		grpc.WithContextDialer(func(ctx context.Context, s string) (net.Conn, error) {
			return tr.DialTLSContext(ctx, "tcp", s, nil)
		}),
		grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, fmt.Errorf("grpc.NewClient: %w", err)
	}
	cli.conn = conn

	return cli, nil
}

// AddTunnel
func (c *clientImpl) AddTunnel(ctx context.Context, newTunnel TunnelAdd) (Tunnel, Auth, error) {
	cli := pbtunnels.NewTunnelServiceClient(c.conn)

	// timeout, cancel := context.WithTimeout(ctx, time.Second*3)
	// defer cancel()

	req := &pbtunnels.CreateTunnelRequest{TunnelName: newTunnel.Name}
	resp, err := cli.CreateTunnel(ctx, req)
	if err != nil {
		return Tunnel{}, nil, fmt.Errorf("failed to execute gRPC create tunnel: %w", err)
	}

	return dtoTunnel(resp.Tunnel), SharedSecret{Key: resp.GetSecretKey()}, nil
}

// GetTunnel
func (c *clientImpl) GetTunnel(ctx context.Context, tunnelName string) (Tunnel, error) {
	cli := pbtunnels.NewTunnelServiceClient(c.conn)

	timeout, cancel := context.WithTimeout(ctx, time.Second*3)
	defer cancel()

	req := &pbtunnels.GetTunnelRequest{Name: tunnelName}

	resp, err := cli.GetTunnel(timeout, req)
	if err != nil {
		return Tunnel{}, fmt.Errorf("failed to execute gRPC get tunnel: %w", err)
	}

	return dtoTunnel(resp.Tunnel), nil
}

// ListTunnels
func (c *clientImpl) ListTunnels(ctx context.Context, args TunnelList) ([]TunnelPartial, error) {
	cli := pbtunnels.NewTunnelServiceClient(c.conn)

	timeout, cancel := context.WithTimeout(ctx, time.Second*3)
	defer cancel()

	req := &pbtunnels.ListTunnelsRequest{
		Limit:      args.Limit,
		Offset:     args.Offset,
		ToComplete: args.Autocomplete,
		Complete:   args.CompleteMe,
	}

	resp, err := cli.ListTunnels(timeout, req)
	if err != nil {
		return nil, fmt.Errorf("failed to execute gRPC list tunnels: %w", err)
	}

	return dtoTunnelPartials(resp.Tunnels), nil
}

// UpdateTunnel
func (c *clientImpl) UpdateTunnel(ctx context.Context, args TunnelUpdate) (Tunnel, error) {
	cli := pbtunnels.NewTunnelServiceClient(c.conn)

	timeout, cancel := context.WithTimeout(ctx, time.Second*3)
	defer cancel()

	req := &pbtunnels.UpdateTunnelRequest{
		TunnelUpdate: &pbtunnels.TunnelUpdate{
			OldName: args.OldName,
			NewName: args.Name,
		},
	}

	resp, err := cli.UpdateTunnel(timeout, req)
	if err != nil {
		return Tunnel{}, fmt.Errorf("failed to execute gRPC update tunnel: %w", err)
	}

	return dtoTunnel(resp.Tunnel), nil
}

// DelTunnel
func (c *clientImpl) DelTunnel(ctx context.Context, tunnelName string) error {

	cli := pbtunnels.NewTunnelServiceClient(c.conn)

	timeout, cancel := context.WithTimeout(ctx, time.Second*3)
	defer cancel()

	req := &pbtunnels.DeleteTunnelRequest{Name: tunnelName}
	_, err := cli.DeleteTunnel(timeout, req)
	if err != nil {
		return fmt.Errorf("failed to execute gRPC delete tunnel: %w", err)
	}

	return nil
}

// AddRoute
func (c *clientImpl) AddRoute(ctx context.Context, args RouteAdd) (Route, error) {

	cli := pbroutes.NewRouteServiceClient(c.conn)

	timeout, cancel := context.WithTimeout(ctx, time.Second*3)
	defer cancel()

	req := &pbroutes.CreateRouteRequest{
		Create: &pbroutes.RouteCreate{
			Tunnel:       args.Tunnel,
			Hostname:     args.Hostname,
			DestProtocol: args.DestinationProtocol,
			DestAddr:     args.DestinationIP,
			DestPort:     args.DestinationPort,
		},
	}

	resp, err := cli.CreateRoute(timeout, req)
	if err != nil {
		return Route{}, fmt.Errorf("failed to execute gRPC crete route: %w", err)
	}

	return dtoRoute(resp.Route), nil
}

// GetRoute
func (c *clientImpl) GetRoute(ctx context.Context, hostname string) (Route, error) {

	cli := pbroutes.NewRouteServiceClient(c.conn)

	timeout, cancel := context.WithTimeout(ctx, time.Second*3)
	defer cancel()

	req := &pbroutes.GetRouteRequest{Hostname: hostname}
	resp, err := cli.GetRoute(timeout, req)
	if err != nil {
		return Route{}, fmt.Errorf("failed to execute gRPC get route: %w", err)
	}

	return dtoRoute(resp.Route), nil
}

// ListRoutes
func (c *clientImpl) ListRoutes(ctx context.Context, tunnel string, limit int32, offset int32) ([]RoutePartial, error) {
	cli := pbroutes.NewRouteServiceClient(c.conn)

	timeout, cancel := context.WithTimeout(ctx, time.Second*3)
	defer cancel()

	req := &pbroutes.ListRoutesRequest{Tunnel: tunnel, Limit: limit, Offset: offset}
	resp, err := cli.ListRoutes(timeout, req)
	if err != nil {
		return nil, fmt.Errorf("failed to execute GRPC list routes: %w", err)
	}

	return dtoRoutePartials(resp.Partials), nil
}

// UpdateRoute
func (c *clientImpl) UpdateRoute(ctx context.Context, args RouteUpdate) (Route, error) {
	cli := pbroutes.NewRouteServiceClient(c.conn)

	timeout, cancel := context.WithTimeout(ctx, time.Second*3)
	defer cancel()

	req := &pbroutes.UpdateRouteRequest{
		Update: &pbroutes.RouteUpdate{
			Uid:          args.UID.String(),
			Hostname:     args.Hostname,
			DestProtocol: args.DestinationProtocol,
			DestAddr:     args.DestinationIP,
			DestPort:     args.DestinationPort,
		},
	}

	resp, err := cli.UpdateRoute(timeout, req)
	if err != nil {
		return Route{}, fmt.Errorf("failed to execute gRPC update route: %w", err)
	}

	return dtoRoute(resp.Route), nil
}

// DelRoute
func (c *clientImpl) DelRoute(ctx context.Context, hostname string) error {
	cli := pbroutes.NewRouteServiceClient(c.conn)

	timeout, cancel := context.WithTimeout(ctx, time.Second*3)
	defer cancel()

	req := &pbroutes.DeleteRouteRequest{Hostname: hostname}
	_, err := cli.DeleteRoute(timeout, req)
	if err != nil {
		return fmt.Errorf("failed to execute gRPC delete route: %w", err)
	}

	return nil
}

// Close
func (c *clientImpl) Close() error {
	return c.conn.Close()
}

func dtoTunnelPartials(pbps []*pbtunnels.TunnelPartial) []TunnelPartial {
	tps := make([]TunnelPartial, 0, len(pbps))
	for _, p := range pbps {
		tps = append(tps, dtoTunnelPartial(p))
	}
	return tps
}

func dtoTunnelPartial(t *pbtunnels.TunnelPartial) TunnelPartial {
	return TunnelPartial{
		Name: t.Name,
	}
}

func dtoTunnel(t *pbtunnels.Tunnel) Tunnel {
	var updatedAtPtr *time.Time
	if t.UpdatedAt.IsValid() {
		updatedAt := t.UpdatedAt.AsTime()
		updatedAtPtr = &updatedAt
	}
	return Tunnel{
		UID:       uuid.MustParse(t.Id),
		Name:      t.Name,
		CreatedAt: t.CreatedAt.AsTime(),
		UpdatedAt: updatedAtPtr,
	}
}

func dtoRoute(r *pbroutes.Route) Route {
	return Route{
		Name:    r.Name,
		Enabled: r.Enabled,
		Tunnel:  r.Tunnel,
	}
}

func dtoRoutePartials(pbrs []*pbroutes.RoutePartial) []RoutePartial {
	rps := make([]RoutePartial, 0, len(pbrs))
	for _, p := range pbrs {
		rps = append(rps, dtoRoutePartial(p))
	}
	return rps
}

func dtoRoutePartial(p *pbroutes.RoutePartial) RoutePartial {
	return RoutePartial{
		UID:      uuid.MustParse(p.Uid),
		Hostname: p.Hostname,
	}
}
