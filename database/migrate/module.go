package migrate

import (
	"embed"
	"errors"
	"fmt"

	"github.com/golang-migrate/migrate/v4"
	pgx "github.com/golang-migrate/migrate/v4/database/pgx/v5"
	"github.com/golang-migrate/migrate/v4/source/iofs"
	"go.uber.org/fx"
	"go.uber.org/zap"
	"soft.structx.io/dino/setup"
)

// Params
type Params struct {
	fx.In

	Logger *zap.Logger

	Cfg *setup.DB

	EmbedFS []embed.FS `group:"migrations"`
}

type zapAdapter struct {
	logger *zap.Logger
}

// Printf implements migrate.Logger.
func (z *zapAdapter) Printf(format string, v ...interface{}) {
	z.logger.Info(format, zap.Any("v", v))
}

// Verbose implements migrate.Logger.
func (z *zapAdapter) Verbose() bool {
	return true
}

var _ migrate.Logger = (*zapAdapter)(nil)

// Module
var Module = fx.Module("database_migrate", fx.Invoke(invokeModule))

func invokeModule(p Params) error {

	dsn := p.Cfg.Dial()

	px := pgx.Postgres{}
	databaseDriver, err := px.Open(dsn)
	if err != nil {
		return fmt.Errorf("px.Open: %w", err)
	}

	for _, embedFS := range p.EmbedFS {

		sourceDriver, err := iofs.New(embedFS, "fixtures")
		if err != nil {
			return fmt.Errorf("iofs.New: %w", err)
		}

		m, err := migrate.NewWithInstance("iofs", sourceDriver, "pgx/v5", databaseDriver)
		if err != nil {
			return fmt.Errorf("migrate.NewWithInstance: %w", err)
		}

		m.Log = &zapAdapter{logger: p.Logger.Named("go_migrate")}

		if err := m.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
			return fmt.Errorf("m.Up: %w", err)
		}
	}

	return nil
}
