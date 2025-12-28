package pubsub

import "go.uber.org/fx"

// Result
type Result struct {
	fx.Out

	Broker Broker
}

// Module
var Module = fx.Module("pubsub", fx.Provide(newModule))

func newModule() Result {
	return Result{
		Broker: newBroker(),
	}
}
