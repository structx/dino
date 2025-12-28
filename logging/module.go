package logging

import (
	"go.uber.org/fx"
	"go.uber.org/fx/fxevent"
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
var Module = fx.Module("zap_logger", fx.Provide(newModule), fx.WithLogger(func(logger *zap.Logger) fxevent.Logger {
	return &fxevent.ZapLogger{Logger: logger}
}))

func newModule(p Params) Result {
	log, _ := zap.NewDevelopment()
	return Result{
		Logger: log,
	}
}
