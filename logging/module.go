package logging

import (
	"go.uber.org/fx"
	"go.uber.org/zap"
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

	Logger *zap.Logger
}

// Module
var Module = fx.Module("zap_logger", fx.Provide(newModule))

func newModule(p Params) Result {
	log, _ := zap.NewDevelopment()
	return Result{
		Logger: log,
	}
}
