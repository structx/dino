package main

import (
	"github.com/structx/teapot"
	teafx "github.com/structx/teapot/adapter/fx"
	"go.uber.org/fx"
	"go.uber.org/fx/fxevent"
	"soft.structx.io/dino/auth"
	"soft.structx.io/dino/database"
	"soft.structx.io/dino/database/migrate"
	"soft.structx.io/dino/gateway"
	"soft.structx.io/dino/gateway/interceptors"
	"soft.structx.io/dino/internal/routes"
	"soft.structx.io/dino/internal/tunnel"
	"soft.structx.io/dino/logging"
	"soft.structx.io/dino/migrations"
	"soft.structx.io/dino/proxy"
	"soft.structx.io/dino/pubsub"
	"soft.structx.io/dino/sessions"
	"soft.structx.io/dino/setup"
	tunnelgateway "soft.structx.io/dino/tunnel/gateway"
	tunnelinterceptors "soft.structx.io/dino/tunnel/gateway/interceptors"
	tunnelserver "soft.structx.io/dino/tunnel/rpc/server"
	"soft.structx.io/dino/tunnel/verifier"
)

var opts = fx.Options(
	setup.Module, // server config
	logging.Module, fx.WithLogger(func(l *teapot.Logger) fxevent.Logger {
		return teafx.New(l)
	}),
	database.Module, // pgx connector
	auth.Module,     // jwt authenticator
	verifier.Module, // jwt verifier
	pubsub.Module,   // pubsub broker

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
