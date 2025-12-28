package main

import (
	"go.uber.org/fx"
	"soft.structx.io/dino/logging"
	"soft.structx.io/dino/setup"
	"soft.structx.io/dino/tunnel/router"
	"soft.structx.io/dino/tunnel/rpc/client"
	"soft.structx.io/dino/tunnel/sessions"
)

var opts = fx.Options(
	setup.Module,
	logging.Module,
	router.Module,
	client.Module,
	sessions.Module,
)

func main() {
	fx.New(opts).Run()
}
