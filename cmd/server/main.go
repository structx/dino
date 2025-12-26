package main

import (
	"github.com/structx/dino/auth"
	"github.com/structx/dino/database"
	"github.com/structx/dino/database/migrate"
	"github.com/structx/dino/gateway"
	"github.com/structx/dino/gateway/interceptors"
	"github.com/structx/dino/internal/routes"
	"github.com/structx/dino/internal/tunnel"
	"github.com/structx/dino/logging"
	"github.com/structx/dino/migrations"
	"github.com/structx/dino/proxy"
	"github.com/structx/dino/sessions"
	"github.com/structx/dino/setup"
	tunnelgateway "github.com/structx/dino/tunnel/gateway"
	tunnelinterceptors "github.com/structx/dino/tunnel/gateway/interceptors"
	tunnelserver "github.com/structx/dino/tunnel/rpc/server"
	"github.com/structx/dino/tunnel/verifier"
	"go.uber.org/fx"
)

var opts = fx.Options(
	setup.Module,    // server config
	logging.Module,  // uber/zap logger
	database.Module, // pgx connector
	auth.Module,     // jwt authenticator
	verifier.Module, // jwt verifier

	migrations.Module, // database migrations fixtures
	migrate.Module,    // pgx database migrations
	tunnel.Module,     // tunnel service logic
	routes.Module,     // routes service logic
	sessions.Module,   // tunnel session manager

	proxy.Module,        // http proxy handler
	interceptors.Module, // gateway interceptors
	gateway.Module,      // api and proxy gateway

	tunnelserver.Module,       // tunnel gRPC server
	tunnelinterceptors.Module, // tunnel gateway interceptors
	tunnelgateway.Module,      // tunnel gateway
)

func main() {
	fx.New(opts).Run()
}
