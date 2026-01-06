package database

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/structx/teapot"
	"go.uber.org/fx"
	"soft.structx.io/dino/setup"
)

// CommandTag
type CommandTag = pgconn.CommandTag

// Row
type Row = pgx.Row

// Rows
type Rows = pgx.Rows

// DBTX
type DBTX interface {
	Exec(context.Context, string, ...interface{}) (CommandTag, error)
	Query(context.Context, string, ...interface{}) (Rows, error)
	QueryRow(context.Context, string, ...interface{}) Row

	Ping(context.Context) error
}

// Params
type Params struct {
	fx.In

	Logger *teapot.Logger

	Lc fx.Lifecycle

	Cfg *setup.DB
}

// Result
type Result struct {
	fx.Out

	DB DBTX
}

// Module
var Module = fx.Module("database", fx.Provide(newModule))

func newModule(p Params) (Result, error) {

	dsn := p.Cfg.Dial()
	dbconfig, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		return Result{}, fmt.Errorf("pgxpool.ParseConfig: %w", err)
	}

	pool, err := pgxpool.NewWithConfig(context.TODO(), dbconfig)
	if err != nil {
		return Result{}, fmt.Errorf("pgxpool.NewWithConfig: %w", err)
	}

	p.Lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			p.Logger.Info("ping database connection")
			return pool.Ping(ctx)
		},
		OnStop: func(ctx context.Context) error {
			p.Logger.Info("close database connection")
			pool.Close()
			return nil
		},
	})

	return Result{
		DB: pool,
	}, nil
}
