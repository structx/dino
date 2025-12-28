package tunnel

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"time"

	"github.com/google/uuid"
	"soft.structx.io/dino/database"
	"soft.structx.io/dino/internal/tunnel/queries"
)

// TunnelCreate
type TunnelCreate struct {
	Name string
}

// Tunnel
type Tunnel struct {
	ID        string
	Name      string
	CreatedAt time.Time
	UpdatedAt *time.Time
}

// TunnelPartial
type TunnelPartial struct {
	ID        string
	Name      string
	CreatedAt time.Time
}

// TunnelUpdate
type TunnelUpdate struct {
	OldName string
	Name    string
}

// SecretKey
type SecretKey struct {
	Secret string
}

// Service
type Service interface {
	// Create
	Create(context.Context, TunnelCreate) (Tunnel, SecretKey, error)
	// Get
	Get(context.Context, string) (Tunnel, error)
	// List
	List(context.Context, int32, int32) ([]TunnelPartial, error)
	// Update
	Update(context.Context, TunnelUpdate) (Tunnel, error)
	// Delete
	Delete(context.Context, string) error

	// VerifyToken
	VerifyToken(context.Context, string) (string, error)

	// Health
	Health(context.Context) error
}

type serviceImpl struct {
	dbtx database.DBTX
}

// interface compliance
var _ Service = (*serviceImpl)(nil)

func newService(db database.DBTX) Service {
	return &serviceImpl{dbtx: db}
}

// Create
func (s *serviceImpl) Create(ctx context.Context, create TunnelCreate) (Tunnel, SecretKey, error) {
	timeout, cancel := context.WithTimeout(ctx, time.Second*3)
	defer cancel()

	// generate random secret
	secretKey := make([]byte, 11)
	if _, err := rand.Read(secretKey); err != nil {
		return Tunnel{}, SecretKey{}, fmt.Errorf("rand.Read: %w", err)
	}
	secretStr := hex.EncodeToString(secretKey)

	tokenHash, err := hashToken(secretStr)
	if err != nil {
		return Tunnel{}, SecretKey{}, err
	}

	sqlTunnel, err := queries.New(s.dbtx).InsertTunnel(timeout, queries.InsertTunnelParams{
		Identifier: create.Name,
		TokenHash:  tokenHash,
	})
	if err != nil {
		return Tunnel{}, SecretKey{}, fmt.Errorf("failed to execute insert tunnel query: %w", err)
	}

	return dtoTunnel(sqlTunnel), SecretKey{Secret: sqlTunnel.TokenHash}, nil
}

// Delete implements Service.
func (s *serviceImpl) Delete(ctx context.Context, tunnelName string) error {
	timeout, cancel := context.WithTimeout(ctx, time.Second*3)
	defer cancel()

	_, err := queries.New(s.dbtx).DeleteTunnel(timeout, tunnelName)
	if err != nil {
		return fmt.Errorf("failed to execute delete tunnel query: %w", err)
	}

	return nil
}

// Get implements Service.
func (s *serviceImpl) Get(ctx context.Context, tunnelName string) (Tunnel, error) {
	timeout, cancel := context.WithTimeout(ctx, time.Second*3)
	defer cancel()

	sqlTunnel, err := queries.New(s.dbtx).SelectTunnel(timeout, tunnelName)
	if err != nil {
		return Tunnel{}, fmt.Errorf("failed to execute select tunnel query: %w", err)
	}

	return dtoTunnel(sqlTunnel), nil
}

// List implements Service.
func (s *serviceImpl) List(ctx context.Context, limit int32, offset int32) ([]TunnelPartial, error) {
	timeout, cancel := context.WithTimeout(ctx, time.Second*3)
	defer cancel()

	partials, err := queries.New(s.dbtx).ListTunnels(timeout, queries.ListTunnelsParams{
		Limit:  limit,
		Offset: offset,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to execute list tunnels query: %w", err)
	}
	fmt.Println(partials)
	return dtoTunnels(partials), nil
}

func dtoTunnels(rows []queries.ListTunnelsRow) []TunnelPartial {
	tps := make([]TunnelPartial, 0, len(rows))
	for _, r := range rows {
		tps = append(tps, TunnelPartial{
			ID:        r.ID.String(),
			Name:      r.Identifier,
			CreatedAt: r.CreatedAt.Time,
		})
	}
	return tps
}

// Update implements Service.
func (s *serviceImpl) Update(ctx context.Context, args TunnelUpdate) (Tunnel, error) {
	timeout, cancel := context.WithTimeout(ctx, time.Second*3)
	defer cancel()

	params := queries.UpdateTunnelParams{
		Identifier:   args.OldName,
		Identifier_2: args.Name,
	}

	sqlTunnel, err := queries.New(s.dbtx).UpdateTunnel(timeout, params)
	if err != nil {
		return Tunnel{}, fmt.Errorf("failed to execute update tunnel query: %w", err)
	}

	return dtoTunnel(sqlTunnel), nil
}

// VerifyToken
func (s *serviceImpl) VerifyToken(ctx context.Context, tunnelID string) (string, error) {
	timeout, cancel := context.WithTimeout(ctx, time.Second*3)
	defer cancel()

	tunnelUID, err := uuid.Parse(tunnelID)
	if err != nil {
		return "", fmt.Errorf("uuid.Parse: %w", err)
	}

	combinedHash, err := queries.New(s.dbtx).SelectTunnelToken(timeout, tunnelUID)
	if err != nil {
		return "", fmt.Errorf("failed to execute select tunnel with token query: %w", err)
	}

	return combinedHash, nil
}

// Health
func (s *serviceImpl) Health(ctx context.Context) error {
	return s.dbtx.Ping(ctx)
}

func dtoTunnel(t queries.DinoTunnel) Tunnel {

	var updatedAt *time.Time
	if t.UpdatedAt.Valid {
		updatedAt = &t.UpdatedAt.Time
	}

	return Tunnel{
		ID:        t.ID.String(),
		Name:      t.Identifier,
		CreatedAt: t.CreatedAt.Time,
		UpdatedAt: updatedAt,
	}
}
