package container

import (
	"context"
	"fmt"
	"net"
	"testing"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/ory/dockertest/v3"
	"github.com/ory/dockertest/v3/docker"
	"github.com/structx/dino/database"
	"github.com/structx/dino/setup"
	"go.uber.org/fx"
)

// Params
type Params struct {
	fx.In

	Lc fx.Lifecycle
	Tb testing.TB

	Cfg *setup.DB
}

// Result
type Result struct {
	fx.Out

	DB database.DBTX
}

// Module
var Module = fx.Module("database_container_test", fx.Provide(newModule))

func newModule(p Params) (Result, error) {
	ctx := p.Tb.Context()

	pool, err := dockertest.NewPool("")
	if err != nil {
		return Result{}, fmt.Errorf("dockertest.NewPool: %w", err)
	}

	err = pool.Client.Ping()
	if err != nil {
		return Result{}, fmt.Errorf("docker client ping: %w", err)
	}

	opts := &dockertest.RunOptions{
		Repository: "postgres",
		Tag:        "18-alpine3.22",
		Env: []string{
			fmt.Sprintf("POSTGRES_PASSWORD=%s", p.Cfg.Password),
			fmt.Sprintf("POSTGRES_USER=%s", p.Cfg.Username),
			fmt.Sprintf("POSTGRES_DB=%s", p.Cfg.Name),
			"listen_addresses='*'",
		},
	}
	resource, err := pool.RunWithOptions(opts, func(config *docker.HostConfig) {
		config.AutoRemove = true
		config.RestartPolicy = docker.NeverRestart()
	})
	if err != nil {
		return Result{}, fmt.Errorf("pool.RunWithOptions: %w", err)
	}

	hostAndPort := resource.GetHostPort("5432/tcp")
	host, port, err := net.SplitHostPort(hostAndPort)
	if err != nil {
		return Result{}, fmt.Errorf("net.SplitHostPort: %w", err)
	}

	p.Cfg.Host = host
	p.Cfg.Port = port

	_ = resource.Expire(45)

	pool.MaxWait = time.Second * 30
	var dbtx database.DBTX
	if err := pool.Retry(func() error {
		cfg, err := pgxpool.ParseConfig(p.Cfg.Dial())
		if err != nil {
			return fmt.Errorf("pgxpool.ParseConfig: %w", err)
		}

		dbPool, err := pgxpool.NewWithConfig(ctx, cfg)
		if err != nil {
			return fmt.Errorf("pgxpool.NewWithConfig: %w", err)
		}

		dbtx = dbPool
		return dbPool.Ping(ctx)
	}); err != nil {
		return Result{}, fmt.Errorf("pool.Retry: %w", err)
	}

	p.Lc.Append(fx.Hook{
		OnStop: func(ctx context.Context) error {
			return pool.Purge(resource)
		},
	})

	return Result{
		DB: dbtx,
	}, nil
}
