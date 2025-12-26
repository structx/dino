package main

import (
	"github.com/structx/dino/logging"
	"github.com/structx/dino/setup"
	"github.com/structx/dino/tunnel/router"
	"github.com/structx/dino/tunnel/rpc/client"
	"github.com/structx/dino/tunnel/sessions"
	"go.uber.org/fx"
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
