package logging

import (
	"os"
	"strings"

	"github.com/structx/teapot"
	"go.uber.org/fx"
	"soft.structx.io/dino/setup"
)

// Params
type Params struct {
	fx.In

	Lc fx.Lifecycle

	Cfg *setup.Logger
}

// Result
type Result struct {
	fx.Out

	Logger *teapot.Logger
}

// Module
var Module = fx.Module("teapot_logger", fx.Provide(newModule))

func newModule(p Params) Result {

	var level teapot.Level
	switch strings.ToLower(p.Cfg.Level) {
	case "fatal":
		level = teapot.FATAL
	case "error":
		level = teapot.ERROR
	case "info":
		level = teapot.INFO
	default:
		level = teapot.DEBUG
	}

	log := teapot.New(
		teapot.WithLevel(level),
		teapot.WithWriter(os.Stdout),
	)
	return Result{
		Logger: log,
	}
}
